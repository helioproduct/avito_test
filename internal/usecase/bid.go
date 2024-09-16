package usecase

import (
	"avito_api/internal/entities/bid"
	"avito_api/internal/entities/tender"
	"errors"
	"log"
	"reflect"
)

var (
	ErrInvalidAuthorType        = errors.New("invalid author type")
	ErrUserNotResponsibleForBid = errors.New("user is not responsible for this bid")
	ErrBidNotFound              = errors.New("bid not found")
	ErrInvalidStatus            = errors.New("invalid status")
	ErrInvalidDecision          = errors.New("invalid status")
)

type bidUseCase struct {
	bidRepo    BidRepo
	orgRepo    OrganizationRepo
	userRepo   UserRepo
	tenderRepo TenderRepo
}

func NewBidUseCase(bidRepo BidRepo,
	orgRepo OrganizationRepo, userRepo UserRepo,
	tenderRepo TenderRepo) BidUseCase {

	return &bidUseCase{
		bidRepo:    bidRepo,
		orgRepo:    orgRepo,
		userRepo:   userRepo,
		tenderRepo: tenderRepo,
	}
}

func (uc *bidUseCase) Create(b *bid.Bid) (*bid.Bid, error) {
	const op = "usecase.Bid.Create:"
	var err error

	user, err := uc.userRepo.GetUserByID(b.AuthorID)
	if err != nil {
		err = ErrUserNotFound
		log.Println(op, err)
		return nil, err
	}

	var responsible bool
	// создает от имени организации
	if b.AuthorType == "Organization" {
		orgID, err := uc.orgRepo.GetUserOrganizationID(user.ID)
		b.AuthorID = orgID
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		responsible, err = uc.orgRepo.IsUserResponsibleForOrganization(user.Username, orgID)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
	} else if b.AuthorType == "User" { // создает от своего имени
		responsible = true
	} else {
		return nil, ErrInvalidAuthorType
	}

	if !responsible {
		return nil, ErrUserNotResponsible
	}

	_, err = uc.tenderRepo.GetByID(b.TenderID)
	if err != nil {
		log.Println(op, ErrNoSuchTender)
		return nil, ErrNoSuchTender
	}

	b.Status = bid.Created
	b.Version = 1

	id, err := uc.bidRepo.CreateBid(b)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	createdBid, err := uc.GetBidByID(id)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	return createdBid, nil
}

// Create, Publish, Cancel (for owner)
func (uc *bidUseCase) ChangeStatus(bidID string, username string, newStatus bid.BidStatus) (*bid.Bid, error) {
	if newStatus != bid.Created && newStatus != bid.Canceled && newStatus != bid.Published {
		return nil, ErrInvalidStatus
	}

	hasPermission := uc.hasPermission(bidID, username)
	if !hasPermission {
		return nil, ErrUserNotResponsibleForBid
	}

	err := uc.bidRepo.ChangeStatus(bidID, username, newStatus)
	if err != nil {
		return nil, err
	}

	updatedBid, err := uc.GetBidByID(bidID)
	if err != nil {
		return nil, err
	}
	return updatedBid, nil
}

func (uc *bidUseCase) Edit(updatedBid *bid.Bid, username string) (*bid.Bid, error) {
	const op = "usecase.bid.Edit"
	hasPermission := uc.hasPermission(updatedBid.ID, username)

	if !hasPermission {
		return nil, ErrUserNotResponsibleForBid
	}

	currentBid, err := uc.GetBidByID(updatedBid.ID)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	updateBidStruct(currentBid, *updatedBid)
	err = uc.bidRepo.EditBid(currentBid)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return uc.GetBidByID(updatedBid.ID)
}

func (uc *bidUseCase) GetBidByID(bidID string) (*bid.Bid, error) {
	bid, err := uc.bidRepo.GetBidByID(bidID)
	if err != nil {
		return nil, ErrBidNotFound
	}
	return bid, nil
}

