APP_NAME = rdpms25-go-rpc-service
DOCKER_REPO = infinityshiv/rdpms25-go-rpc-service
RELEASE_VERSION = occ3.0.0

run:
	@go mod tidy
	@go run main.go

build:
	go build -o bin/$(APP_NAME) main.go

window:
	echo "Compiling for windows"
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/window/$(APP_NAME).exe main.go
	echo "Output stored in ./bin"

linux:
	echo "Compiling for linux"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/linux/$(APP_NAME) main.go
	echo "Output stored in ./bin"

alpine:
	echo "Compiling for docker alpine"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.version=$(RELEASE_VERSION) -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -a -installsuffix cgo -o bin/alpine/$(APP_NAME) main.go
	echo "Output stored in ./bin"

docker:
	$(MAKE) alpine
	docker build -t $(APP_NAME) --platform linux/amd64 --build-arg APP_NAME=$(APP_NAME) --build-arg RELEASE_VERSION=$(RELEASE_VERSION) .

hub:
	docker tag $(APP_NAME) $(DOCKER_REPO):$(RELEASE_VERSION)
	docker push $(DOCKER_REPO):$(RELEASE_VERSION)

orm:
	sqlboiler psql --add-global-variants

grpc:
	protoc --go_out=. --go-grpc_out=. proto/client.proto

swag:
	swag init -g ./main.go -o ./docs --parseDependency --parseInternal

compose:
	docker compose up --build -d

down:
	docker compose down
