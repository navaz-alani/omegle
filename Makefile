GXX=go
BIN=omegle
GXX_SRC=$(wildcard **/*.go)
PROTO=protoc
PROTO_DEF=./pb/omegle.proto
PROTO_TS_OUT=./web/pb
PROTO_GO_OUT=./pb/go

.PHONY: exec grpc-go grpc-ts clean clean-grpc tidy

exec: $(BIN)
	./$(BIN)

$(BIN): $(GXX_SRC)
	go get
	$(GXX) build -o $(BIN) .

grpc-go: $(PROTO_DEF)
	mkdir -p $(PROTO_GO_OUT)
	$(PROTO) --go_out=$(PROTO_GO_OUT) --go_opt=paths=source_relative \
           --go-grpc_out=$(PROTO_GO_OUT) --go-grpc_opt=paths=source_relative \
		       $(PROTO_DEF)

grpc-ts: $(PROTO_DEF)
	mkdir -p $(PROTO_TS_OUT)
	protoc --js_out=import_style=commonjs,binary:$(PROTO_TS_OUT) $(PROTO_DEF)
	protoc --grpc-web_out=import_style=typescript,mode=grpcwebtext:$(PROTO_TS_OUT) $(PROTO_DEF)

clean:
	rm -rf $(BIN)

clean-grpc:
	rm -rf $(PROTO_TS_OUT) $(PROTO_GO_OUT)

tidy:
	go mod tidy
