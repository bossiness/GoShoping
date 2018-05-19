package center

import (
	"btdxcx.com/walker/apis/common/server"
	"github.com/emicklei/go-restful"

	"os"
)

// APIServer 实现 APIServer
type APIServer struct {
	server.CommonAPIServer
}

// NewAPIServer new api serve
func NewAPIServer() *APIServer {
	container := restful.NewContainer()
	container.Router(restful.CurlyRouter{})
	apis := APIServer{}
	apis.Host = "localhost"
	apis.Port = ":3001"
	apis.Container = container
	apis.SwaggerFilePath = "/tmp/swagger-ui/dist"
	domainHost := os.Getenv("DOMAINHOST")
	if len(domainHost) != 0 {
		apis.Host = domainHost
	}
	return &apis
}
