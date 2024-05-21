package user

import (
	"fmt"
	"net/http"

	"github.com/atharvakadlag/splitfree/config"
	"github.com/atharvakadlag/splitfree/service/auth"
	"github.com/atharvakadlag/splitfree/types"
	"github.com/atharvakadlag/splitfree/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/", h.health).Methods("GET")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid data in payload"))
		return
	}

	u, err := h.store.GetUserByEmail(payload.Email)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid username %v", err))
		return
	}

	isOk := auth.ComparePasswords(u.Password, payload.Password)

	if !isOk {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid username or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error creating JWT"))
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)

	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", u.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("invalid password"))
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", err))
		return
	}

	utils.WriteJson(w, http.StatusCreated, nil)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusOK, "healthy")
}
