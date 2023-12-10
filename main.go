package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/evdnx/unixmint/auth"
	"github.com/evdnx/unixmint/db"
	mw "github.com/evdnx/unixmint/middleware"
	"github.com/evdnx/unixmint/store"
	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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

	// init db
	if err := db.Init(); err != nil {
		glog.Fatalln(err)
	}

	// get last login
	last_login, err := db.Read(db.AuthBucket, "last_login")
	if err == nil {
		lastLogin, err := time.Parse(time.RFC3339, last_login)
		if err == nil {
			now := time.Now().UTC()
			allowedTime := lastLogin.Add(24 * time.Hour)
			if allowedTime.Before(now) {
				// difference between now and allowedTime
				auth.Timer = time.AfterFunc(now.Sub(allowedTime), func() { auth.Logout() })
			}
		}
	}

	// init services
	err = auth.Init()
	if err != nil {
		glog.Fatal(err)
	}

	// create a new echo app
	e := echo.New()

	// use and config recovery middleware with custom stacktrace handler
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			glog.Errorf("\npanic: %v\n%s\n", err, stack)
			return err
		},
	}))

	// initialize rate limiter store
	//Database: "ratelimit.db",
	//Bucket:   "ratelimit",

	// rate limiter
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: store.NewRateLimiterPersistentStoreWithConfig(
			store.RateLimiterPersistentStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
	}))

	// embed ui into program binary
	f, err := fs.Sub(embedFS, "ui/dist/pwa")
	if err != nil {
		glog.Fatal(err)
	}

	assetHandler := http.FileServer(http.FS(f))
	e.GET("/", echo.WrapHandler(assetHandler))

	// auth not required for login
	login := e.Group("/")
	login.POST("/login", func(c echo.Context) error { return nil })

	// auth required for everything else
	api := e.Group("/")
	api.Use(mw.AuthMiddleware())

	// start app
	err = e.Start(fmt.Sprint(":", *port))
	if err != nil {
		glog.Fatal(err)
	}
}
