package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/navaz-alani/oryx/auth"
	authpb "github.com/navaz-alani/oryx/pb/go/pb/auth"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 10000, "The server port")
  secret = flag.String("secret", "", "Service authentication secret")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
  if *secret == "" {
		log.Fatalf("Authentication secret not provided\n")
  }
  addr := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
  var srv authpb.AuthServer
	srv, err = auth.NewAuthService(*secret)
	if err != nil {
		log.Fatalf("failed to instantiate auth service: %v\n", err)
	}
	authpb.RegisterAuthServer(grpcServer, srv)
  log.Printf("Listening on %s", addr)
	grpcServer.Serve(listener)
}
