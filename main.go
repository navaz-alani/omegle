package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/navaz-alani/omegle/pb/go/pb"
	"google.golang.org/grpc"
)

var (
  port = flag.Int("port", 10000, "The server port")
)

func main() {
  flag.Parse()
  rand.Seed(time.Now().UnixNano())
  listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }
  log.Println("listening on port ", *port)
  var opts []grpc.ServerOption
  grpcServer := grpc.NewServer(opts...)
  // TODO: implement newServer() in package main
	pb.RegisterOmegleServer(grpcServer, NewServer())
	grpcServer.Serve(listener)
}
