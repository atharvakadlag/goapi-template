package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/atharvakadlag/splitfree/types"
	"github.com/go-playground/validator/v10"
)

type ContextKey string

const UserKey ContextKey = "user"

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("invalid payload, request body missing")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJson(w, status, map[string]string{"error": err.Error()})
}

func GetUserFromContext(ctx context.Context) types.User {
	user := ctx.Value(UserKey)

	if user == nil {
		return types.User{}
	}

	return user.(types.User)
}
