# Bookwise: Empowering Your Book Experience

Bookwise presents a dynamic RESTful bookstore API, seamlessly integrating with the Flutterwave payment gateway. Explore a wide array of books from Open Library and effortlessly expand your collection. Harnessing cutting-edge technologies, Bookwise ensures secure payments, is powered by Flutterwave, and employs robust authentication using JSON Web Tokens (JWT).

📖 **Features:**
- Swift book searches by title
- Precise search result filtering
- Seamlessly add books to your library
- Secure payments via Flutterwave
- Rock-solid authentication and authorization

🚀 **Technologies:**
- GoLang
- MongoDB
- Docker
- HTTP Session Management (kataras/go-sessions/v3)
- Flutterwave API (Rave)
- Open Library API
- JSON Web Tokens (JWT)

🔧 **Getting Started:**
1. Prerequisites: GoLang, MongoDB
2. Clone the repository: `git clone https://github.com/akinbyte/bookwise.git`
3. Set environment variables: `FLUTTERWAVE_PUBLIC_KEY` & `FLUTTERWAVE_SECRET_KEY`
4. Launch the server.

📚 **API Endpoints:**
- `/view/books` - Browse available books
- `/create/account` - Register a new user
- `/login/account` - User login
- `/api/user/search-book` - Search for book by title
- `/api/user/pay/details` - Submit payment for a book
- `/api/user/validate` - Validate payment details
- `/api/user/view/library` - View all available books
- `/api/user/view/books` - Access book info
- `/api/user/delete/book/{id}` - Delete the book from the catalog
- `/api/user/search/book/{id}` - Search catalog for a book
- `/api/add/new/book` - Add new book (authentication required)

🔒 **Security Measures:**
- JWT-based Authentication and Authorization
- All requests mandate a valid JWT token
- Passwords are bcrypt-hashed before storage

Explore our Postman documentation [here](https://documenter.getpostman.com/view/24714144/2s8Z6yXDMj) for a comprehensive overview.

Embark on an enriched book journey with Bookwise today! 📚✨
