package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

//go:embed ui/dist/pwa/*
var embedFS embed.FS

func main() {
	// define program flags
	// logging
	flag.Bool("logtostderr", false, "")
	flag.Bool("alsologtostderr", false, "")
	flag.String("log_dir", "/var/log/linux-control", "")
	// port
	port := flag.Int("port", 3000, "")

	// read program flags
	flag.Parse()

	// setup logging
	glog.MaxSize = 16777216 // 16 MB

	// create new gin app
	r := gin.New()

	// embed ui into program binary
	r.StaticFS("/app", http.FS(embedFS))

	// use and config recovery middleware with custom stacktrace handler
	r.Use(gin.CustomRecovery(func(c *gin.Context, e any) {
		glog.Errorf("\npanic: %v\n%s\n", e, debug.Stack())
	}))

	// start app
	err := r.Run(fmt.Sprint(":", *port))
	if err != nil {
		glog.Fatal(err)
	}
}
