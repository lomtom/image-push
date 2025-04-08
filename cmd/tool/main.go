package main

import (
	"bytes"
	"flag"
	"github.com/lomtom/image_push/pkg/harbor"
	"log"
	"os"
	"strings"
)

var (
	address   = flag.String("address", "", "address, example: http://localhost:5000")
	username  = flag.String("username", "", "username, example: admin")
	password  = flag.String("password", "", "password")
	project   = flag.String("project", "", "project, example: library")
	archive   = flag.String("file", "", "archive path, example: /root/archive.tar.gz")
	skipTls   = flag.Bool("skipTls", true, "skip tls")
	chunkSize = flag.Int("chunkSize", 1024*1024*10, "chunk size")
)

func main() {
	flag.Parse()
	if *address == "" {
		log.Fatalln("address is required")
		return
	}
	file, err := os.ReadFile(*archive)
	if err != nil {
		log.Fatalln(err)
		return
	}
	reader := bytes.NewReader(file)
	// 获取文件名
	splits := strings.Split(*archive, "/")
	archiveName := splits[len(splits)-1]
	c, err := harbor.NewConfig(*address, *username, *password, *project, archiveName, *skipTls, *chunkSize, reader)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if err = c.Push(); err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("push success")
}
