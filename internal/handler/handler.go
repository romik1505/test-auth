package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/romik1505/auth/docs"
	"github.com/romik1505/auth/internal/mapper"
	"github.com/romik1505/auth/internal/services/auth"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	AuthService *auth.AuthService
}

func NewHandler(as *auth.AuthService) *Handler {
	return &Handler{
		AuthService: as,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/login", ErrorWrapper(h.login)).Methods(http.MethodPost)
	router.HandleFunc("/register", ErrorWrapper(h.register)).Methods(http.MethodPost)

	return router
}

// @Summary login user account
// @Tags auth
// @ID login
// @Accept json
// @Produce json
// @Param input body mapper.LoginRequest true "account info"
// @Success 200 {object} mapper.TokenPair
// @Failure 400,500 {string} string
// @Router /login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	req := mapper.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return ErrorBadRequest
	}

	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidate, err)
	}

	tokens, err := h.AuthService.Login(ctx, req)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		return err
	}
	return nil
}

// @Summary register new user account
// @Tags auth
// @ID register
// @Accept json
// @Produce json
// @Param input body mapper.RegisterRequest true "register user data"
// @Success 200 {string} string
// @Failure 400,500 {string} string
// @Router /register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	req := mapper.RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return ErrorBadRequest
	}

	if err := req.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidate, err)
	}

	err := h.AuthService.Register(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
