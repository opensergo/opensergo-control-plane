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

FROM golang:1.16.6-alpine3.14 AS builder

LABEL maintainer="OpenSergo Maintainers" \
      website="https://www.opensergo.io" \
      source="https://github.com/opensergo/opensergo-control-plane"

# arg goproxy
## if you are in China, in your build command, you can set arg `GOPROXY` like `https://goproxy.cn,direct`
ARG GOPROXY=https://proxy.golang.org,direct

# define work dir
WORKDIR /opensergo/src

# copy source-code
COPY . /opensergo/src

# build opensergo
RUN cd /opensergo/src/pkg/main && \
    go build -o /opensergo/dist/bin/opensergo && \
    chmod +x /opensergo/dist/bin/opensergo


FROM alpine:3.14 AS runtime

# define work dir
WORKDIR /opensergo

# copy binary file
COPY --from=builder /opensergo/dist .

# expose default port
EXPOSE 10246

# start openser-control-plane
ENTRYPOINT ["./bin/opensergo"]