package routes

import (
	"errors"
	"net/http"
	"strconv"

	"avito_api/internal/server"
	uc "avito_api/internal/usecase"

	"avito_api/internal/entities/tender"
	"encoding/json"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handlers struct {
	tenderUseCase uc.TenderUseCase
}

func NewHandlers(uc uc.TenderUseCase) *Handlers {
	return &Handlers{tenderUseCase: uc}
}

func RegisterRoutes(r chi.Router, h *Handlers) {
	r.Get("/", h.GetTendersHandler)
	r.Post("/new", h.CreateTenderHandler)
	r.Get("/my", h.GetMyTendersHandler)
	r.Get("/{tenderId}/status", h.GetTenderStatus)
	r.Put("/{tenderId}/status", h.ChangeTenderStatus)
	r.Patch("/{tenderId}/edit", h.EditTender)
}

func (h *Handlers) GetTendersHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.GetTenders"

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldnt get limit"))
		return
	}

	if limit < 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("limit must be postive"))
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		log.Println(op, err)
		render.JSON(w, r, server.Error("couldn't get offset"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	if offset < 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("offset must be positive"))
		return
	}

	serviceType := r.URL.Query().Get("service_type")
	tenders, err := h.tenderUseCase.GetTenders(limit, offset, serviceType)

	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusInternalServerError)
		// render.JSON(w, r, server.Error("Error fetching tenders"))
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	if len(tenders) == 0 {
		render.JSON(w, r, struct{}{})
		return
	}

	render.JSON(w, r, tenders)
}

type tenderRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationId string `json:"organizationId"`
	ServiceType    string `json:"serviceType"`
	Username       string `json:"creatorUsername"`
	CreatorID      string `json:"creatorID"`
}

func toTenderEntity(req tenderRequest) *tender.Tender {
	return &tender.Tender{
		Name:           req.Name,
		Description:    req.Description,
		ServiceType:    tender.ServiceType(req.ServiceType),
		OrganizationID: req.OrganizationId,
		CreatorID:      req.CreatorID,
	}
}

func (h *Handlers) CreateTenderHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.CreateTender:"

	var tr tenderRequest

	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("invalid request payload"))
		return
	}

	userID, err := h.tenderUseCase.GetUserID(tr.Username)
	if err != nil {
		log.Println(err)
	}

	if userID == "" || err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, server.Error("пользователь не существует или некорректен"))
		return
	}

	tr.CreatorID = userID
	createdTender := toTenderEntity(tr)

	tenderID, err := h.tenderUseCase.CreateTender(createdTender, tr.Username)
	if errors.Is(err, uc.ErrUserNotResponsible) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, server.Error("Пользователь не ответсвенен за организацию"))
		return
	}

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, server.Error("Error creating tender"))
		return
	}

	createdTender.ID = tenderID
	render.JSON(w, r, createdTender)
}

func (h *Handlers) GetMyTendersHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.GetMyTenders:"

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldnt get limit"))
		return
	}

	if limit < 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("limit must be postive"))
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		log.Println(op, err)
		render.JSON(w, r, server.Error("couldn't get offset"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	if offset < 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("offset must be positive"))
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.JSON(w, r, server.Error("couldn't get username"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	userID, err := h.tenderUseCase.GetUserID(username)
	if userID == "" || err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, server.Error("пользователь не существует или некорректен"))
		return
	}

	myTenders, err := h.tenderUseCase.GetMyTenders(username, limit, offset)
	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, server.Error("Error creating tender"))
		return
	}

	if len(myTenders) == 0 {
		render.JSON(w, r, struct{}{})
		return
	}
	render.JSON(w, r, myTenders)
}

func (h *Handlers) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	const op = "routes.GetTenderStatus:"

	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		render.JSON(w, r, server.Error("couldn't get tenderId"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.JSON(w, r, server.Error("couldn't get username"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	userId, err := h.tenderUseCase.GetUserID(username)
	if err != nil || userId == "" {
		log.Println(op, err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, server.Error("Пользователь не существует или некорректен"))
		return
	}

	status, err := h.tenderUseCase.GetStatus(tenderID, username)
	if errors.Is(err, uc.ErrNoSuchTender) {
		log.Println(op, err)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, server.Error("Тендер не найден"))
		return
	}

	if errors.Is(err, uc.ErrUserNotResponsible) {
		log.Println(op, err)
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, server.Error("Недостаточно прав"))
		return
	}
	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, struct{ Status tender.StatusType }{Status: status})
}

func (h *Handlers) ChangeTenderStatus(w http.ResponseWriter, r *http.Request) {
	const op = "routes.ChangeTenderStatus:"

	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		render.JSON(w, r, server.Error("couldn't get tenderId"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.JSON(w, r, server.Error("couldn't get username"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	newStatus := r.URL.Query().Get("status")
	if newStatus == "" {
		render.JSON(w, r, server.Error("couldn't get status"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	userId, err := h.tenderUseCase.GetUserID(username)
	if err != nil || userId == "" {
		log.Println(op, err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, server.Error("Пользователь не существует или некорректен"))
		return
	}

	err = h.tenderUseCase.ChangeStatus(tenderID, username, tender.StatusType(newStatus))
	if errors.Is(err, uc.ErrNoSuchTender) {
		log.Println(op, err)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, server.Error("Тендер не найден"))
		return
	}
	if errors.Is(err, uc.ErrUserNotResponsible) {
		log.Println(op, err)
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, server.Error("Недостаточно прав"))
		return
	}
	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	tender, err := h.tenderUseCase.GetByID(tenderID, username)
	if err != nil {
		log.Println("not hello")
		log.Println(op, err)
		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, tender)
}

func (h *Handlers) EditTender(w http.ResponseWriter, r *http.Request) {
	const op = "routes.EditTender:"

	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		render.JSON(w, r, server.Error("couldn't get tenderId"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.JSON(w, r, server.Error("couldn't get username"))
		render.Status(r, http.StatusBadRequest)
		return
	}

	userId, err := h.tenderUseCase.GetUserID(username)
	if err != nil || userId == "" {
		log.Println(op, err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, server.Error("Пользователь не существует или некорректен"))
		return
	}

	var updatedTender *tender.Tender
	err = json.NewDecoder(r.Body).Decode(&updatedTender)
	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("Данные неправильно сформированы или не соответствуют требованиям"))
		return
	}

	updatedTender.ID = tenderID
	err = h.tenderUseCase.EditTender(updatedTender, username)

	if errors.Is(err, uc.ErrNoSuchTender) {
		log.Println(op, err)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	if errors.Is(err, uc.ErrUserNotResponsible) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, server.Error("недостаточно прав для выполнения действия"))
		return
	}

	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	updatedTender, err = h.tenderUseCase.GetByID(tenderID, username)
	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, server.Error("Ошибка получения обновленного тендера"))
		return
	}

	render.JSON(w, r, updatedTender)
}
