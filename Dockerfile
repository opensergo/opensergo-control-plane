# Build the manager binary
FROM golang:1.17-alpine3.14 as  builder

ENV GOPROXY=https://proxy.golang.com.cn,direct

WORKDIR /workspace

# Copy the go source
COPY / ./
#COPY main.go main.go
#COPY api/ api/
#COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM alpine:3.14
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
