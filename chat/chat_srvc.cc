#include "chat.hpp"

#include <iostream>

#include <grpc++/server.h>
#include <grpc++/server_builder.h>
#include <grpc++/server_context.h>
#include <grpc++/security/server_credentials.h>
#include <grpc++/security/credentials.h>
#include <grpc++/channel.h>
#include <grpc++/create_channel.h>

int main(void) {
  // instantiate chatService with a channel to authentication service
  ChatService *service = new ChatService(grpc::CreateChannel("0.0.0.0:10000",
                                         grpc::InsecureChannelCredentials()));
  std::string serviceAddr = "0.0.0.0:10001";

  grpc::ServerBuilder builder;
  builder.AddListeningPort(serviceAddr, grpc::InsecureServerCredentials());
  builder.RegisterService(service);

  u_ptr<grpc::Server> server = builder.BuildAndStart();
  std::cout << "Listening on " << serviceAddr << std::endl;
  server->Wait();
}
