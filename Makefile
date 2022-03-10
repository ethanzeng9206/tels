gen:
	protoc -I ./proto --go_out=./pb --go-grpc_out=./pb ./proto/*.proto

