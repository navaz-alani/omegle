package main

import (
	"sync"

	"github.com/navaz-alani/omegle/pb/go/pb"
)

type User struct {
  *pb.User

  mu sync.RWMutex
  stream pb.Omegle_JoinServer
  outgoing chan *pb.Payload
  currMatch string
}

func NewUser(username string, stream pb.Omegle_JoinServer) *User {
  return &User{
    User: &pb.User{
      Username: username,
    },
    stream: stream,
    outgoing: make(chan *pb.Payload),
    mu: sync.RWMutex{},
    currMatch: "",
  }
}

func (u *User) SetMatch(m string) {
  u.mu.Lock()
  defer u.mu.Unlock()
  u.currMatch = m
}

func (u *User) CurrentMatch() string {
  u.mu.RLock()
  defer u.mu.RUnlock()
  return u.currMatch
}

func (u *User) Send(p *pb.Payload) error {
  u.mu.RLock()
  defer u.mu.RUnlock()
  return u.stream.Send(p)
}

func (u *User) Outgoing() <-chan *pb.Payload {
  u.mu.RLock()
  defer u.mu.RUnlock()
  return u.outgoing
}
