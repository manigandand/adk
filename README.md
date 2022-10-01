## ADK - Golang API Development Kit

[![CircleCI](https://circleci.com/gh/manigandand/adk/tree/master.svg?style=shield)](https://circleci.com/gh/manigandand/adk/tree/master)
[![Go Report](https://goreportcard.com/badge/github.com/manigandand/adk)](https://goreportcard.com/report/github.com/manigandand/adk)
[![GolangCI](https://golangci.com/badges/github.com/manigandand/adk.svg)](https://golangci.com/r/github.com/manigandand/adk)
[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/manigandand/adk/blob/master/LICENSE)
[![](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/manigandand/adk)

**GoDoc Reference:**
- api: [![](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/manigandand/adk/api)
- errors: [![](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/manigandand/adk/errors)
- respond: [![](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/manigandand/adk/respond)

---

[![GOPHERCON 2022 | ADK - Golang API Development Kit](https://img.youtube.com/vi/opReKsCXsA0/0.jpg)](https://www.youtube.com/watch?v=opReKsCXsA0)

---

Common utilities to write simple apis in golang.

```shell
- Custom API Handlers
  - Custom API Request(JSON) Decoders
  - Custom API URL Query-params Decoders using gorrila schema.
- App Errors
- Response Writers
```

> Ex: **Conventional way of writing api**

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/user", CreateUserHandler)
	http.ListenAndServe(":3000", r)
}

// CreateUserHandler creates a new users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

	reqBts, err := ioutil.ReadAll(r.Body)
	if err != nil {
        w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "Couldn't read request body"}`))
		return
	}

	var createReq createUserReq
	jErr := json.Unmarshal(reqBts, &createReq)
	if jErr != nil {
        w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "invalid request body"}`))
		return
	}
    if err := validateCreateUserReq(&createReq); err != nil {
        w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "invalid request body"}`))
		return
    }
    if err := store.CreateUser(&createReq); err != nil {
        w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"error": "invalid request body"}`))
		return
    }

    // respond to the client
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "user created successfully", "id": 123}`))
}
```

### After using custom handlers and ADK

```go
package main

import (
	"net/http"

    "github.com/manigandand/adk/api"
    "github.com/manigandand/adk/errors"
    "github.com/manigandand/adk/middleware"
    "github.com/manigandand/adk/respond"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

	r.Method(http.MethodPost, "/user", api.Handler(CreateUserHandler))
	http.ListenAndServe(":3000", r)
}

type createUserReq struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

func (c *createUserReq) Validate() *errors.AppError {
    if c.Email == "" {
        return errors.KeyRequired("email")
    }
    return nil
}

// CreateUserHandler creates a new users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) *errors.AppError{
    ctx := r.Context()

	var createReq createUserReq

	if err := api.Decode(r, &createReq); err != nil {
		return err
	}
    if err := store.CreateUser(&createReq); err != nil {
		return err
    }

    // respond to the client
    return respond.OK(w, map[string]interface{}{
        "message": "user created successfully",
        "id": 123,
    })
}
```

> NOTE:
> Decoder currently will work on only on the `"application/json"` body.
