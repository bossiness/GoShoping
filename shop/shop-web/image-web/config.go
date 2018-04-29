package imageweb

import (
    "encoding/json"
    "log"
    "os"
)

type Config struct {
    ListenAddr string
    Storage string
}

var conf Config
var ConfPath = "image-web/config.json"

func LoadConf() {
	r, err := os.Open(ConfPath)
    if err != nil {
        log.Fatalln(err)
    }
    decoder := json.NewDecoder(r)
    err = decoder.Decode(&conf)
    if err != nil {
        log.Fatalln(err)
    }
    log.Printf("Load Config : %s", ConfPath)
}