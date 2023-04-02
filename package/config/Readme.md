# Package Config

This is a package used in the Bookwise API to hold reusable configurations for the application. This package exports a struct type `CatalogueConfig` that holds various field types.

### Installation

You can install this package using go get command:

```go
go get github.com/yusuf/bookwise/pkg/config
```

### Usage

```go

import (
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/kataras/go-sessions/v3"
	"github.com/bookwise/bookwise-api/pkg/config"
)

func main() {
    // Initialize CatalogueConfig
    catConfig := config.CatalogueConfig{
        InfoLogger:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
        ErrorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
        Client:      &http.Client{},
        Session:     sessions.New(sessions.Config{}),
        Validate:    validator.New(),
    }

    // Use catConfig in your application
}
```

### Fields

#### InfoLogger

The `InfoLogger` field is a `*log.Logger` instance used for logging informational messages in the application.

#### ErrorLogger

The `ErrorLogger` field is a `*log.Logger` instance used for logging error messages in the application.

#### Client

The `Client` field is an `*http.Client` instance used for making HTTP requests in the application.

#### Session

The `Session` field is a `*sessions.Sessions` instance used for managing sessions in the application.

#### Validate

The `Validate` field is a `*validator.Validate` instance used for validating structs and fields in the application.
