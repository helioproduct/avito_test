package postgres

import (
	"avito_api/internal/entities/bid"
	uc "avito_api/internal/usecase"

	"database/sql"
	"log"
)

type postgresAdvancedBidRepo struct {
	DB       *sql.DB
	baseRepo uc.BidRepo
}

func NewPostgresBidAdvancedRepo(db *sql.DB, baseRepo uc.BidRepo) uc.AdvancedBidRepo {
	return &postgresAdvancedBidRepo{
		DB:       db,
		baseRepo: baseRepo,
	}
}

/*
	using underlying base repo
*/

func (r *postgresAdvancedBidRepo) GetDecisionCount(bidId string, decision bid.Decision) (int, error) {
	const op = "postgres.AdvancedBid.GetDecisionCount:"

	query := `
		SELECT COUNT(*)
		FROM decisions
		WHERE decision_value = $1 AND bid_id = $2
	`

	var count int
	err := r.DB.QueryRow(query, decision, bidId).Scan(&count)
	if err != nil {
		log.Println(op, err)
		return 0, err
	}
	return count, nil
}

func (r *postgresAdvancedBidRepo) MakeDecision(bidID string, username string, decision bid.Decision) error {
	const op = "postgres.advancedBid.ChangeDecison:"
	query := `
		INSERT INTO decisions (bid_id, username, decision_value)
		VALUES ($1, $2, $3)
		ON CONFLICT (bid_id, username, decision_value) DO NOTHING;
	`
	_, err := r.DB.Exec(query, bidID, username, decision)
	if err != nil {
		log.Println(op, err)
		return err
	}
	return r.baseRepo.MakeDecision(bidID, username, decision)
}

/*
	using underlying base repo
*/

func (r *postgresAdvancedBidRepo) ChangeStatus(bidID string, username string, newStatus bid.BidStatus) error {
	return r.baseRepo.ChangeStatus(bidID, username, newStatus)
}

func (r *postgresAdvancedBidRepo) CreateBid(bid *bid.Bid) (string, error) {
	return r.baseRepo.CreateBid(bid)
}

func (r *postgresAdvancedBidRepo) EditBid(updatedBid *bid.Bid) error {
	return r.baseRepo.EditBid(updatedBid)
}

func (r *postgresAdvancedBidRepo) GetBidByID(bidID string) (*bid.Bid, error) {
	return r.baseRepo.GetBidByID(bidID)
}

func (r *postgresAdvancedBidRepo) GetBids(tenderID string, username string, limit int, offset int) ([]*bid.Bid, error) {
	return r.baseRepo.GetBids(tenderID, username, limit, offset)
}

func (r *postgresAdvancedBidRepo) GetBidsByUser(username string, limit int, offset int) ([]*bid.Bid, error) {
	return r.baseRepo.GetBidsByUser(username, limit, offset)
}