func (uc *bidUseCase) MakeDecision(bidID string, username string, decision bid.Decision) (*bid.Bid, error) {

	// Ensure the decision is valid
	if decision != bid.Approved && decision != bid.Rejected {
		return nil, ErrInvalidDecision
	}

	b, err := uc.bidRepo.GetBidByID(bidID)
	if err != nil {
		return nil, ErrBidNotFound
	}

	t, err := uc.tenderRepo.GetByID(b.TenderID)
	if err != nil {
		return nil, ErrNoSuchTender
	}

	// Verify if the user is responsible for the organization that owns the tender
	// responsible := uc.hasPermission()
	isResponsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, t.OrganizationID)
	if err != nil || !isResponsible {
		return nil, ErrUserNotResponsible
	}

	// If the decision is "Approved", close the tender
	if decision == bid.Approved {
		err = uc.tenderRepo.ChangeStatus(b.TenderID, username, tender.Closed)
		if err != nil {
			return nil, err
		}
	}

	err = uc.bidRepo.MakeDecision(bidID, username, decision)
	if err != nil {
		return nil, ErrInvalidDecision
	}
	return uc.GetBidByID(bidID)
}

func (uc *bidUseCase) GetBidsForTender(tenderID string, username string, limit, offset int) ([]*bid.Bid, error) {
	const op = "usecase.bid.GetBidsForTender:"

	// check if tenders exists
	tender, err := uc.tenderRepo.GetByID(tenderID)
	if err != nil {
		log.Println(op, err)
		return nil, ErrNoSuchTender
	}

	// check if user is responsible for organization who created tender
	responsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, tender.OrganizationID)
	if !responsible {
		log.Println(op, err)
		err = ErrUserNotResponsible
		return nil, err
	}

	var publishedBids []*bid.Bid

	allBids, err := uc.bidRepo.GetBids(tenderID, username, limit, offset)
	if err != nil {
		return nil, err
	}
	for _, b := range allBids {
		if b.Status == bid.Published {
			publishedBids = append(publishedBids, b)
		}
	}
	return publishedBids, nil
}

func (uc *bidUseCase) GetBidsByUser(username string, limit, offset int) ([]*bid.Bid, error) {
	const op = "useccase.bids.GetBidsByUser:"

	_, err := uc.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Println(op, err)
		return nil, ErrUserNotFound
	}

	return uc.bidRepo.GetBidsByUser(username, limit, offset)
}

func (uc *bidUseCase) GetBidStatus(bidID, username string) (bid.BidStatus, error) {
	const op = "usecase.Bid.GetBidStatus:"

	bid, err := uc.GetBidByID(bidID)
	if err != nil {
		log.Println(op, err)
		return "", err
	}

	user, err := uc.userRepo.GetUserByUsername(username)
	if user == nil {
		return "", ErrUserNotFound
	}
	if err != nil {
		return "", err
	}

	responsible := uc.hasPermission(bidID, username)
	if !responsible {
		return "", ErrUserNotResponsibleForBid
	}

	return bid.Status, nil
}

// hasPermission checks if the user has permission to act on the bid
// user from organization who placed a bid or from organization who created tender or the author
func (uc *bidUseCase) hasPermission(bidID string, username string) bool {

	const op = "usecase.bid.hasPermission"
	log.Println(op, bidID, username)

	b, err := uc.bidRepo.GetBidByID(bidID)
	if err != nil {
		log.Println(op, err)
		return false
	}
	user, err := uc.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Println(op, err)
		return false
	}

	// author
	if b.AuthorType == bid.User {
		return b.AuthorID == user.ID
	}
	// other responsible users
	if b.AuthorType == bid.Organization {
		// author id
		responsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, b.AuthorID)
		if err != nil {
			log.Println(op, err)
			return false
		}
		return responsible
	}
	return false
}
func updateBidStruct(a *bid.Bid, b bid.Bid) {
	vA := reflect.ValueOf(a).Elem()
	vB := reflect.ValueOf(b)

	for i := 0; i < vA.NumField(); i++ {
		fieldA := vA.Field(i)
		fieldB := vB.Field(i)

		if fieldA.CanSet() {
			switch fieldA.Kind() {
			case reflect.String:
				if fieldB.String() != "" {
					fieldA.SetString(fieldB.String())
				}
			case reflect.Int:
				if fieldB.Int() != 0 {
					fieldA.SetInt(fieldB.Int())
				}
			case reflect.Bool:
				fieldA.SetBool(fieldB.Bool())
			}
		}
	}
}
