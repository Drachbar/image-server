APP_NAME=image-server
BUILD_DIR=build

build: build-binary

build-binary:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME) main.go

clean:
	rm -rf $(BUILD_DIR)

run: build-binary
	./$(BUILD_DIR)/$(APP_NAME) -apikey=din-nyckel -dir=./bilder -baseurl=http://localhost:8080/images
