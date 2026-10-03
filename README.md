# bookwise

A bookstore REST API in Go (2023). Users search books from the Open Library API,
add them to a personal library, and pay for them through Flutterwave.

- Go with the chi router, MongoDB, Docker
- JWT authentication; passwords hashed with bcrypt
- Flutterwave for payments, Open Library for the catalogue

## Running it

Requires Go and MongoDB.

```sh
git clone https://github.com/TheBraveByte/bookwise.git
cd bookwise
export FLUTTERWAVE_PUBLIC_KEY=... FLUTTERWAVE_SECRET_KEY=...
go run ./web/cmd
```

## Endpoints

| Route | Purpose |
|---|---|
| `POST /create/account`, `POST /login/account` | Register and log in |
| `GET /view/books` | Browse books |
| `POST /api/user/search-book` | Search Open Library by title |
| `POST /api/user/pay/details`, `GET /api/user/pay/validate` | Pay for a book and confirm the charge |
| `GET /api/user/view/books` | The user's library |
| `GET /api/user/delete/book/{id}` | Remove a book from the library |

Full request examples: [Postman documentation](https://documenter.getpostman.com/view/24714144/2s8Z6yXDMj).

Archived. An early project, kept for reference; not maintained.
