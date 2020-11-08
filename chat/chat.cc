#include "chat.hpp"

#include <map>

#include <grpc++/client_context.h>

using grpc::Status;
using grpc::StatusCode;
using grpc::ClientContext;

AuthStatus::AuthStatus(Status status, grpc::string_ref token, grpc::string_ref username)
    : status{ status }, token{ token }, username{ username } {}

ChatService::ChatService(s_ptr<grpc::Channel> authChan)
  : authStub_{ auth::Auth::NewStub(authChan) } {}

AuthStatus ChatService::authenticateContext(grpc::ServerContext *context) {
  auto meta = context->client_metadata();
  auto token = meta.find("_jwt_");
  auto username = meta.find("_username_");

  if (token == meta.end() || username == meta.end())
    return AuthStatus(grpc::Status(StatusCode::UNAUTHENTICATED,
                                   "incomplete credentials"), "", "");

  // verify client context using rpc call to auth service
  ClientContext ctx;
  auth::Cert cert;
  auth::CertStatus certStatus;
  cert.set_jwt(token->second.data()); cert.set_username(username->second.data());
  Status status = authStub_->VerifCert(&ctx, cert, &certStatus);

  if (!status.ok())
    return AuthStatus(grpc::Status(StatusCode::INTERNAL,
                                   "authentication rpc failure"), "", "");
  else if (certStatus.status() != auth::CertStatus_Status_VALID)
    return AuthStatus(grpc::Status(StatusCode::UNAUTHENTICATED,
                                   "invalid credentials"), "", "");
  return AuthStatus(Status::OK, token->second, username->second);
}

grpc::Status ChatService::Join(grpc::ServerContext* context,
                               const chat::JoinReq* request,
                               chat::Receipt* response) {
  const AuthStatus as = authenticateContext(context);
  if (as.status.error_code() != StatusCode::OK) return as.status;
  const std::string username = as.username.data();

  // add user to service
  mu_.lock();
  users_[username] = new User(username);
  mu_.unlock();
  response->set_status(chat::Receipt_Status::Receipt_Status_OK);
  return grpc::Status::OK;
}

grpc::Status ChatService::Send(grpc::ServerContext* context,
                               const chat::Payload* p,
                               chat::Receipt* response) {
  const AuthStatus as = authenticateContext(context);
  if (as.status.error_code() != StatusCode::OK) return as.status;

  // copy payload and add it to user's message queue
  chat::Payload *payload = new chat::Payload(*p);
  mu_.lock();
  users_[as.username.data()]->receive(payload);
  mu_.unlock();
  response->set_status(chat::Receipt_Status::Receipt_Status_OK);
  return grpc::Status::OK;
}

grpc::Status ChatService::Poll(grpc::ServerContext* context,
                               const chat::PollReq* request,
                               chat::PollUpdate* pollUpdate) {
  const AuthStatus as = authenticateContext(context);
  if (as.status.error_code() != StatusCode::OK) return as.status;

  mu_.lock();
  users_[as.username.data()]->poll(pollUpdate);
  mu_.unlock();
  return grpc::Status::OK;
}
