kvexpress
===============

## Why?


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

This replaces our previous custom Ruby/shell scripts with a single Go binary we can use to get data in and out of Consul's KV store.

## How does it work? - 1000 foot view

### In

`kvexpress in --key hosts --file /etc/consul-template/output/hosts.consul --length 100 --sorted=true`

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

### Out

`kvexpress out -k hosts -f /etc/hosts.consul -l 100 -e 'sudo pkill -HUP dnsmasq'`

1. Check that at least `--file` and `--key` are passed along with the command. Quit if they're not present - there are no safe defaults for those flags.
2. Check for the existence of a `stop` key - if it's there - stop and exit.
3. Pull the `data` and `checksum` keys out of Consul.
4. If `data` is long enough and the `checksum` as computed on this side matches the `checksum` key - then continue.
5. Write the contents of `data` to the passed `--file` location.
6. If `--exec` is passed - run that command.

## Where can I get it?

Official releases can be downloaded from [the releases page](https://github.com/DataDog/kvexpress/releases).

## Commands Available

A detailed list of commands is available [here](https://github.com/DataDog/kvexpress/tree/master/docs/cli.md).

## Consul KV Structure

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

## Logging

All logs are sent to syslog and are tagged with `kvexpress`. To enable debug logs, please `export KVEXPRESS_DEBUG=1`

## Build

To build: `make deps && make`

To run integration tests: `make deps && make && make test` - it will spin up a blank Consul and kill it after the run.

Because we use `user.Current()` - you can't cross compile this. If you want to build for Linux - you must build on Linux. [Closed Issue](https://github.com/DataDog/kvexpress/issues/51#issuecomment-170307910)

To install Consul - [there are instructions here](https://www.consul.io/intro/getting-started/install.html).

To launch an empty [Consul](https://www.consul.io/) instance: `make consul`
