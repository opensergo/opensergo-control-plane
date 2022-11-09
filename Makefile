# Copyright 2022, OpenSergo Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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