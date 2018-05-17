package image

import (
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-log"
	"btdxcx.com/walker/apis/common/errors"
)

// Upload image
func Upload(req *restful.Request, resp *restful.Response) {

	imgid := NewImageID()

	r := req.Request
	w := resp.ResponseWriter

	//上传参数为uploadfile
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		log.Log(err)
		w.Write([]byte("Error:Upload Error."))
		return
	}
	defer file.Close()

	//检测文件类型
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Log(err)
		w.Write([]byte("Error:Upload Error."))
		return
	}
	filetype := http.DetectContentType(buff)
	if (filetype != "image/jpeg") && (filetype != "image/png") {
		w.Write([]byte("Error:Not JPEG."))
		return
	}

	//回绕文件指针
	log.Log(filetype)
	if _, err = file.Seek(0, 0); err != nil {
		log.Log(err)
	}

	if err = BuildDir(imgid); err != nil {
		log.Log(err)
	}

	f, err := os.OpenFile(PathForm(imgid), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Log(err)
		w.Write([]byte("Error:Save Error."))
		return
	}
	defer f.Close()
	io.Copy(f, file)
	w.Write([]byte(imgid))
}

// Download image
func Download(req *restful.Request, resp *restful.Response) {

	imageid := req.PathParameter("imgid")
	if len([]rune(imageid)) != 16 {
		errors.Response(resp, errors.BadRequest("api.image.download", "Image ID incorrect."))
		return
	}
	imgpath := PathForm(imageid)
	if !FileExist(imgpath) {
		errors.Response(resp, errors.NotFound("api.image.download", "Image Not Found."))
		return
	}
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		imgpath)
}

// Home image
func Home(req *restful.Request, resp *restful.Response) {
	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatalf("Template gave: %s", err)
	}
	t.Execute(resp.ResponseWriter, nil)
}

// List images
func List(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		PathForm(""))
}