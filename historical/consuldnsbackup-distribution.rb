#!/usr/bin/env ruby
require 'diffy'
require 'digest'
require 'diplomat'
require 'dogapi'

Diffy::Diff.default_format = :text

Diplomat.configure do |config|
  config.acl_token =  '<%= @token %>'
end

api_key = '<%= @api_key %>'
ENV['DATADOG_HOST'] = '<%= @url %>'
dog = Dogapi::Client.new(api_key)

# If there is a `consuldnsbackup/stop` KV value - stop right here.
if reason = begin Diplomat::Kv.get('consuldnsbackup/stop') rescue nil end
  system("logger -t consul consuldnsbackup: Did NOT update hosts - stopped via consuldnsbackup/stop.")
  text = "Consul Hosts - Stopped via consuldnsbackup/stop"
  dog.emit_event(Dogapi::Event.new(
    "Reason Stopped: #{reason}", :msg_title => text, :tags => 'consuldnsbackup_stopped'
  ))
  raise reason
else
  system("logger -t consul consuldnsbackup: No explicit stop - proceeding.")
end

current_file = '/etc/consul-template/output/consul-hosts-generated.ini'
compare_file = "#{current_file}.compare"
last_file = "#{current_file}.last"
line_count_expected = '<%= @expected_count %>'

system("cat #{current_file} | sort | uniq > #{compare_file}")

# FIXME: This is a hack to prevent the script from running if there are too few lines in the file.
# In reality, when consul is running everywhere, prod value should be much higher.
compare_file_line_count = `wc -l #{compare_file}`.split.first.to_i || 0
if compare_file_line_count < line_count_expected.to_i
  text = "Consul Hosts - Not enough lines in the compare_file file, exiting"
  dog.emit_event(Dogapi::Event.new(
    "compare_file: #{compare_file}\nline count: #{compare_file_line_count}", :msg_title => text, :tags => 'consuldnsbackup_failure'
  ))
  raise text
end

hosts = File.read(compare_file) if File.exist?(compare_file)
last_hosts = File.read(last_file) if File.exist?(last_file)

if hosts && last_hosts
  hosts_sha = Digest::SHA256.hexdigest hosts
  last_hosts_sha = Digest::SHA256.hexdigest last_hosts
  diff = Diffy::Diff.new(last_hosts, hosts, :context => 2)

  text = "New File: #{hosts_sha}\nOld File: #{last_hosts_sha}\n\nDiff:\n\n#{diff}"
  title = 'Consul Hosts Update'

  if hosts_sha != last_hosts_sha
    dog.emit_event(Dogapi::Event.new(text, :msg_title => title, :tags => 'consuldnsbackup_update'))
    # Copy the file.
    FileUtils.cp "#{compare_file}", "#{last_file}"
    system("logger -t consul consuldnsbackup: Updated hosts: #{hosts_sha}")
    # Put the services into Consul.
    result = Diplomat::Kv.put('consuldnsbackup/data', hosts)
  else
    system("logger -t consul consuldnsbackup: Did NOT update hosts - files match.")
  end
end
