Stupidly simple command line tool to resolve and ping hostnames.

I just use the binary to troubleshoot connectivity issues in scratch/alpine containers that don't have ping/dig/curl/wget.
Just copy the binary to the container, and problem solved.

The binary has to be built with `CGO_ENABLED=0` to avoid problems with alpine-based images.

Example of use:

```
pingish www.google.es
```

# Server

You can start a server that will accept requests on `/ping?host=<target>`, using the `--server` flag. e.g.

```
pingish --server
```

# TROUBLESHOOTING

To run it as a normal user in ubuntu, you might need to configure your host first: `sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"` 

And/or set the capabilities of the binary: `sudo setcap cap_net_raw=ep pingish`


