#ifndef __CHAT_H__
#define __CHAT_H__

#include "defs.h"
#include <string>
#include <map>
#include <mutex>

#include <grpc++/grpc++.h>
#include <grpc++/server.h>
#include <grpc++/server_context.h>

#include "pb/chat/chat.grpc.pb.h"
#include "pb/auth/auth.grpc.pb.h"

#include "user.hpp"

struct AuthStatus {
  const grpc::Status status;
  const grpc::string_ref token;
  const grpc::string_ref username;

  AuthStatus(grpc::Status status,
             grpc::string_ref token,
             grpc::string_ref username);
};

class ChatService final : public chat::Chat::Service {
    std::mutex mu_;
    std::map<std::string, User*> users_;
    u_ptr<auth::Auth::Stub> authStub_;

    AuthStatus authenticateContext(grpc::ServerContext *context);
  public:
    // chan is a channel to the authentication service
    explicit ChatService(s_ptr<grpc::Channel> authChan);
    ~ChatService() = default;

    virtual grpc::Status Join(grpc::ServerContext* context,
                              const chat::JoinReq* request,
                              chat::Receipt* response) override;
    virtual grpc::Status Send(grpc::ServerContext* context,
                              const chat::Payload* request,
                              chat::Receipt* response) override;
    virtual grpc::Status Poll(grpc::ServerContext* context,
                              const chat::PollReq* request,
                              chat::PollUpdate* response) override;
};

#endif
