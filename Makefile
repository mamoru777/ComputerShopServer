$(info $(SHELL))
USER_PROTO_OUT_DIR = pkg
USER_API_PATH = gRPC

.PHONY: gen-user

gen-user:

	protoc \
		-I ${USER_API_PATH} \
		--go_out=$(USER_PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(USER_PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		./${USER_API_PATH}/*.proto