PROTO_DIR=./protos
SRC_DIR=$(PROTO_DIR)/src

PROTOC=protoc

GO_OUT=--go_out
GO_GRPC_OUT=--go-grpc_out

PROTO_FILES=$(shell find $(SRC_DIR) -name "*.proto")

generate-proto:
	@for file in $(PROTO_FILES); do \
		BASE_NAME=$(basename $$file .proto); \
		OUTPUT_DIR=$(PROTO_DIR)/$$BASE_NAME; \
		mkdir -p $$OUTPUT_DIR; \
		$(PROTOC) $$file $(GO_OUT)=$$OUTPUT_DIR $(GO_GRPC_OUT)=$$OUTPUT_DIR; \
	done

clean-proto:
	@find $(PROTO_DIR) -mindepth 1 \
		-not -path "$(SRC_DIR)" \
		-not -name "go.mod" \
		-not -name "go.sum" \
		-not -name "*.proto" \
		-exec rm -rf {} +
