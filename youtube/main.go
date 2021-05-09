package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
	"youtube/config"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const baseURL = "https://www.googleapis.com/youtube/v3/"

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

type searchResult struct {
	Items []struct {
		ID struct {
			Videoid string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Publishedat time.Time `json:"publishedAt"`
			Channelid   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

func main() {
	e := echo.New()

	renderer := &TemplateRenderer{
		templates: template.Must(template.New("").Delims("[[", "]]").ParseGlob("views/*.html")), // vue.jsとdelimsがかぶるので変更
	}
	e.Renderer = renderer

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"time":"${time_rfc3339}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}"` + "\n",
	}))
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
	})
	e.GET("/api/search", func(c echo.Context) error {
		q := c.QueryParam("q")
		if q == "" {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		url, _ := url.Parse(baseURL)
		url.Path = path.Join(url.Path, "search") // https://www.googleapis.com/youtube/v3/search になるように path を設定

		// パラメータをセット
		//https://developers.google.com/youtube/v3/docs/search/list?apix=true&apix_params=%7B%22part%22%3A%22snippet%22%2C%22q%22%3A%22YouTube%20Data%20API%22%2C%22type%22%3A%22video%22%2C%22videoCaption%22%3A%22closedCaption%22%7D
		queryParams := url.Query()
		queryParams.Set("key", config.Config.APIKey)
		queryParams.Set("part", "snippet")
		queryParams.Set("order", "relevance") // rating, viewCount
		queryParams.Set("regionCode", "jp")
		queryParams.Set("type", "video")
		queryParams.Set("q", q) //ここだけはUI状で入力した値をセット
		url.RawQuery = queryParams.Encode()

		resp, err := http.Get(url.String())
		if err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		jsonBytes := ([]byte)(buf.String())
		searchResult := new(searchResult)
		if err := json.Unmarshal(jsonBytes, searchResult); err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, searchResult)
	})

	e.Logger.Fatal(e.Start(config.Config.Port))
}
