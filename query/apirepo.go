package query

import "github.com/yusuf/p-catalogue/api"

type APIRepo interface {
	AddBook(title string, bookData api.Book) (int64, error)
}
