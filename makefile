

dep:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest	
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	cp ~/go/bin/protoc-gen-go /usr/local/bin/
	cp ~/go/bin/protoc-gen-go-grpc /usr/local/bin/



proto:
	rm -rf pkg/protocol/proto/pb/*
	protoc --proto_path=./pkg/protocol/proto --go_out=. --go-grpc_out=. ./pkg/protocol/proto/*.proto

