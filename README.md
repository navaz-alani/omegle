# Oryx

This is a service similar to Omegle, implemented using Go & C++ (for the
backend) and React & NextJS (for the frontend) with gRPC.

## Features

- [x] Conersations with (possibly random) people who are also online
- [x] Text messaging
- [ ] Multimedia & file sharing
- [ ] Voice calling
- [ ] Video calling

## Build

There are two gRPC services (`auth_srvc_grpc` and `chat_srvc_grpc`) and one
server for the frontend (`frontend`). These services can all be seen in the
`docker-compose.yml` file and have to be build independently.

The `auth_srvc_grpc` is implemented in Golang and is quite simple to build and
run. The command `make auth_srvc_grpc` will produce a binary with the service
name.

The `chat_srvc_grpc` is implemented in C++ and depending on the chosen build
method, the build process can be a bit involved. There are two ways to build
this service. If you already have `grpc` and `protobuf` installed on your
system and configured with `pkgconfig`, then you may simply run `make
chat_srvc_grpc` and this will produce the service binary (if there are no
compilation/linking errors). The second way, which may be easier, is to use
Bazel: the command `bazel build //chat:chat_srvc_grpc` will produce the service
binary. This second method is quite time consuming however since `grpc` and
`protobuf` have to be built from source.

The frontend server is easy to run with the command `cd web && yarn dev`.

More information about builing & deployment can be obtained from the following:
`docker-compose.yml` (network setup), `nginx.conf` (for deployment with the
services running in docker-compose) and `dockerfiles/*/Dockerfile` (for building
services).
