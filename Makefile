PROTOC=protoc

.PHONY: proto

proto: proto/*.proto
	$(PROTOC) --go_out=proto/ proto/*.proto
