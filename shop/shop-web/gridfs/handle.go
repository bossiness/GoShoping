package gridfsweb

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/seehuhn/mt19937"
)

func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {

	var imgid string
	for {
		imgid = MakeImageID()
		if !FileExist(ImageID2Path(imgid)) {
			break
		}
	}
	//上传参数为uploadfile
	r.ParseMultipartForm(32 << 20)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error."))
		return
	}
	defer file.Close()
	//检测文件类型
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Upload Error."))
		return
	}
	filetype := http.DetectContentType(buff)
	if (filetype != "image/jpeg") && (filetype != "image/png") {
		w.Write([]byte("Error:Not JPEG."))
		return
	}
	//回绕文件指针
	log.Println(filetype)
	if _, err = file.Seek(0, 0); err != nil {
		log.Println(err)
	}

	f, err := h.db.GridFS("fs").Create(imgid)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error:Save Error."))
		return
	}
	defer f.Close()
	io.Copy(f, file)
	w.Write([]byte(imgid))
}

func (h *Handler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageid := vars["imgid"]
	if len([]rune(imageid)) != 16 {
		w.Write([]byte("Error:ImageID incorrect."))
		return
	}
	imgpath := ImageID2Path(imageid)
	if !FileExist(imgpath) {
		w.Write([]byte("Error:Image Not Found."))
		return
	}
	http.ServeFile(w, r, imgpath)
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(
		`<html><body>
			<form id="upload-form" action="./image" method="post" enctype="multipart/form-data" >
		　　　<input type="file" id="uploadfile" name="uploadfile" /> <br />
		　　　<input type="submit" value="Upload" />
			</form></body>
		</html>`,
	))
}

func makeImageID() string {
	mt := mt19937.New()
	mt.Seed(time.Now().UnixNano())
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, mt.Uint64())
	return strings.ToUpper(hex.EncodeToString(buf))
}
