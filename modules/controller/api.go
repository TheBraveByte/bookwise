package controller

import (
	"encoding/json"
	"fmt"
	"github.com/yusuf/p-catalogue/model"
	"io"
	"net/http"
	"os"
	"strings"
)

func (cg *Catalogue) SearchBookTitle(wr http.ResponseWriter, rq *http.Request) {
	var docs model.Docs
	if err := rq.ParseForm(); err != nil {
		cg.App.ErrorLogger.Fatalf("cannot parse post info :  %v \n", err)
	}
	title := strings.ToLower(rq.Form.Get("title"))
	searchBook := strings.Replace(strings.TrimSpace(title), " ", "+", -1)

	resp, err := cg.Client.Get(fmt.Sprintf("https://openlibrary.org/search.json?q=%v", searchBook))
	if err != nil {
		cg.App.ErrorLogger.Fatalf("url error : %v \n", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			cg.App.ErrorLogger.Fatalln(err)
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		cg.App.ErrorLogger.Println(err)
	}
	apiData := string(data)

	err = os.WriteFile("./modules/json/api.json", []byte(apiData), 0666)
	if err != nil {
		cg.App.ErrorLogger.Println(err)
	}

	bookData, err := os.Open("./modules/json/api.json")
	if err != nil {
		cg.App.ErrorLogger.Println(err)
	}

	defer func(bookData *os.File) {
		err := bookData.Close()
		if err != nil {
			cg.App.ErrorLogger.Println(err)
			return
		}
	}(bookData)

	byteData, err := io.ReadAll(bookData)
	if err != nil {
		cg.App.ErrorLogger.Println(err)
	}
	err = json.Unmarshal(byteData, &docs)
	if err != nil {
		cg.App.ErrorLogger.Println(err)
	}
	book := docs.Docs[0]

	count, bookID, err := cg.CatDB.AddBook(book.Title, book)

	cg.App.Session.Put(rq.Context(), "book_id", bookID)
	if count >= 1 && err != nil {
		cg.App.ErrorLogger.Fatalln("cannot add book to catalogue, it exist already")
	} else if count == 0 && err != nil {
		cg.App.ErrorLogger.Fatalln("error while counting database documents")
	} else {
		err := json.NewEncoder(wr).Encode(&book)
		if err != nil {
			return
		}
	}
}
