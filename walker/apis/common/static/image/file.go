package image

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/seehuhn/mt19937"
)

var imagedir = "/tmp/image"

// NewImageID make image id
func NewImageID() string {
	var imgid string
	for {
		imgid = makeID()
		if !FileExist(PathForm(imgid)) {
			break
		}
	}
	return imgid
}

func makeID() string {
	mt := mt19937.New()
	mt.Seed(time.Now().UnixNano())
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, mt.Uint64())
	return strings.ToUpper(hex.EncodeToString(buf))
}

// FileExist file exist
func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	}
	return true
}

// BuildTree build image dir
func BuildDir(imageid string) error {
	return os.MkdirAll(fmt.Sprintf("%s", imagedir), 0777)
}

// PathForm path form image_id
func PathForm(imageid string) string {
	return fmt.Sprintf("%s/%s", imagedir, imageid)
}