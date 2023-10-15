
.PHONY: run
run:
	go run cmd/main.go

# ----------------------------------------------------------------


.PHONY:generate
gen: 
	protoc -I=. \
	--go_out=. \
	--go_opt=paths=source_relative my_api.proto  \
 
# protoc -I ./proto \
# --go_out ./proto --go_opt paths=source_relative \
# --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
#  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
# ./proto/my_api.proto

# protoc --proto_path=.\
# --go_out=. \
# --go_opt=source-relative=. my_api.proto 

.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
