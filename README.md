# Inabox

`inaboxctl` is service that allows VanillaVerse developers to deploy VanillaVerse apps to a remote machine
and manage services running on that machine, by wrapping the `docker-compose` commandset and using `rsync`.

This app is written as a successor to the old Inabox design.

## Installation

To install this service to your machine, do:
```
go install github.com/vanillaverse/inaboxctl
```
Also, make sure that `GOBIN` or (`GOPATH/bin`) are in your `PATH` variable!
