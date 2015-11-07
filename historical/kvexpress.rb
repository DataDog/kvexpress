#!/usr/bin/env ruby
require 'diffy'
require 'digest'
require 'diplomat'
require 'dogapi'
require 'trollop'
require 'syslogger'
require 'fileutils'
require 'statsd'

log = Syslogger.new('kvexpress', Syslog::LOG_PID, Syslog::LOG_LOCAL1)
log.level = Logger::DEBUG

Diffy::Diff.default_format = :text

# CLI Configuration from ARGV

opts = Trollop::options do
  opt :filename, 'Filename', :type => :string
  opt :key, 'Key location', :type => :string
  opt :min_lines, 'Minimum number of lines.', :default => 10
  opt :sorted, 'Should the file be sorted and uniqued', :type => :string, :default => 'no'
  opt :token, 'Alternate Token', :type => :string, :default => '<%= @token %>'
  opt :prefix, 'Alternate Prefix', :type => :string, :default => 'kvexpress'
end

Trollop::die :filename, 'Needs a filename' unless opts[:filename]
Trollop::die :key, 'Needs a Consul key location' unless opts[:key]
Trollop::die :filename, 'Filename needs to exist' unless File.exist?(opts[:filename]) if opts[:filename]

# Configuration

Diplomat.configure do |config|
  config.acl_token =  opts[:token]
end

sorted = opts[:sorted]
min_lines = opts[:min_lines]
filename = opts[:filename]
filename_compare = "#{filename}.compare"
filename_last = "#{filename}.last"
prefix = opts[:prefix]
key_short = opts[:key]
key_path = "#{prefix}/#{key_short}"
data_key = "#{key_path}/data"
checksum_key = "#{key_path}/checksum"
updated_key = "#{key_path}/updated"
stop_key = "#{key_path}/stop"
min_lines_key = "#{key_path}/min_lines"

# Datadog Config
api_key = '<%= @api_key %>'
ENV['DATADOG_HOST'] = '<%= @url %>'
dog = Dogapi::Client.new(api_key)

statsd = Statsd.new('localhost', 8125)

# Methods

def dog_event_stop(reason, key, stop_key, dog)
  text = "#{key}: kvexpress stopped via #{stop_key}"
  dog.emit_event(Dogapi::Event.new("Reason Stopped: #{reason}", :msg_title => text, :tags => [key, 'kvexpress']))
  raise reason
end

def dog_event_fatal(reason, key, dog)
  dog.emit_event(Dogapi::Event.new("#{key}: Reason Stopped: #{reason}", :msg_title => reason, :tags => [key, 'kvexpress']))
  raise reason
end

