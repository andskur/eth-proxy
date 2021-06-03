COMMON_PATH	?= $(shell pwd)

tidy:
	go mod tidy

update:
	go get -u ./...

generate-api:
	swagger generate server -A proxy-api --server-package server -f ./api/swagger.yaml --exclude-main --keep-spec-order --flag-strategy pflag --target ./api --spec ./api/swagger.yaml

generate-api-docker:
	docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=${HOME}/Code/Go:/go -v ${HOME}:${HOME} -w $(shell pwd) quay.io/goswagger/swagger:v0.25.0 generate server -A proxy-api --server-package server -f ./api/swagger.yaml --exclude-main --keep-spec-order --flag-strategy pflag --target ./api --spec ./api/swagger.yaml

protoc:
	protoc -I $(COMMON_PATH)/proto $(COMMON_PATH)/proto/*.proto  --go_out=plugins=grpc:$(COMMON_PATH)/proto

run_service:
	go run proxy-service/cmd/proxy-service.go serve

run_api:
	API_HOST=127.0.0.1 go run api/cmd/api.go serve