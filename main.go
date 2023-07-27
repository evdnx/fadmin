package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/evdnx/unixmint/constants"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/bbolt"
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

	// create data file if it doesn't exist yet
	if _, err := os.Stat(constants.DataFileName); err != nil && os.IsNotExist(err) {
		err = os.WriteFile(constants.DataFileName, []byte{}, 0600)
		if err != nil {
			glog.Fatal(err)
		}
	}

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

	// initialize rate limiter store
	storage := bbolt.New(bbolt.Config{
		Database: "unixmint.db",
		Bucket:   "ratelimit",
	})

	// rate limiter
	app.Use(limiter.New(limiter.Config{
		Max:        1,
		Expiration: 1 * time.Second,
		//LimiterMiddleware: limiter.SlidingWindow{},
		Storage: storage,
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
