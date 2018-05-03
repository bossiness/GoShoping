package gridfsweb

import (
	"github.com/micro/cli"
	"github.com/micro/go-web"
	"github.com/micro/go-log"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"

	"time"
)

var (
	// DBUrl mongodb URL
	DBUrl = "localhost:27017"
)

// Commands add auth api command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "gridfs",
			Usage:  "Run image web",
			Action: aweb,
		},
	}
}

func aweb(ctx *cli.Context) {

	service := web.NewService(
		web.Name("go.micro.web.gridfs"),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	handler := new(Handler)

	r := mux.NewRouter()
	r.Path("/gridfs")
	r.HandleFunc("/", handler.HomeHandler).Methods("GET")
	r.HandleFunc("/", handler.UploadHandler).Methods("POST")
	r.HandleFunc("/{imgid}", handler.DownloadHandler).Methods("GET")

	service.Handle("/", r)

	handler.Init()
	defer handler.Deinit()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

// Handler DB
type Handler struct {
	session *mgo.Session
	db *mgo.Database
}

// Init 数据库初始化
func (h *Handler) Init() error {
	session, err := mgo.Dial(DBUrl)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	h.session = session
	h.db = h.session.DB("image_fiels")
	return nil
}

// Deinit 资源释放
func (h *Handler) Deinit() {
	if h.session != nil {
		h.session.Close()
	}
}