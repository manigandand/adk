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
	OrgID uint
}

func (c *createUserReq) Validate() *errors.AppError {
	if c.Email == "" {
		return errors.KeyRequired("email")
	}

	return nil
}

// CreateUserHandler creates a new users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) *errors.AppError {
	ctx := r.Context()
	orgID, ok := ctx.Value("orgID").(uint)
	if !ok {
		return errors.InternalServer("org id not set in context")
	}

	var createReq createUserReq
	if err := api.Decode(r, &createReq); err != nil {
		return err
	}

	createReq.OrgID = orgID

	if err := storeCreateUser(&createReq); err != nil {
		return err
	}

	// respond to the client
	return respond.OK(w, map[string]interface{}{
		"message": "user created successfully",
		"id":      123,
	})
}

func storeCreateUser(req *createUserReq) *errors.AppError {
	// save db
	// if err := db.Insert(req); err != nil {
	// 	return errors.InternalServer("couldn't able to create user").
	// 		Wrap(err, "create user")
	// }
	return nil
}
