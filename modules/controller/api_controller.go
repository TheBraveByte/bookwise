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

func (ct *Catalogue) SearchBookTitle(wr http.ResponseWriter, rq *http.Request) {
	var docs model.Docs
	if err := rq.ParseForm(); err != nil {
		ct.App.ErrorLogger.Fatalf("cannot parse post info :  %v \n", err)
	}
	title := strings.ToLower(rq.Form.Get("title"))
	searchBook := strings.Replace(strings.TrimSpace(title), " ", "+", -1)

	resp, err := ct.Client.Get(fmt.Sprintf("https://openlibrary.org/search.json?q=%v", searchBook))
	if err != nil {
		ct.App.ErrorLogger.Fatalf("url error : %v \n", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ct.App.ErrorLogger.Fatalln(err)
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}
	apiData := string(data)

	err = os.WriteFile("./modules/json/api.json", []byte(apiData), 0666)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}

	bookData, err := os.Open("./modules/json/api.json")
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}

	defer func(bookData *os.File) {
		err := bookData.Close()
		if err != nil {
			ct.App.ErrorLogger.Println(err)
			return
		}
	}(bookData)

	byteData, err := io.ReadAll(bookData)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}
	err = json.Unmarshal(byteData, &docs)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}
	book := docs.Docs[0]

	//Check if Book is in the Catalogue/ Library

	count, bookID, err := ct.CatDB.CheckLibrary(book.Title, book)
	if err != nil {
		ct.App.ErrorLogger.Fatalln("error while checking for book in the library")
	}

	ct.App.Session.Put(rq.Context(), "book_id", bookID)

	if count >= 1 {
		msg := model.ResponseMessage{
			StatusCode: http.StatusOK,
			Message:    fmt.Sprintf("%v :Book Found in Library ! \n %v", book.Title, &book),
		}
		err = json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	} else if count == 0 {
		msg := model.ResponseMessage{
			StatusCode: http.StatusOK,
			Message:    fmt.Sprintf(" %v : Book not found in Library \n Adding new Book to Library .... \n Search again", book.Title),
		}
		err = json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	} else {
		msg := model.ResponseMessage{
			StatusCode: http.StatusInternalServerError,
			Message:    "error while ",
		}
		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	}
}
