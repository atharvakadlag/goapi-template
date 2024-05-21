package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/atharvakadlag/splitfree/config"
	"github.com/atharvakadlag/splitfree/types"
	"github.com/atharvakadlag/splitfree/utils"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID":    strconv.Itoa(userID),
			"exipredAt": time.Now().Add(expiration).Unix(),
		})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)

		token, err := validateToken(tokenString)

		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		str := claims["userID"].(string)

		userID, _ := strconv.Atoi(str)

		u, err := store.GetUserByID(userID)

		if err != nil {
			log.Printf("user not found")
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserKey, u)

		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusNetworkAuthenticationRequired, fmt.Errorf("permission denied"))
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth == "" {
		return ""
	}

	return tokenAuth
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected siging method %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}
