package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed ui/dist/pwa/*
var embedFS embed.FS

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// embed ui into program binary
	app.Use("/app", filesystem.New(filesystem.Config{
		Root:       http.FS(embedFS),
		PathPrefix: "app",
		Browse:     true,
	}))

	port := flag.Int("port", 3000, "")
	log.Fatal(app.Listen(":" + fmt.Sprint(port)))
}
