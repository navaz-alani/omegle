#include "chat.hpp"

#include <iostream>

#include <grpcpp/server.h>
#include <grpcpp/server_builder.h>
#include <grpcpp/server_context.h>
#include <grpcpp/security/server_credentials.h>
#include <grpcpp/security/credentials.h>
#include <grpcpp/channel.h>
#include <grpcpp/create_channel.h>

int main(void) {
  // instantiate chatService with a channel to authentication service
  ChatService *service = new ChatService(grpc::CreateChannel("0.0.0.0:4002",
                                         grpc::InsecureChannelCredentials()));
  std::string serviceAddr = "0.0.0.0:4001";

  grpc::ServerBuilder builder;
  builder.AddListeningPort(serviceAddr, grpc::InsecureServerCredentials());
  builder.RegisterService(service);

  u_ptr<grpc::Server> server = builder.BuildAndStart();
  std::cout << "Listening on " << serviceAddr << std::endl;
  server->Wait();
}
