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

	publicStorage, pubErr := files.NewHostStorage(fmt.Sprintf("%s/public", dir), "http://fs.local/public")
	if pubErr != nil {
		log.WithError(pubErr).Fatal("Init public storage failed")
	}
	privateStorage, privErr := files.NewHostStorage(fmt.Sprintf("%s/private", dir), "http://fs.local/private")
	if privErr != nil {
		log.WithError(privErr).Fatal("Init public storage failed")
	}

	publicHandler := files.NewStorageHttpHandler(publicStorage, "/public", nil)
	privateHandler := files.NewStorageHttpHandler(privateStorage, "/private", func(s string, request *http.Request) bool {
		return request.Header.Get("h-auth") != ""
	})
	publicHandler.Logger = log.WithField("channel", "fs-server-public")
	privateHandler.Logger = log.WithField("channel", "fs-server-private")
	http.Handle("/public/", publicHandler)
	http.Handle("/private/", privateHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
