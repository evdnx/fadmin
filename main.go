package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"runtime/debug"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang/glog"
)

//go:embed ui/dist/pwa/*
var embedFS embed.FS

func main() {
	// define and set program flags
	// logging
	flag.Set("logtostderr", "false")
	flag.Set("alsologtostderr", "false")
	flag.Set("log_dir", "/var/log/unixmint")
	// port
	port := flag.Int("port", 3000, "")

	// read program flags
	flag.Parse()

	// setup logging
	glog.MaxSize = 16777216 // 16 MB

	// create new fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// use and config recovery middleware with custom stacktrace handler
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			glog.Errorf("\npanic: %v\n%s\n", e, debug.Stack())
		},
	}))

	// embed ui into program binary
	f, err := fs.Sub(embedFS, "ui/dist/pwa")
	if err != nil {
		glog.Fatal(err)
	}

	app.Use(filesystem.New(filesystem.Config{
		Root:         http.FS(f),
		NotFoundFile: "404.html",
	}))

	// start app
	err = app.Listen(fmt.Sprint(":", *port))
	if err != nil {
		glog.Fatal(err)
	}
}
