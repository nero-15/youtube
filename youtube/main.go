package main

import (
	"bytes"
	"log"
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

		// パラメータをセット
		//https://developers.google.com/youtube/v3/docs/search/list?apix=true&apix_params=%7B%22part%22%3A%22snippet%22%2C%22q%22%3A%22YouTube%20Data%20API%22%2C%22type%22%3A%22video%22%2C%22videoCaption%22%3A%22closedCaption%22%7D
		queryParams := baseURL.Query()
		queryParams.Set("key", config.Config.APIKey)
		queryParams.Set("part", "snippet")
		queryParams.Set("order", "relevance") // rating, viewCount
		queryParams.Set("q", "inter")         // TODO: 入力値で検索できるようにする
		queryParams.Set("regionCode", "jp")
		queryParams.Set("type", "video")
		baseURL.RawQuery = queryParams.Encode()
		resp, err := http.Get(baseURL.String())
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)

		return c.String(http.StatusOK, buf.String())
	})

	e.Logger.Fatal(e.Start(config.Config.Port))
}
