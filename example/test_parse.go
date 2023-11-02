package main

import (
	"fmt"
	"vault-uploader/cloud"
)

func main() {
	f, _ := cloud.ToKerberosFormat("hobao", "2006-01-02_15-04-05-893403.mp4")
	fmt.Println("Fu: " + f)
}
