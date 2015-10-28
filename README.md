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
5. Run some other random command after writing the file.
6. Verify that the file we put into the KV was the same file that was written on the other end.
7. Stop the process - in or out - if we want everything to stay as it is for the moment.

**in**

```
Usage:
  kvexpress in [flags]

Flags:
  -f, --file="": filename to read data from
  -k, --key="": key to push data to
  -S, --sorted[=false]: sort the input file

Global Flags:
  -c, --chmod=416: permissions for the file
  -e, --exec="": Execute this command after
  -l, --length=10: minimum amount of lines in the file
  -p, --prefix="kvexpress": prefix for the key
  -s, --server="localhost:8500": Consul server location
  -t, --token="": Token for Consul access
```

Example: `kvexpress in -k hosts -f /etc/consul-template/output/hosts.consul -l 100 --sorted=true`

**out**

```
Usage:
  kvexpress out [flags]

Flags:
  -f, --file="": where to write the data
  -k, --key="": key to pull data from

Global Flags:
  -c, --chmod=416: permissions for the file
  -e, --exec="": Execute this command after
  -l, --length=10: minimum amount of lines in the file
  -p, --prefix="kvexpress": prefix for the key
  -s, --server="localhost:8500": Consul server location
  -t, --token="": Token for Consul access
```

Example: `kvexpress out -k nginx -f /etc/nginx/nginx.conf -e 'sudo service restart nginx'`

**Build**

To build: `make deps && make`

To build for Linux: `make deps && make linux`

Logs to to Syslog.

`./kvexpress out -h` shows you the flags you need to use.
