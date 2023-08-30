package main

import (
	"log"

	stream_plugin "github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin/stream"

	"github.com/opensergo/opensergo-control-plane/pkg/plugin/pl"
	"github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin"
)

//nolint:gosimple
func main() {
	pluginServer := pl.NewPluginServer()
	//defer func() {
	//	if err := pluginServer.RunShutdownFuncs(); err != nil {
	//		log.Fatalln("Error:", err.Error())
	//	}
	//	log.Println("Server shutdown")
	//	return
	//}()
	err := pluginServer.InitPlugin()
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}

	client, err := pluginServer.GetPluginClient(builtin.StreamServicePluginName)
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}
	raw, ok := client.(stream_plugin.Stream)
	if !ok {
		log.Fatalln("Error: can't convert rpc plugin to normal wrapper")
	}

	sa := &say{}
	greet, err := raw.Greeter("这是一个前缀", sa)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
	log.Println("Greeter:", greet)

	for {
		select {
		case <-pluginServer.ShutdownCh:
			pluginServer.ContextCancel()
			if err := pluginServer.RunShutdownFuncs(); err != nil {
				log.Fatalln("Error:", err.Error())
			}
			log.Println("Server shutdown")
			return
		}
	}

}

type say struct {
}

func (s *say) Say(ss string) string {
	return ss + "这是一个后缀v2"
}
