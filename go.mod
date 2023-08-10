module github.com/opensergo/opensergo-control-plane

go 1.14

require (
	github.com/alibaba/sentinel-golang v1.0.3
	github.com/envoyproxy/go-control-plane v0.10.3-0.20221109183938-2935a23e638f
	github.com/envoyproxy/protoc-gen-validate v0.6.7
	github.com/go-logr/logr v0.4.0
	github.com/golang/protobuf v1.5.2
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	go.uber.org/atomic v1.7.0
	google.golang.org/genproto v0.0.0-20220329172620-7be39ac1afc7
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	sigs.k8s.io/controller-runtime v0.9.7
)
