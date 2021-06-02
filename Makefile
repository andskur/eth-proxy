COMMON_PATH	?= $(shell pwd)

clean:
	rm -f $(APP)

tidy:
	go mod tidy

update:
	go get -u ./...

protoc:
	protoc -I $(COMMON_PATH)/proto $(COMMON_PATH)/proto/*.proto  --go_out=plugins=grpc:$(COMMON_PATH)/proto
