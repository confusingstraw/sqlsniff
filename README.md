# MySQL Sniff

This repo contains Golang code for performing basic reconnaissance of a known MySQL host. It accepts as arguments an IP address and a port. It will output plain text to stdout which describes the information (if any) it was able to glean.

## Building

Ensure you have a version of Golang installed which is able to support Gomodules (though this doesn't _really_ use them). To build, simply run `go build`.

## Usage

To see the commands available, run the binary with the `-h` flag:

```
./sqlsniff.exe -h
>>>Usage of sqlsniff:
>>>  -address string
>>>        Host address to connect to, without port (default "127.0.0.1")
>>>  -port uint
>>>        Port to connect to (default 3306)
```

An example of running this against a MySQL server available via loopback:
```
./sqlsniff -address="127.0.0.1" -port=3306
>>>2020/10/27 18:56:29 AuthenticationPlugin: caching_sha2_password
>>>2020/10/27 18:56:29 Capabilities: 3355443199
>>>2020/10/27 18:56:29 ConnectionId: 66
>>>2020/10/27 18:56:29 ProtocolVersion: 10
>>>2020/10/27 18:56:29 ServerVersion: 8.0.22
>>>2020/10/27 18:56:29 Status: 2
```

More detailed information on the underlying struct/its types can be found in `constants.go`.
