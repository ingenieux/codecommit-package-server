package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/go-zoo/bone"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"text/template"
	"time"
)

type RepoInfo struct {
	RepoId     string
	RepoUrl    string
	SourcePath string
	Protocol   string
}

var pageTemplate = template.Must(template.New("").Parse(`
<html>
  <head>
    <meta name="go-import" content="{{.RepoId}} git {{.RepoUrl}}">
    <title>Package Redirect for {{.RepoId}}</title>
  </head>
  <body>
    <p>go get {{.RepoId}}</p>
    <p>Source: <a href="{{.SourcePath}}">{{.SourcePath}}</a></p>
    <h2>Setup Instructions</h2>
    <ul>
      <li>See <a href="https://alestic.com/2015/11/aws-codecommit-iam-role/">this guide</a></li>
      <li>Or <a href="http://docs.aws.amazon.com/codecommit/latest/userguide/setting-up.html">read the AWS Docs</a></li>
      <li><b>Ubuntu 14.04 Users</b>: <a href="https://askubuntu.com/questions/186847/error-gnutls-handshake-failed-when-connecting-to-https-servers">Beware</a></li>
    </ul>
  </body>
</html>
`))

var hostname string = "codecommit.ingenieux.io"
var region string = "us-east-1"
var listenAddr string = "127.0.0.1:3001"
var defaultProto string = "https"

func main() {
	flag.StringVar(&hostname, "hostname", hostname, "hostname to use")
	flag.StringVar(&region, "region", region, "region to use")
	flag.StringVar(&listenAddr, "listenAddr", listenAddr, "Address to Listen")
	flag.StringVar(&defaultProto, "defaultProto", defaultProto, "Default Protocol (ssh / https)")

	flag.Parse()

	router := bone.New()

	router.Get("/", http.HandlerFunc(RootHandler))
	router.Get("/repo/:id", http.HandlerFunc(RepoHandler))
	router.Get("/repo/:id/*", http.HandlerFunc(RepoHandler))

	httpServer := &http.Server{
		Handler:      handlers.LoggingHandler(os.Stderr, router),
		Addr:         listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info("Starting")

	log.Fatal(httpServer.ListenAndServe())
}

func RootHandler(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "https://github.com/ingenieux/codecommit-package-server", 301)
}

func RepoHandler(response http.ResponseWriter, request *http.Request) {
	// https://git-codecommit.us-east-1.amazonaws.com/v1/repos/cowsurfing ||
	// ssh://git-codecommit.us-east-1.amazonaws.com/v1/repos/cowsurfing

	repoId := bone.GetValue(request, "id")

	sourcePath := fmt.Sprintf("https://console.aws.amazon.com/codecommit/home?region=%s#/repository/%s/browse/", region, repoId)

	if region != "us-east-1" {
		sourcePath = fmt.Sprintf("https://%s.console.aws.amazon.com/codecommit/home?region=%s#/repository/%s/browse/", region, region, repoId)
	}

	protocol := defaultProto

	if protocolParameter := request.URL.Query().Get("protocol"); "" != protocolParameter {
		protocol = protocolParameter
	}

	repoInfo := RepoInfo{
		RepoId:     fmt.Sprintf("%s/repo/%s", hostname, repoId),
		RepoUrl:    fmt.Sprintf("%s://git-codecommit.%s.amazonaws.com/v1/repos/%s", protocol, region, repoId),
		SourcePath: sourcePath,
		Protocol:   protocol,
	}

	response.WriteHeader(200)
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	pageTemplate.Execute(response, repoInfo)
}
