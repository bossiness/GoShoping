package server

import (
	"net/http"

	"github.com/micro/go-log"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
)

// APIServer api server
type APIServer struct {
	host            string
	port            string
	container       *restful.Container
	swaggerFilePath string
}

// NewAPIServer new api serve
func NewAPIServer() *APIServer {
	container := restful.NewContainer()
	container.Router(restful.CurlyRouter{})
	return &APIServer{
		host:            "localhost",
		port:            ":3001",
		container:       container,
		swaggerFilePath: "/tmp/swagger-ui/dist",
	}
}

// Container GET restful.Container
func (as *APIServer) Container() *restful.Container {
	return as.container
}

// Start api server
func (as *APIServer) Start() {
	as.configSwagger()

	restful.TraceLogger(as)

	log.Log("start listening on " + as.host + as.port)
	server := &http.Server{Addr: as.port, Handler: as.container}
	log.Log(server.ListenAndServe())
}

func (as *APIServer) configSwagger() {
	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open <your_app_id>.appspot.com/apidocs and enter http://<your_app_id>.appspot.com/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    as.container.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: as.getGaeURL(),
		ApiPath:        "/apidocs.json",
		SwaggerPath:    "/apidocs/",
		Info: swagger.Info{
			Title:       "Shop APIs",
			Description: "Shop APIs",
		},
		ApiVersion:      "v1.0.0",
		SwaggerFilePath: as.swaggerFilePath}
	swagger.RegisterSwaggerService(config, as.container)

}

func (as *APIServer) getGaeURL() string {
	return "http://" + as.host + as.port
}

// Print Trace Logger
func (as APIServer) Print(v ...interface{}) {
	log.Log(v)
}

// Printf Trace Logger
func (as APIServer) Printf(format string, v ...interface{}) {
	log.Logf(format, v)
}
