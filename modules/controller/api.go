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

// SearchForBook : this method searched for requested books using OpenLibrary API and
// storing the response data in a proper json format for use.
func (ct *Catalogue) SearchForBook(wr http.ResponseWriter, rq *http.Request) {
	var docs model.Docs

	// parse the form input
	if err := rq.ParseForm(); err != nil {
		ct.App.ErrorLogger.Fatalf("cannot parse post info :  %v \n", err)
	}

	// remove redundant space from the input values
	title := strings.ToLower(rq.Form.Get("title"))
	searchBook := strings.Replace(strings.TrimSpace(title), " ", "+", -1)

	// using the http client to get result/data from the API url
	resp, err := ct.Client.Get(fmt.Sprintf("https://openlibrary.org/search.json?q=%v", searchBook))
	if err != nil {
		ct.App.ErrorLogger.Fatalf("url error : %v \n", err)
	}

	// close/stop reading the response body from the request URL
	// after this function is done executing
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			ct.App.ErrorLogger.Fatalln(err)
		}
	}(resp.Body)

	// Read from the response body of the API request
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}

	apiData := string(data)

	// write the data read from the body to json file (api.json)
	err = os.WriteFile("./modules/json/api.json", []byte(apiData), 0666)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}

	// Open the JSON file and check for any error that might occur
	bookData, err := os.Open("./modules/json/api.json")
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}

	// close/stop reading data from the json file after this function is done executing
	defer func(bookData *os.File) {
		err := bookData.Close()
		if err != nil {
			ct.App.ErrorLogger.Println(err)
			return
		}
	}(bookData)

	// Read all the data stored in the JSON file as bytes of data
	byteData, err := io.ReadAll(bookData)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}

	// Encode the data into struct model format similar to the JSON format
	// Note : the data is the details of the searched book
	err = json.Unmarshal(byteData, &docs)
	if err != nil {
		ct.App.ErrorLogger.Println(err)
	}

	book := docs.Docs[0]

	//Check if Book is in the Library if not add the book to the library collections
	count, bookID, err := ct.CatDB.CheckLibrary(book.Title, book)
	if err != nil {
		ct.App.ErrorLogger.Fatalln("error while checking for book in the library")
	}

	ct.App.Session.Put(rq.Context(), "book_id", bookID)

	// conditions : check if the searched book is available in the Main Library/ Store
	// , not available or if an error pop up in the server
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
			Message:    "error while searching for book",
		}
		err := json.NewEncoder(wr).Encode(msg)
		if err != nil {
			return
		}
	}
}
