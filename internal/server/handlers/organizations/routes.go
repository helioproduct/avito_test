package routes

import (
	"avito_api/internal/server"

	uc "avito_api/internal/usecase"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handlers struct {
	orgUC uc.OrganizationUseCase
}

func NewHandlers(uc uc.OrganizationUseCase) *Handlers {
	return &Handlers{orgUC: uc}
}

func (h *Handlers) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	const op = "routes.organizations:"

	allOrganizations, err := h.orgUC.GetAllOrganizations()

	if err != nil {
		log.Println(op, err)
		render.Status(r, 500)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}
	render.JSON(w, r, allOrganizations)
}

func (h *Handlers) GetOrganizationsResponsible(w http.ResponseWriter, r *http.Request) {
	const op = "routes.OrganizationsResponsible:"

	orgResp, err := h.orgUC.GetResponsibleUsersForOrganization()

	if err != nil {
		log.Println(op, err)
		render.Status(r, 500)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}
	render.JSON(w, r, orgResp)
}

func RegisterRoutes(r chi.Router, h *Handlers) {
	r.Get("/all", h.GetOrganizations)
	r.Get("/responsible", h.GetOrganizationsResponsible)
}
