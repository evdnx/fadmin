package main

import (
	"embed"
	"flag"
	"fmt"
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

	// create new fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// embed ui into program binary
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(embedFS),
		//PathPrefix: "app",
		//Browse: true,
	}))

	// use and config recovery middleware with custom stacktrace handler
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			glog.Errorf("\npanic: %v\n%s\n", e, debug.Stack())
		},
	}))

	// start app
	err := app.Listen(fmt.Sprint(":", *port))
	if err != nil {
		glog.Fatal(err)
	}
}
