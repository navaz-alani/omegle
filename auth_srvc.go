package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/navaz-alani/omegle/auth"
	authpb "github.com/navaz-alani/omegle/pb/go/pb/auth"
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
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	srv, err := auth.NewAuthService(*secret)
	if err != nil {
		log.Fatalf("failed to instantiate auth service: %v\n", err)
	}
	authpb.RegisterAuthServer(grpcServer, srv)
	grpcServer.Serve(listener)
}
