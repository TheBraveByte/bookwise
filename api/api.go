package api

import (
	"encoding/json"
	"fmt"
	"github.com/yusuf/p-catalogue/dependencies/config"
	"github.com/yusuf/p-catalogue/query"
	"github.com/yusuf/p-catalogue/query/api"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"os"
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
	ApiDB  query.APIRepo
}

//var OpenAPI *OpenLibraryAPI

func NewOpenLibraryAPI(app *config.CatalogueConfig, db *mongo.Client) *OpenLibraryAPI {
	return &OpenLibraryAPI{
		Client,
		app,
		api.NewDBRepoAPI(app, db),
	}
}

func (ap *OpenLibraryAPI) SearchBookTitle(wr http.ResponseWriter, rq *http.Request) {
	var docs Docs
	if err := rq.ParseForm(); err != nil {
		ap.App.ErrorLogger.Fatalln("cannot parse post info :  %v", err)
	}
	title := strings.ToLower(rq.Form.Get("title"))
	searchBook := strings.Replace(strings.TrimSpace(title), " ", "+", -1)

	resp, err := ap.Client.Get(fmt.Sprintf("https://openlibrary.org/search.json?q=%v", searchBook))
	if err != nil {
		ap.App.ErrorLogger.Fatalln("url error : %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ap.App.ErrorLogger.Fatalln(err)
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ap.App.ErrorLogger.Println(err)
	}
	apiData := string(data)

	err = os.WriteFile("./api/api.json", []byte(apiData), 0666)
	if err != nil {
		ap.App.ErrorLogger.Println(err)
	}

	bookData, err := os.Open("./api/api.json")
	if err != nil {
		ap.App.ErrorLogger.Println(err)
	}

	defer func(bookData *os.File) {
		err := bookData.Close()
		if err != nil {
			ap.App.ErrorLogger.Println(err)
			return
		}
	}(bookData)

	byteData, err := io.ReadAll(bookData)
	if err != nil {
		ap.App.ErrorLogger.Println(err)
	}
	err = json.Unmarshal(byteData, &docs)
	if err != nil {
		ap.App.ErrorLogger.Println(err)
	}
	book := docs.Docs[0]

	count, err := ap.ApiDB.AddBook(book.Title, book)

	if count >= 1 && err != nil {
		ap.App.ErrorLogger.Fatalln("cannot add book to catalogue, it exist already")
	} else if count == 0 && err != nil {
		ap.App.ErrorLogger.Fatalln("error while counting database documents")
	} else {
		err := json.NewEncoder(wr).Encode(&book)
		if err != nil {
			return
		}
	}

}

func (ap *OpenLibraryAPI) ServeApiToController(wr http.ResponseWriter, rq *http.Request) {
	//	this api handlers will find/fetch the searched book from the
	//	database and make it available to the user userHandler
}
