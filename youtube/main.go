package main

import (
	"net/http"
	"net/url"
	"path"
	"youtube/config"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"time":"${time_rfc3339}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}"` + "\n",
	}))
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {

		baseURL, _ := url.Parse(config.Config.BaseURL)
		// https://www.googleapis.com/youtube/v3/search になるように path を設定
		baseURL.Path = path.Join(baseURL.Path, "search")

		// TODO: データ取得実装
		//https://developers.google.com/youtube/v3/sample_requests
		//https://developers.google.com/youtube/v3/docs/search/list?apix=true&apix_params=%7B%22part%22%3A%22snippet%22%2C%22q%22%3A%22YouTube%20Data%20API%22%2C%22type%22%3A%22video%22%2C%22videoCaption%22%3A%22closedCaption%22%7D

		return c.String(http.StatusOK, baseURL.Path)
	})

	e.Logger.Fatal(e.Start(config.Config.Port))
}
