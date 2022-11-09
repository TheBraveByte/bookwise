package api

import (
	"encoding/json"
	"fmt"
	"github.com/yusuf/p-catalogue/pkg/config"
	"io"
	"net/http"
	"strings"
	"time"
)

var Client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:       3,
		DisableCompression: true,
		MaxConnsPerHost:    3,
		IdleConnTimeout:    100 * time.Second,
	},
}

type OpenLibraryAPI struct {
	Client *http.Client
	App    *config.CatalogueConfig
}

var OpenAPI *OpenLibraryAPI

func NewOpenLibraryAPI(app *config.CatalogueConfig) *OpenLibraryAPI {
	return &OpenLibraryAPI{
		Client,
		app,
	}
}

func (ap *OpenLibraryAPI) SearchBookTitle(wr http.ResponseWriter, rq *http.Request) {
	var book Book
	if err := rq.ParseForm(); err != nil {
		ap.App.ErrorLogger.Fatalln("cannot parse post info :  %v", err)
	}
	title := strings.ToTitle(rq.Form.Get("title"))
	searchBook := strings.Replace(strings.TrimSpace(title), " ", "+", -1)

	resp, err := ap.Client.Get(fmt.Sprint("http://openlibrary.org/search.json?q=%v", searchBook))
	if err != nil {
		ap.App.ErrorLogger.Fatalln("url error : %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ap.App.ErrorLogger.Fatalln(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(&book)
		if err != nil {
			return
		}
		return
	}

}

func (ap *OpenLibraryAPI) ConvertXMLToJSON() {
	return
}

func (ap *OpenLibraryAPI) ServerToController() {

}
