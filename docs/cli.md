## Command Flags

### Available commands

```
darron@: kvexpress -h
Small Go program to put and pull configuration data out of Consul and write to filesystem.

Usage:
  kvexpress [flags]
  kvexpress [command]

Available Commands:
  clean       Clean local cache files.
  copy        Copy a Consul key to another location.
  in          Put configuration into Consul.
  lock        Lock a file on a single node so it stays the way it is.
  out         Write a file based on kvexpress organized data stored in Consul.
  raw         Write a file pulled from any Consul KV data.
  stop        Put stop value into Consul.
  unlock      Unock a file on a single node so it updates.
```

### Global Flags

```
Global Flags:
  -c, --chmod int                  permissions for the file (default 416)
  -z, --compress                   gzip in and out of the KV store
  -C, --config string              Config file location
  -a, --datadog_api_key string     Datadog API Key
  -A, --datadog_app_key string     Datadog App Key
  -d, --dogstatsd                  send metrics to dogstatsd
  -D, --dogstatsd_address string   address for dogstatsd server (default "localhost:8125")
  -e, --exec string                Execute this command after
  -l, --length int                 minimum amount of lines in the file (default 10)
  -o, --owner string               who to write the file as
  -p, --prefix string              prefix for the key (default "kvexpress")
  -s, --server string              Consul server location (default "localhost:8500")
  -t, --token string               Token for Consul access (default "anonymous")
      --verbose                    log output to stdout
```

* [clean](#clean-command-flags)
* [copy](#copy-command-flags)
* [in](#in-command-flags)
* [lock](#lock-command-flags)
* [out](#out-command-flags)
* [raw](#raw-command-flags)
* [stop](#stop-command-flags)
* [unlock](#unlock-command-flags)

### `clean` command flags

```
darron@: kvexpress clean -h
clean is for cleaning up local cache files.

Usage:
  kvexpress clean [flags]

Flags:
  -f, --file string   file to clean
```

Example Command:

`kvexpress clean -f /etc/consul-template/output/hosts.consul`

### `copy` command flags

```
darron@: kvexpress copy -h
copy is for copying already existing keys.

Usage:
  kvexpress copy [flags]

Flags:
      --keyfrom string   key to pull data from
      --keyto string     key to write the data to
```

Example Command:

`kvexpress copy --keyfrom "hosts" --keyto "hosts_alternate"`

### `in` command flags

```
darron@: kvexpress in -h
in is for putting data into a Consul key so that you can write it on another networked node.

Usage:
  kvexpress in [flags]

Flags:
  -f, --file string   filename to read data from
  -k, --key string    key to push data to
  -S, --sorted        sort the input file
  -u, --url string    url to read data from
```

Example Command:

`kvexpress in -d true -k hosts -f /etc/consul-template/output/hosts.consul -l 100 --sorted=true`


### `lock` command flags

```
darron@: kvexpress lock -h
Lock is a convenient way to stop a file from being updated on a single node.

Usage:
  kvexpress lock [flags]

Flags:
  -f, --file string     file to lock
  -r, --reason string   reason to lock
```

Example Command:

`kvexpress lock -f /etc/hosts.consul -r "I need this file to be locked for an hour."`

### `out` command flags

```
darron@: kvexpress out -h
out is for writing a file based on a Consul kvexpress key and checksum.

Usage:
  kvexpress out [flags]

Flags:
  -f, --file string   where to write the data
      --ignore_stop   ignore stop key
  -k, --key string    key to pull data from
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
### `raw` command flags

```
darron@: kvexpress raw -h
raw is for writing a file based on any Consul key.

Usage:
  kvexpress raw [flags]

Flags:
  -f, --file string   where to write the data
  -k, --key string    Raw key to pull data from
```

Example Command:

`kvexpress raw -f /etc/hosts.consul -k kvexpress/hosts/data`

### `stop` command flags

```
darron@: kvexpress stop -h
stop is a convenient way to put stop values in Consul.  Stops ALL nodes from updating.

Usage:
  kvexpress stop [flags]

Flags:
  -k, --key string      key to stop
  -r, --reason string   reason to stop
```

Example Command:

`kvexpress stop -k hosts -r "Restarting all Consul nodes - want this to be locked for all nodes."`

### `unlock` command flags

```
darron@: kvexpress unlock -h
Unlock is a convenient way to allow a previously locked file to be updated.

Usage:
  kvexpress unlock [flags]

Flags:
  -f, --file string   file to unlock
```

Example Command:

`kvexpress unlock -f /etc/hosts.consul`
