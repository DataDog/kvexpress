kvexpress
===============

**Why?**

Small Go utility to:

1. Put data into Consul's KV store.
2. Pull data out of Consul's KV store and write it to a file.

Why a dedicated utility though? Can't I just do it with curl?

Yes you can - but we kept wanting to:

1. Make sure the file was long enough. 0-length configuration files are bad.
2. Load the file from some other custom templating process - not just from straight KV files.
3. Put the file into any location in the filesystem.
4. Restart/reload/stop/start daemon after writing the file.
5. Run some other custom command after writing the file.
6. Verify that the file we put into the KV was the same file that was written on the other end.
7. Stop the process on all nodes - in or out - if we want everything to stay as it is for the moment.

We did this at first with some custom Ruby scripts - but the pattern was apparent and could be applied to many other files as well.

This replaces all the custom Ruby/shell scripts with a single Go binary we can use to get data in and out.

**How does it work? - 1000 foot view**

**In:** `kvexpress in --key hosts --file /etc/consul-template/output/hosts.consul --length 100 --sorted=true`

1. Check that at least `--file` and `--key` are passed along with the command. Quit if they're not present - there are no safe defaults for those flags.
2. Check for the existence of a `stop` key - if it's there - stop and exit.
3. Read the file into a string, and sort the string if requested.
4. Check if the file is long enough - if not - stop and exit.
5. Save the file to a `.compare` file - we will use this data from now on.
6. Check for the existence of a `.last` file - if it's not there - create it.
7. Are the `.compare` and `.last` files blank? If not - let's continue.
8. Compare the checksums of the `.compare` and `.last` files - if they're different - continue.
9. Grab the checksum from Consul and compare with the `.compare` file - if it's different - then let's update. This is to guard against it running on multiple server nodes that might have different `.last` files.
10. Save `data`, and `checksum` keys.
11. Copy `.compare` to `.last`
12. If `--exec` is passed - run that command.

**Out:** `kvexpress out -k hosts -f /etc/hosts.consul -l 100 -e 'sudo pkill -HUP dnsmasq'`

1. Check that at least `--file` and `--key` are passed along with the command. Quit if they're not present - there are no safe defaults for those flags.
2. Check for the existence of a `stop` key - if it's there - stop and exit.
3. Pull the `data` and `checksum` keys out of Consul.
4. If `data` is long enough and the `checksum` as computed on this side matches the `checksum` key - then continue.
5. Write the contents of `data` to the passed `--file` location.
6. If `--exec` is passed - run that command.

**`in` command flags**

```
Usage:
  kvexpress in [flags]

Flags:
  -f, --file="": filename to read data from
  -k, --key="": key to push data to
  -S, --sorted[=false]: sort the input file

Global Flags:
  -c, --chmod=416: permissions for the file
  -d, --dogstatsd[=false]: send metrics to dogstatsd
  -D, --dogstatsd_addr="localhost:8125": address for dogstatsd server
  -e, --exec="": Execute this command after
  -l, --length=10: minimum amount of lines in the file
  -p, --prefix="kvexpress": prefix for the key
  -s, --server="localhost:8500": Consul server location
  -t, --token="": Token for Consul access
```

Example: `kvexpress in -d true -k hosts -f /etc/consul-template/output/hosts.consul -l 100 --sorted=true`

**`out` command flags**

```
Usage:
  kvexpress out [flags]

Flags:
  -f, --file="": where to write the data
  -k, --key="": key to pull data from

Global Flags:
  -c, --chmod=416: permissions for the file
  -d, --dogstatsd[=false]: send metrics to dogstatsd
  -D, --dogstatsd_addr="localhost:8125": address for dogstatsd server
  -e, --exec="": Execute this command after
  -l, --length=10: minimum amount of lines in the file
  -p, --prefix="kvexpress": prefix for the key
  -s, --server="localhost:8500": Consul server location
  -t, --token="": Token for Consul access
```

Example `out` as a Consul watch:

```
{
  "watches": [
    {
      "type":"key",
      "key":"/kvexpress/hosts/checksum",
      "handler":"kvexpress out -d true -k hosts -f /etc/hosts.consul -l 100 -e 'sudo pkill -HUP dnsmasq'"
    }
  ]
}
```

**Consul KV Structure**

How are keys organized in Consul's KV store to work with kvexpress?

Underneath a global prefix `/kvexpress/` - each directory represents a specific file we are distributing through the KV store.

Each directory is named for the unique key and has the following keys underneath it:

1. `data` - where the configuration file is stored.
2. `checksum` - where the SHA256 of the data is stored.

For example - the `hosts` file is arranged like this:

```
/kvexpress/hosts/data
/kvexpress/hosts/checksum
```

There is an optional `stop` key - that if present - will cause all `in` and `out` processes to stop before writing anything. Allows us to freeze the automatic process if we need to.

**Build**

To build: `make deps && make`

To build for Linux: `make deps && make linux`

To launch an empty [Consul](https://www.consul.io/) instance: `make consul`

Logs to to Syslog.

`./kvexpress out -h` shows you the flags you need to use.
