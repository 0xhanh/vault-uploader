package main

import (
	"log"
	"os"
	"vault-uploader/cloud"
	"vault-uploader/config"
)

func main() {
	conf, err := config.ReadConf("vault-config.yml")
	if err != nil {
		log.Println("vault-uploader: error read config", err.Error())
		return
	}
	args := os.Args[1:]
	if len(args) != 2 {
		log.Println("invalid args")
		return
	}
	// $MTX_PATH $MTX_SEGMENT_PATH
	path := args[0]
	filename := args[1]

	// path := "hobao"
	// filename := "/vms/data/hobao/2023-10-30_16-43-17.893403.mp4"

	cloud.HandleUpload(path, filename, conf)
}
