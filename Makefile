SYSTEM=$(shell uname | cut -f 1 -d_)

GXX=go
CXX=g++
CXXFLAGS=-std=c++14 -Wall
CPPFLAGS=`pkg-config --cflags protobuf grpc` \
				 -I${CURDIR}/pb/cpp # grpc & protobuf generated code
CXX_LD_FLAGS=-L/usr/local/lib `pkg-config --libs protobuf grpc++`\
					-pthread\
          -ldl
ifeq ($(SYSTEM),Darwin)
	CXX_LD_FLAGS+=-framework CoreFoundation \
          			-lgrpc++_reflection
else
	CXX_LD_FLAGS+=-Wl,--no-as-needed -lgrpc++_reflection -Wl,--as-needed
endif

GXX_SRC=$(wildcard *.go pb/go/pb/*.go)
CXX_SRC=$(wildcard chat/*.cc pb/cpp/pb/auth/*.cc pb/cpp/pb/chat/*.cc)
CXX_OBJS=${CXX_SRC:.cc=.o}

PROTOC=protoc
PROTO_DEF=$(wildcard ./pb/*/*.proto)
PROTO_TS_OUT=./web/pb
PROTO_GO_OUT=./pb/go
PROTO_CPP_OUT=./pb/cpp

CHAT_SRVC=chat_srvc
AUTH_SRVC=auth_srvc

.PHONY: proto \
				grpc-cpp grpc-go grpc-ts \
				clean clean-cpp clean-grpc tidy \
				exec_chat_srvc exec_auth_srvc

exec_chat_srvc: $(CHAT_SRVC)
	./$(CHAT_SRVC)

exec_auth_srvc: $(AUTH_SRVC)
	./$(AUTH_SRVC)

$(AUTH_SRVC): $(GXX_SRC)
	$(GXX) get
	$(GXX) build -o $@ .

$(CHAT_SRVC): $(CXX_OBJS)
	$(CXX) $(CXXFLAGS) $(CXX_LD_FLAGS) $(CXX_OBJS) -o $@

%.o: %.cc
	$(CXX) $(CXXFLAGS) $(CPPFLAGS) -c -o $@ $<

proto: grpc-cpp grpc-go grpc-ts

grpc-cpp:
	mkdir -p $(PROTO_CPP_OUT)
	$(PROTOC) --cpp_out=$(PROTO_CPP_OUT) \
				 --grpc_out=$(PROTO_CPP_OUT) \
		     --plugin=$(PROTOC)-gen-grpc=`which grpc_cpp_plugin` \
				 $(PROTO_DEF)

grpc-go: $(PROTO_DEF)
	mkdir -p $(PROTO_GO_OUT)
	$(PROTOC) --go_out=$(PROTO_GO_OUT) --go_opt=paths=source_relative \
         --go-grpc_out=$(PROTO_GO_OUT) --go-grpc_opt=paths=source_relative \
		     $(PROTO_DEF)

grpc-ts: $(PROTO_DEF)
	mkdir -p $(PROTO_TS_OUT)
	$(PROTOC) --js_out=import_style=commonjs,binary:$(PROTO_TS_OUT) $(PROTO_DEF)
	$(PROTOC) --grpc-web_out=import_style=typescript,mode=grpcwebtext:$(PROTO_TS_OUT) $(PROTO_DEF)

clean:
	rm -rf $(BIN)

clean-grpc:
	rm -rf $(PROTO_TS_OUT) $(PROTO_GO_OUT)

clean-cpp:
	rm $(CXX_OBJS)

tidy:
	go mod tidy