def dog_event_success(compare_sha, last_sha, diff, dog, key)
  message = %(
  %%%
  **New File:** #{compare_sha}


  **Last File:** #{last_sha}


  **Diff:**


  %%%
  #{diff})
  dog.emit_event(Dogapi::Event.new(message, :msg_title => "#{key}: Updated File", :tags => [key, 'kvexpress']))
end

def dog_event_too_short(min_lines, compare_file_lines, key, dog)
  message = %(
  %%%
  **#{key}**: Too Short


  **New File:** #{compare_file_lines}


  **Minimum Lines:** #{min_lines}
  %%%
  )
  title = "#{key}: New File NOT Long Enough"
  dog.emit_event(Dogapi::Event.new(message, :msg_title => title, :tags => [key, 'kvexpress']))
  raise title
end

def sort_and_uniq_compare_file(file, compare)
  system("cat #{file} | sort | uniq > #{compare}")
end

def create_compare_file(file, compare)
  FileUtils.cp file, compare
end

def read_file(filename)
  File.read(filename) if File.exist?(filename)
end

def sha_variable(data)
  Digest::SHA256.hexdigest data
end

def variable_diff(last, current)
  Diffy::Diff.new(last, current, :context => 2)
end

def line_count(filename)
  `wc -l #{filename}`.split.first.to_i || 0
end

def set(key, value, log)
  result = Diplomat::Kv.put(key, value)
  log.debug("method='set' key='#{key}' result='#{result.to_s}'")
end

###### Main Logic Starts Here ######

log.debug("key='#{key_short}' filename='#{filename}' key='#{key_short}' min_lines='#{min_lines}' sorted='#{sorted}'")

# Check for stop key.
if reason = begin Diplomat::Kv.get(stop_key) rescue nil end
  log.debug("key='#{key_short}' Stop key present at #{stop_key}")
  dog_event_stop(reason, key_short, stop_key, dog)
else
  log.debug("key='#{key_short}' NO Stop key present - proceeding.")
end

# Create .compare file.
unless sorted == 'no'
  sort_and_uniq_compare_file(filename, filename_compare)
  log.debug("key='#{key_short}' sorted='#{sorted}' filename='#{filename}'")
else
  create_compare_file(filename, filename_compare)
  log.debug("key='#{key_short}' copied='yes' sorted='#{sorted}' filename='#{filename}'")
end

# Check for the existance of .last - create it if it doesn't exist.
unless File.exist?(filename_last)
  log.debug("key='#{key_short}' touched='yes' filename_last='#{filename_last}'")
  FileUtils.touch(filename_last)
end

# Is the file long enough?
compare_file_lines = line_count(filename_compare)
if compare_file_lines < min_lines.to_i
  log.debug("key='#{key_short}' too_short='yes' count='#{compare_file_lines}' min_lines='#{min_lines}'")
  dog_event_too_short(min_lines, compare_file_lines, key_short, dog)
end

# Read compare and last files.
compare_data = read_file(filename_compare)
compare_data_bytes = compare_data.length
last_data = read_file(filename_last)

if compare_data && last_data
  # Compare the SHA values.
  compare_data_sha = sha_variable(compare_data)
  last_data_sha = sha_variable(last_data)

  # If there are changes - then diff them and update the key.
  if compare_data_sha != last_data_sha
    log.debug("key='#{key_short}' files_different='yes' compare_data_sha='#{compare_data_sha}' last_data_sha='#{last_data_sha}' bytes=#{compare_data_bytes} lines=#{compare_file_lines}")

    # Diff
    diff = variable_diff(last_data, compare_data)

    # Get the SHA from Consul if it exists.
    if consul_checksum = begin Diplomat::Kv.get(checksum_key) rescue nil end
      log.debug("key='#{key_short}' consul_checksum='#{consul_checksum}' compare_data_sha='#{compare_data_sha}'")
    else
      log.debug("key='#{key_short}' NO Checksum key present - proceeding.")
    end

    # Update the key as long as they don't match.
    if consul_checksum != compare_data_sha
      updated_time = Time.now.utc.to_s
      log.debug("key='#{key_short}' key_checksum_different='yes' consul_checksum='#{consul_checksum}' compare_data_sha='#{compare_data_sha}' updated_time='#{updated_time}'")
      set(data_key, compare_data, log)
      set(checksum_key, compare_data_sha, log)
      set(updated_key, updated_time, log)
      # TODO: Why does adding this break everything?
      # set(min_lines_key, min_lines, log)

      # Throw some metrics at statsd.
      statsd.batch do |s|
        s.increment('dd.kvexpress.updates', :tags => ["kvkey:#{key_short}"])
        s.gauge("dd.kvexpress.bytes", compare_data_bytes, :tags => ["kvkey:#{key_short}"])
        s.gauge("dd.kvexpress.lines", compare_file_lines, :tags => ["kvkey:#{key_short}"])
      end

      # Send to Datadog.
      dog_event_success(compare_data_sha, last_data_sha, diff, dog, key_short)

    else
      log.debug("key='#{key_short}' key_checksum_different='no' consul_checksum='#{consul_checksum.to_s}' compare_data_sha='#{compare_data_sha}' last_data_sha='#{last_data_sha}'")
    end

    # Copy compare to last.
    FileUtils.cp filename_compare, filename_last
  else
    log.debug("key='#{key_short}' files_different='no' compare_data_sha='#{compare_data_sha}' last_data_sha='#{last_data_sha}'")
  end
else
  fatal_message = 'Missing a file - this should NEVER happen.'
  log.debug("key='#{key_short}' #{fatal_message}")
  dog_event_fatal(fatal_message, key_short, dog)
end
