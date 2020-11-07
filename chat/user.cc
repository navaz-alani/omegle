#include "user.hpp"
#include "pb/auth/auth.grpc.pb.h"

User::User(const std::string &username)
  : username_{ username } {}

const std::string User::username() const { return username_; }

void User::receive(chat::Payload* p) {
  mu_.lock();
  incoming_.push(p);
  mu_.unlock();
}

void User::poll(chat::PollUpdate *p) {
  // load the update with incoming messages
  mu_.lock();
  chat::Payload *newPayload;
  while (!incoming_.empty()) {
    newPayload = p->add_incoming();
    newPayload = incoming_.front();
    incoming_.pop();
  }
  mu_.unlock();
}
