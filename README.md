# go-net
Golang net wrapper for HTTP_PROXY and NAT with UPnP.

This library is wrapper for golang net library.

This support
* Dial and DialContext support connection with HTTP Proxy (Need set environment variable "HTTPS_PROXY").
* ListenTCP set the PortMapping in UPnP Router.

TODO
* When using UPnP Router, the global address always show 8.8.8.8.