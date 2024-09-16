package postgres

import (
	"database/sql"
	"errors"
	"log"

	"avito_api/internal/entities/bid"
	uc "avito_api/internal/usecase"
)

type postgresBidRepository struct {
	DB *sql.DB
}

func NewPostgresBidRepo(db *sql.DB) uc.BidRepo {
	return &postgresBidRepository{DB: db}
}

func (r *postgresBidRepository) CreateBid(b *bid.Bid) (string, error) {
	const op = "posrgres.bid.CreateBid:"

	query := `
		INSERT INTO bids (name, description, tender_id, author_type, author_id, status, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id
	`

	var id string
	err := r.DB.QueryRow(query, b.Name, b.Description, b.TenderID, b.AuthorType, b.AuthorID, b.Status, b.Version).Scan(&id)
	if err != nil {
		log.Println(op, err)
		return "", err
	}

	return id, nil
}

func (r *postgresBidRepository) GetBidByID(bidID string) (*bid.Bid, error) {
	const op = "posrgres.bid.GetBidByID:"

	query := `
		SELECT id, name, description, tender_id, author_type, author_id, status, decision, version, created_at, updated_at
		FROM bids
		WHERE id = $1
	`

	var b bid.Bid
	err := r.DB.QueryRow(query, bidID).Scan(
		&b.ID, &b.Name, &b.Description, &b.TenderID, &b.AuthorType, &b.AuthorID,
		&b.Status, &b.Decision, &b.Version, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(op, err)
			return nil, errors.New("bid not found")
		}
		return nil, err
	}

	return &b, nil
}

// EditBid updates an existing bid
func (r *postgresBidRepository) EditBid(updatedBid *bid.Bid) error {
	const op = "posrgres.bid.EditBid:"

	query := `
		UPDATE bids
		SET name = $2, description = $3, version = version + 1, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.Exec(query, updatedBid.ID, updatedBid.Name, updatedBid.Description)
	if err != nil {
		log.Println(op, err)
	}
	return err
}

// ChangeStatus changes the status of a bid
func (r *postgresBidRepository) ChangeStatus(bidID, username string, newStatus bid.BidStatus) error {
	const op = "posrgres.bid.ChangeStatus:"

	query := `
		UPDATE bids
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.Exec(query, bidID, newStatus)
	if err != nil {
		log.Println(op, err)
	}
	return err
}

// ChangeDecision updates the decision on a bid (Approve/Reject)
func (r *postgresBidRepository) MakeDecision(bidID, username string, decision bid.Decision) error {
	const op = "posrgres.bid.ChangeDecision:"

	query := `
		UPDATE bids
		SET decision = $2, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.DB.Exec(query, bidID, decision)
	if err != nil {
		log.Println(op, err)
	}
	return err
}

func (r *postgresBidRepository) GetBids(tenderID string, username string, limit, offset int) ([]*bid.Bid, error) {
	const op = "posrgres.bid.GetBids:"

	query := `
		SELECT b.id, b.name, b.description, b.tender_id, b.author_type, b.author_id, b.status, b.decision, b.version, b.created_at, b.updated_at
		FROM bids b
		JOIN tender t ON b.tender_id = t.id
		JOIN organization_responsible org_res ON t.organization_id = org_res.organization_id
		JOIN employee e ON org_res.user_id = e.id
		WHERE b.tender_id = $1 AND e.username = $2
		ORDER BY b.created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.Query(query, tenderID, username, limit, offset)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	defer rows.Close()

	var bids []*bid.Bid
	for rows.Next() {
		var b bid.Bid
		err := rows.Scan(
			&b.ID, &b.Name, &b.Description, &b.TenderID, &b.AuthorType, &b.AuthorID,
			&b.Status, &b.Decision, &b.Version, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		bids = append(bids, &b)
	}

	return bids, nil
}

func (r *postgresBidRepository) GetBidsByUser(username string, limit, offset int) ([]*bid.Bid, error) {
	const op = "posrgres.bid.GetBidsByUser:"

	// only personal bids
	query := `
		SELECT b.id, b.name, b.description, b.tender_id, b.author_type, b.author_id, b.status, b.decision, b.version, b.created_at, b.updated_at
		FROM bids b
		JOIN employee e ON b.author_id = e.id
		WHERE e.username = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(query, username, limit, offset)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	defer rows.Close()

	var bids []*bid.Bid
	for rows.Next() {
		var b bid.Bid
		err := rows.Scan(
			&b.ID, &b.Name, &b.Description, &b.TenderID, &b.AuthorType, &b.AuthorID,
			&b.Status, &b.Decision, &b.Version, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		bids = append(bids, &b)
	}

	// organization's bids
	query = `
		SELECT b.id, b.name, b.description, b.tender_id, b.author_type, b.author_id, b.status, b.decision, b.version, b.created_at, b.updated_at
		FROM bids b
		JOIN organization_responsible org_resp ON b.author_id = org_resp.organization_id
		JOIN employee e ON e.id = org_resp.user_id
		
		WHERE e.username = $1
		ORDER BY b.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err = r.DB.Query(query, username, limit, offset)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	for rows.Next() {
		var b bid.Bid
		err := rows.Scan(
			&b.ID, &b.Name, &b.Description, &b.TenderID, &b.AuthorType, &b.AuthorID,
			&b.Status, &b.Decision, &b.Version, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		bids = append(bids, &b)
	}

	return bids, nil
}
