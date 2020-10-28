package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/navaz-alani/omegle/pb/go/pb"
)

const UsernameLen = 5

type OmegleServer struct {
  pb.UnimplementedOmegleServer

  mu        sync.RWMutex
  unmatched map[string]bool
  available map[string]*User
  engaged   map[string]*User
}

func NewServer() *OmegleServer {
  return &OmegleServer{
    mu: sync.RWMutex{},
    unmatched: make(map[string]bool),
    available: make(map[string]*User),
    engaged: make(map[string]*User),
  }
}

func (o *OmegleServer) match(usrA, usrB string) error {
  o.mu.Lock()
  defer o.mu.Unlock()

  if uA, ok := o.available[usrA]; !ok {
    return fmt.Errorf("%s not available", usrA)
  } else if uB, ok := o.available[usrB]; !ok {
    return fmt.Errorf("%s not available", usrB)
  } else {
    // match two users
    uA.SetMatch(usrB)
    uB.SetMatch(usrA)
    // move them to appropriate maps
    o.engaged[usrA] = o.available[usrA]
    delete(o.available, usrA)
    delete(o.unmatched, usrA)
    o.engaged[usrB] = o.available[usrB]
    delete(o.available, usrB)
    delete(o.unmatched, usrB)
  }

  return nil
}

func (o *OmegleServer) unmatch(usrA, usrB string) error {
  o.mu.Lock()
  defer o.mu.Unlock()

  if uA, ok := o.engaged[usrA]; !ok {
    return fmt.Errorf("%s not engaged", usrA)
  } else if uB, ok := o.engaged[usrB]; !ok {
    return fmt.Errorf("%s not engaged", usrB)
  } else {
    // unmatch two users
    uA.SetMatch("")
    uB.SetMatch("")
    // move them to appropriate maps
    o.available[usrA] = o.engaged[usrA]
    o.unmatched[usrA] = true
    delete(o.engaged, usrA)
    o.available[usrB] = o.engaged[usrB]
    o.unmatched[usrB] = true
    delete(o.engaged, usrB)
  }

  return nil
}

func (o *OmegleServer) usernameTaken(uname string) bool {
  o.mu.RLock()
  defer o.mu.RUnlock()

  if _, ok := o.unmatched[uname]; ok {
    return true
  } else if _, ok := o.engaged[uname]; ok {
    return true
  }
  return false
}

func (o *OmegleServer) Join(req *pb.JoinReq, stream pb.Omegle_JoinServer) error {
  o.mu.Lock()
  defer o.mu.Unlock()

  username := req.Username
  if req.Username == "" {
    username = randStr(UsernameLen, AlphaLU)
  } else if o.usernameTaken(req.Username) {
    return fmt.Errorf("Username taken.")
  }
  u := NewUser(username, stream)
  o.unmatched[req.Username] = true
  o.available[req.Username] = u
  return nil
}

func (o *OmegleServer) Chat(ctx context.Context, req *pb.ChatReq) (*pb.ChatResp, error) {
  o.mu.RLock()
  defer o.mu.RUnlock()

  match := req.Requested
  if match == "" {
    // random match
    for m := range o.unmatched {
      match = m
      break
    }
  }

  if err := o.match(req.Requestor, match); err != nil {
    return &pb.ChatResp{
      Status: pb.ChatResp_UNAVAILABLE,
      Msg: err.Error(),
    }, nil
  }

  return &pb.ChatResp{
    Status: pb.ChatResp_OK,
    Msg: "",
  }, nil
}

func (o *OmegleServer) Send(ctx context.Context, p *pb.Payload) (*pb.SendReceipt, error) {
  o.mu.Lock()
  defer o.mu.Unlock()

  if to, ok := o.engaged[p.To]; !ok || to.CurrentMatch() != p.From {
    return &pb.SendReceipt{
      Status: pb.SendReceipt_ERROR,
      Msg: "Not matched with recipient.",
    }, nil
  } else {
    if err := to.Send(p); err != nil {
      return &pb.SendReceipt{
        Status: pb.SendReceipt_ERROR,
        Msg: err.Error(),
      }, nil
    }
    return &pb.SendReceipt{
      Status: pb.SendReceipt_OK,
    }, nil
  }
}
