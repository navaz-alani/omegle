#ifndef __USER_H__
#define __USER_H__

#include <string>
#include <queue>
#include <mutex>

#include "pb/chat/chat.grpc.pb.h"

class User {
    std::mutex mu_;
    const std::string username_;
    std::queue<chat::Payload*> incoming_;

  public:
    User(const std::string &username);
    ~User() = default;
    std::string username() const;

    // Users receive payloads, which are queued.
    void receive(chat::Payload* p);
    // Users constantly poll the service for their message updates. Updates are
    // put in the object at p.
    void poll(chat::PollUpdate *p);
};

#endif
