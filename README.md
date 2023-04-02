# Bookwise

Bookwise is a RESTful bookstore API that allows users to search for books and securely add them to their library or book catalogue through the integration of the Flutterwave payment system gateway. It interacts with an external search endpoint provided by Open Library to enable users to filter their search results based on the book title.

## Features
- Search for books by title
- Filter search results by book title
- Add books to your library or book catalogue
- Secure payment processing through the Flutterwave payment gateway
- Authentication and Authorization using JSON Web Tokens (JWT)

## Technologies Used
- GoLang
- MongoDB
- Docker
- HTTP Session Management for Go (kataras/go-sessions/v3)
- Flutterwave API (Rave)
- Open Library API
- JSON Web Tokens (JWT)

## Getting Started
### Prerequisites
- GoLang
- MongoDB

## Installation

clone the repository 
```go
git clone https://github.com/dev-ayaa/bookwise.git
```

## Install dependencies:

- Set up environment variables:

- Replace the values for **FLUTTERWAVE_PUBLIC_KEY** and **FLUTTERWAVE_SECRET_KEY** with your Flutterwave API keys
- Start the server:

## API Endpoints
### Books

`/view/books` - view available books

### User Access

`/create/account` - Register a new user

`/login/account` - Login a user

### Authentication

`/api/user/search-book` - search for book title

`/api/user/pay/details` - submit user details to pay for a specific book

`/api/user/validate` - validate user payment details

`/api/user/view/library` - Check for all available book in the library

`/api/user/view/books` - Check for a book info from the library

`/api/user/delete/book/{id}` - user delete book from their catalogue

`/api/user/search/book/{id}` - user search for book in the catalogue

### Authentication to add new book

`/api/add/new/book` - add new book for authorize user

### Error Handling

- All errors are returned in a standardized JSON format with a status and message property.
Security
- Authentication and Authorization are handled using JSON Web Tokens (JWT).
- All requests to the API require a valid JWT token in the Authorization header.
- User passwords are hashed using the bcrypt library before being stored in the database.
Deployment
- This API can be deployed to any cloud platform that supports Node.js and MongoDB, such as AWS or Heroku.

### Conclusion
BookwiseAPI is a powerful and secure RESTful API that allows users to search for and add books to their library using the Flutterwave payment gateway. With its integration with the Open Library API, users can easily filter search results to find the books they are looking for.




