

dep:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest	
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	cp ~/go/bin/protoc-gen-go /usr/local/bin/
	cp ~/go/bin/protoc-gen-go-grpc /usr/local/bin/



proto:
	rm -rf pkg/protocol/proto/pb/*
	protoc --proto_path=./pkg/protocol/proto --go_out=. --go-grpc_out=. ./pkg/protocol/proto/*.proto

build:
	rm -rf bin/*
	go build -o bin/gateway cmd/connection/gateway/main.go
	go build -o bin/state cmd/connection/state/main.go
	go build -o bin/logicHttp cmd/connection/logicHttp/main.go
	go build -o bin/ipconfig cmd/connection/ipconfig/main.go
	go build -o bin/logic cmd/logic/main.go
	mkdir log -p
	touch log/gateway.log
	touch log/state.log
	touch log/logicHttp.log
	touch log/ipconfig.log
	touch log/logic.log

gateway:
	nohup ./bin/gateway > log/gateway.log &

state:
	nohup ./bin/state > log/state.log &

logicHttp:
	nohup ./bin/logicHttp > log/logicHttp.log &
	
ipconfig:
	nohup ./bin/ipconfig > log/ipconfig.log &

logic:
	nohup ./bin/logic > log/logic.log &

run: gateway state logicHttp ipconfig logic

all:  proto build
