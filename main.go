package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func healthCheck(appUrl string) {
	resp, err := http.Get(appUrl + "/api/health")
	if resp.StatusCode == http.StatusOK && err == nil {
		return
	}
	fmt.Println(resp, err, resp.StatusCode)
}

func main() {
	backendURL := os.Getenv("bakckend_url")

	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/upload", handleUpload,
			apis.RequireAdminAuth(),
			apis.ActivityLogger(app))
		return nil
	})

	// scheduler start
	if backendURL == "" {
		backendURL = app.Settings().Meta.AppUrl
	}
	ticker := time.NewTicker(14 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				healthCheck(backendURL)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	// scheduler end

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
