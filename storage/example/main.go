package main

import (
	"fmt"
	"github.com/chipx/go-ext/files"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	//http.Handle("/foo", fooHandler)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	publicStorage := storage.NewHostStorage(fmt.Sprintf("%s/public", dir))
	privateStorage := storage.NewHostStorage(fmt.Sprintf("%s/private", dir))
	publicHandler := storage.NewHttpHandler(publicStorage, "/public", nil)
	privateHandler := storage.NewHttpHandler(privateStorage, "/private", func(s string, request *http.Request) bool {
		return request.Header.Get("h-auth") != ""
	})
	publicHandler.Logger = log.WithField("channel", "fs-server-public")
	privateHandler.Logger = log.WithField("channel", "fs-server-private")
	http.Handle("/public/", publicHandler)
	http.Handle("/private/", privateHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}