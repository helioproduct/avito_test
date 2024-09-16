package routes

import (
	"avito_api/internal/server"
	"encoding/json"
	"errors"
	"strconv"

	"avito_api/internal/entities/bid"

	"avito_api/internal/usecase"

	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type BidHandlers struct {
	bidsUC usecase.BidUseCase
}

func NewHandlers(uc usecase.BidUseCase) *BidHandlers {
	return &BidHandlers{bidsUC: uc}
}

func (h *BidHandlers) CreateBidHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.bids.CreateBid:"

	var bid bid.Bid

	if err := json.NewDecoder(r.Body).Decode(&bid); err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("Неверный формат запроса"))
		return
	}

	createdBid, err := h.bidsUC.Create(&bid)
	if err != nil {
		log.Println(op, err)
	}

	if errors.Is(err, usecase.ErrUserNotFound) {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if errors.Is(err, usecase.ErrUserNotResponsible) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if errors.Is(err, usecase.ErrNoSuchTender) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	render.JSON(w, r, createdBid)
}

func (h *BidHandlers) MyBidsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.bids.MyBids:"
	// log.Println("hello")

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

	bids, err := h.bidsUC.GetBidsByUser(username, limit, offset)

	if err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	if len(bids) == 0 {
		render.JSON(w, r, struct{}{})
		return
	}

	render.JSON(w, r, bids)
}

func (h *BidHandlers) GetBidsForTenderHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.bid.GetBdisForTender:"

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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get offset"))
		return
	}

	if offset < 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("offset must be positive"))
		return
	}

	tenderID := r.PathValue("tenderId")
	if tenderID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get tenderId"))
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get username"))
		return
	}

	bids, err := h.bidsUC.GetBidsForTender(tenderID, username, limit, offset)
	if err != nil {
		log.Println(op, err)
	}

	if errors.Is(err, usecase.ErrNoSuchTender) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if errors.Is(err, usecase.ErrUserNotResponsible) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	if len(bids) == 0 {
		render.JSON(w, r, struct{}{})
		return
	}
	render.JSON(w, r, bids)
}

func (h *BidHandlers) GetBidStatusHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.bids.GetBidStatusHandler:"

	bidID := r.PathValue("bidId")
	if bidID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get bidID"))
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get username"))
		return
	}

	status, err := h.bidsUC.GetBidStatus(bidID, username)
	if err != nil {
		log.Println(op, err)
	}

	if errors.Is(err, usecase.ErrBidNotFound) || errors.Is(err, usecase.ErrUserNotFound) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if errors.Is(err, usecase.ErrUserNotResponsibleForBid) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, server.Error(err.Error()))
		return
	}

	render.JSON(w, r, status)
}

func (h *BidHandlers) ChangeBidStatusHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.bids.ChangeBidStatusHandler:"

	bidID := r.PathValue("bidId")
	if bidID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get bidID"))
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get username"))
		return
	}

	newStatus := r.URL.Query().Get("status")
	if newStatus == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get status"))
		return
	}
	updatedBid, err := h.bidsUC.ChangeStatus(bidID, username, bid.BidStatus(newStatus))
	if err != nil {
		log.Println(op, err)
	}

	if errors.Is(err, usecase.ErrInvalidStatus) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrUserNotResponsibleForBid) {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, err.Error())
		return
	} else if errors.Is(err, usecase.ErrBidNotFound) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, err.Error())
		return
	}
	render.JSON(w, r, updatedBid)
}

func (h *BidHandlers) EditBidHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.bids.ChangeBidStatusHandler:"

	bidID := r.PathValue("bidId")
	if bidID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get bidID"))
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get username"))
		return
	}

	var updatedInfo bid.Bid
	updatedInfo.ID = bidID

	if err := json.NewDecoder(r.Body).Decode(&updatedInfo); err != nil {
		log.Println(op, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("Неверный формат запроса"))
		return
	}

	bid, err := h.bidsUC.Edit(&updatedInfo, username)

	if errors.Is(err, usecase.ErrInvalidStatus) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrUserNotResponsibleForBid) {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, err.Error())
		return
	} else if errors.Is(err, usecase.ErrBidNotFound) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, bid)
}

func (h *BidHandlers) SubmitDesisionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "routes.bids.Submitdecision:"

	bidID := r.PathValue("bidId")
	if bidID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get bidID"))
		return
	}

	decision := r.URL.Query().Get("decision")
	if decision == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get decision"))
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error("couldn't get username"))
		return
	}

	bid, err := h.bidsUC.ChangeDecision(bidID, username, bid.Decision(decision))

	if err != nil {
		log.Println(op, err)
	}
	if errors.Is(err, usecase.ErrInvalidDecision) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, server.Error(err.Error()))
		return
	} else if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrUserNotResponsibleForBid) {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, err.Error())
		return
	} else if errors.Is(err, usecase.ErrBidNotFound) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, err.Error())
		return
	}

	render.JSON(w, r, bid)
}

func RegisterRoutes(r chi.Router, h *BidHandlers) {
	r.Post("/new", h.CreateBidHandler)
	r.Get("/my", h.MyBidsHandler)
	r.Get("/{tenderId}/list", h.GetBidsForTenderHandler)
	r.Get("/{bidId}/status", h.GetBidStatusHandler)
	r.Put("/{bidId}/status", h.ChangeBidStatusHandler)
	r.Patch("/{bidId}/edit", h.EditBidHandler)
	r.Put("/{bidId}/submit_decision", h.SubmitDesisionHandler)
}
