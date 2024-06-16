# ErrorEX

ErrorEX is a custom Go library for error handling that provides custom error types with additional details. It enables greater control over the creation, handling, and checking of errors, making it easier to identify and resolve issues in complex systems.

## Motivation

The motivation behind Errorex is to improve error management in Go by providing a means to associate unique codes and detailed messages with errors. This helps to standardize error responses in APIs and backend services, making the errors more expressive and easier to track.

## Usage

To use Errorex, you need to register error codes and their corresponding descriptions. After registration, you can create new custom errors and check if an error matches a particular error code.

## Advantages

- Error Standardization: Provides a uniform structure for errors in Go applications.
- Additional Details: Allows associating additional details with errors for better diagnostics.
- Prevention of Duplicate Codes: Ensures that each error has a unique code.

## Possible Pitfalls

- **Panic in Cases of Unregistered Error:** The librarywill panic if you try to create an error with an unregistered code.
- **Strict Detail Type:** When create an error, the detail type must match the registered detail type.

### Installation

```bash
go get github.com/fkmatsuda/errorex
```

### Example

```go
package main

import (
    "fmt"
    "github.com/fkmatsuda/errorex"
)

init() {
    errorex.Register("E001", "Invalid input", ErrInvalidInputDetail{})
}

type ErrInvalidInputDetail struct {
    Field string `json:"field"`
    Value string `json:"value"`
}

func main() {
    err := errorex.New("E001", ErrInvalidInputDetail{Field: "name", Value: "John"})

    if errorex.Is(err, "E001") {
        fmt.Println("Invalid input error")
        fmt.Println(err)
        return

        // Output:
        // Invalid input error
        // {
        //   "code": "E001",
        //   "message": "Invalid input",
        //   "detail": {
        //     "field": "name",
        //     "value": "John"
        //   }
        // }
    }

    fmt.Println("Other error")

}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
