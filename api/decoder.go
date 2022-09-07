package api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/manigandand/adk/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// custom validator interface
type ok interface {
	Validate() *errors.AppError
}

// Decode - decodes the request body and extends the validator interface with the Validate() method
//
// EX:
// type User struct {
// 	Email        string       `json:"email"`
// 	Name         string       `json:"name"`
// }
//
// func (c *Component) Validate() *errors.AppError {
// 	if c.Email == "" {
// 		return errors.IsRequiredErr("email")
// 	}
// 	return nil
// }
func Decode(r *http.Request, v interface{}) *errors.AppError {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.UnprocessableEntity("unmarshal request payload").
			AddDebug(err)
	}

	// custom validator interface
	if payload, ok := v.(ok); ok {
		return payload.Validate()
	}
	return nil
}

// JustDecode just decodes the request body
func JustDecode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

// FormDecoder gorrila schema decoder is used to decode form data & query params into structs.
// And it validates the data against the struct types.
// supported types:
// - bool,
// - float32, float64,
// - int, int8, int16, int32, int64,
// - string, []string,
// - uint, uint8, uint16, uint32, uint64,
// - struct
// and we can write custom decoders for custom types. example time.Time, uuid.UUID, etc..
var FormDecoder *schema.Decoder

// init the decoder
func init() {
	FormDecoder = schema.NewDecoder()
	FormDecoder.ZeroEmpty(true)
	FormDecoder.IgnoreUnknownKeys(true)
	FormDecoder.RegisterConverter(time.Time{}, parseFilterTime)
	FormDecoder.RegisterConverter(uuid.NullUUID{}, parseFilterUUID)
	FormDecoder.RegisterConverter(primitive.NilObjectID, parseFilterObjectID)
}

// register custom decoder for time.Time
func parseFilterTime(date string) reflect.Value {
	if s, err := time.Parse(time.RFC3339, date); err == nil {
		return reflect.ValueOf(s)
	}

	return reflect.Value{}
}

// register custom decoder for uuid.UUID
func parseFilterUUID(id string) reflect.Value {
	if s, err := uuid.Parse(id); err == nil {
		return reflect.ValueOf(s)
	}

	return reflect.Value{}
}

func parseFilterObjectID(id string) reflect.Value {
	if s, err := primitive.ObjectIDFromHex(id); err == nil {
		return reflect.ValueOf(s)
	}

	return reflect.Value{}
}
