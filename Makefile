BIN_NAME=opensergo-control-plane
BUILD_DIR=build
SRC_MAIN=pkg/main/main.go
.DEFAULT_GOAL=build

build:
	go build -o ${BUILD_DIR}/${BIN_NAME} ${SRC_MAIN}

run:
	go run ${SRC_MAIN}


clean:
	go clean
	rm -rf ${BUILD_DIR}

test:
	go test