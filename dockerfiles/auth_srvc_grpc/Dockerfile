FROM golang:latest

ARG AUTH_SECRET
ENV AUTH_SECRET ${AUTH_SECRET}

COPY ./ /go/src/auth_srvc_grpc

RUN make --directory=src/auth_srvc_grpc exec_auth_srvc
