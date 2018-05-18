package server

import (
	"net/http"

	"github.com/micro/go-log"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
)

// APIServer api server
type APIServer interface {
	GetHost() string
	GetPort() string
	GetContainer() *restful.Container
	GetSwaggerFilePath() string
	Start()
}

// CommonAPIServer 实现 APIServer
type CommonAPIServer struct {
	Host            string
	Port            string
	Container       *restful.Container
	SwaggerFilePath string
}

// NewCommonAPIServer new api serve
func NewCommonAPIServer() *CommonAPIServer {
	container := restful.NewContainer()
	container.Router(restful.CurlyRouter{})
	return &CommonAPIServer{
		Host:            "localhost",
		Port:            ":3001",
		Container:       container,
		SwaggerFilePath: "/tmp/swagger-ui/dist",
	}
}

// GetContainer GET restful.Container
func (as *CommonAPIServer) GetContainer() *restful.Container {
	return as.Container
}

// GetHost GET
func (as *CommonAPIServer) GetHost() string {
	return as.Host
}

// GetPort GET
func (as *CommonAPIServer) GetPort() string {
	return as.Port
}

// GetSwaggerFilePath GET
func (as *CommonAPIServer) GetSwaggerFilePath() string {
	return as.SwaggerFilePath
}

// Start api server
func (as *CommonAPIServer) Start() {
	as.configSwagger()

	restful.TraceLogger(as)

	log.Log("start listening on " + as.GetHost() + as.GetPort())
	server := &http.Server{Addr: as.GetPort(), Handler: as.GetContainer()}
	log.Log(server.ListenAndServe())
}

func (as *CommonAPIServer) configSwagger() {
	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open <your_app_id>.appspot.com/apidocs and enter http://<your_app_id>.appspot.com/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    as.GetContainer().RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: as.getGaeURL(),
		ApiPath:        "/apidocs.json",
		SwaggerPath:    "/apidocs/",
		Info: swagger.Info{
			Title:       "Shop APIs",
			Description: "Shop APIs",
		},
		ApiVersion:      "v1.0.0",
		SwaggerFilePath: as.GetSwaggerFilePath()}
	swagger.RegisterSwaggerService(config, as.GetContainer())

}

func (as *CommonAPIServer) getGaeURL() string {
	return "http://" + as.GetHost() + as.GetPort()
}

// Print Trace Logger
func (as CommonAPIServer) Print(v ...interface{}) {
	log.Log(v)
}

// Printf Trace Logger
func (as CommonAPIServer) Printf(format string, v ...interface{}) {
	log.Logf(format, v)
}
