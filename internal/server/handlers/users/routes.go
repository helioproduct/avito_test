package routes

import (
	"avito_api/internal/server"

	"avito_api/internal/usecase"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handlers struct {
	userUC usecase.UserUseCase
}

func NewHandlers(uc usecase.UserUseCase) *Handlers {
	return &Handlers{userUC: uc}
}

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "routes.users:"

	users, err := h.userUC.GetAllUsers()
	if err != nil {
		log.Println(op, err)
		render.Status(r, 500)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}
	render.JSON(w, r, users)
}

func RegisterRoutes(r chi.Router, h *Handlers) {
	r.Get("/", h.GetAllUsers)
}
