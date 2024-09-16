package usecase

import (
	"avito_api/internal/entities/bid"
	"log"
)

// wrapper over base usecase where submit_decision can
type advancedBidUseCase struct {
	baseUseCase BidUseCase
	bidRepo     AdvancedBidRepo
	orgRepo     OrganizationRepo
	tenderRepo  TenderRepo
	quorum      int
}

type AdvancedBidRepo interface {
	BidRepo
	GetDecisionCount(bidId string, decision bid.Decision) (int, error)
}

func NewAdvancedBidUseCase(bidUseCase BidUseCase,
	bidRepo AdvancedBidRepo,
	orgRepo OrganizationRepo,
	tenderRepo TenderRepo, quorum int) BidUseCase {

	return &advancedBidUseCase{
		baseUseCase: bidUseCase,
		bidRepo:     bidRepo,
		orgRepo:     orgRepo,
		quorum:      quorum,
		tenderRepo:  tenderRepo,
	}
}

// create also sets aproves and rejects count to zeros
// advanced change decision method based on quorum
func (uc *advancedBidUseCase) MakeDecision(bidID string, username string, decision bid.Decision) (*bid.Bid, error) {
	const op = "usecase.advancedBid.ChangeDecision:"

	b, err := uc.bidRepo.GetBidByID(bidID)
	log.Println(op, "bid", b.ID)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	t, err := uc.tenderRepo.GetByID(b.TenderID)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	// user from organization that created a tender
	isResponsible, err := uc.orgRepo.IsUserResponsibleForOrganization(username, t.OrganizationID)
	if err != nil || !isResponsible {
		return nil, ErrUserNotResponsible
	}

	// nothing changes
	if decision == bid.Rejected {
		return uc.baseUseCase.MakeDecision(bidID, username, decision)
	}

	// save decision to repo
	err = uc.bidRepo.MakeDecision(bidID, username, decision)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	// updated approve count
	approveCount, err := uc.bidRepo.GetDecisionCount(bidID, decision)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	responsibleUsers, err := uc.orgRepo.GetResponsibleUsers(t.OrganizationID)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}

	quorum := min(uc.quorum, len(responsibleUsers))
	log.Println("current quorum is ", quorum)

	if approveCount >= quorum {
		// closes tender
		return uc.baseUseCase.MakeDecision(bidID, username, decision)
	}
	return b, nil
}

/*
	using underlying base usecase
*/

func (uc *advancedBidUseCase) Create(b *bid.Bid) (*bid.Bid, error) {
	return uc.baseUseCase.Create(b)
}

func (uc *advancedBidUseCase) ChangeStatus(bidID string, username string, newStatus bid.BidStatus) (*bid.Bid, error) {
	return uc.baseUseCase.ChangeStatus(bidID, username, newStatus)
}

func (uc *advancedBidUseCase) Edit(updatedBid *bid.Bid, username string) (*bid.Bid, error) {
	return uc.baseUseCase.Edit(updatedBid, username)
}

func (uc *advancedBidUseCase) GetBidByID(bidID string) (*bid.Bid, error) {
	return uc.baseUseCase.GetBidByID(bidID)
}

func (uc *advancedBidUseCase) GetBidsForTender(tenderID string, username string, limit, offset int) ([]*bid.Bid, error) {
	return uc.baseUseCase.GetBidsForTender(tenderID, username, limit, offset)
}

func (uc *advancedBidUseCase) GetBidsByUser(username string, limit, offset int) ([]*bid.Bid, error) {
	return uc.baseUseCase.GetBidsByUser(username, limit, offset)
}

func (uc *advancedBidUseCase) GetBidStatus(bidID, username string) (bid.BidStatus, error) {
	return uc.baseUseCase.GetBidStatus(bidID, username)
}
