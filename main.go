package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/upload", handleUpload,
			apis.RequireAdminAuth(),
			apis.ActivityLogger(app))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func handleUpload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	backupPath := "pb_data/backups"
	if err := os.MkdirAll(backupPath, 0777); err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join(backupPath, file.Filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields.</p>", file.Filename))
}
