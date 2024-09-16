package posgres

import (
	"avito_api/internal/entities/tender"
	uc "avito_api/internal/usecase"
	"database/sql"
	"log"
)

type postgresTenderRepo struct {
	DB *sql.DB
}

func NewPostgresTenderRepo(db *sql.DB) uc.TenderRepo {
	return &postgresTenderRepo{
		DB: db,
	}
}
func (r *postgresTenderRepo) CreateTender(tender *tender.Tender, username string) (string, error) {
	const op = "postgres.tender.CreateTender:"

	query := `
		INSERT INTO tender (name, description, serviceType, organization_id, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`
	err := r.DB.QueryRow(query, tender.Name, tender.Description, tender.ServiceType, tender.OrganizationID, tender.CreatorID).Scan(&tender.ID)
	if err != nil {
		log.Println(op, err)
		return "", err
	}
	return tender.ID, nil
}

func (r *postgresTenderRepo) GetStatus(tenderID string) (tender.StatusType, error) {
	const op = "postgres.tender.GetStatus:"

	query := `SELECT status FROM tender WHERE id = $1`
	var status tender.StatusType
	err := r.DB.QueryRow(query, tenderID).Scan(&status)
	if err != nil {
		log.Println(op, err)
		return "", err
	}
	return status, nil
}

func (r *postgresTenderRepo) GetByID(tenderID string) (*tender.Tender, error) {
	const op = "postgres.tender.GetByID:"
	query := `
		SELECT id, name, description, serviceType, organization_id, creator_id, status, current_version, created_at, updated_at
		FROM tender
		WHERE id = $1
	`
	var tender tender.Tender
	err := r.DB.QueryRow(query, tenderID).Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.OrganizationID, &tender.CreatorID, &tender.Status, &tender.CurrentVersion, &tender.CreatedAt, &tender.UpdatedAt)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return &tender, nil
}

func (r *postgresTenderRepo) GetVersionByID(tenderID string, version int) (*tender.Tender, error) {
	const op = "postgres.tender.GetVersionByID:"

	query := `
		SELECT name, description, serviceType
		FROM tender_versions
		WHERE tender_id = $1 AND version = $2
	`
	var oldTender tender.Tender
	err := r.DB.QueryRow(query, tenderID, version).Scan(&oldTender.Name, &oldTender.Description, &oldTender.ServiceType)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return &oldTender, nil
}

func (r *postgresTenderRepo) ChangeStatus(tenderID string, username string, newStatus tender.StatusType) error {
	const op = "postgres.tender.ChangeStatus"

	query := `UPDATE tender SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.DB.Exec(query, tenderID, newStatus)
	if err != nil {
		log.Println(op, err)
	}
	return err
}

func (r *postgresTenderRepo) EditTender(updatedTender *tender.Tender, username string) error {
	const op = "postgres.tender.Editender:"

	query := `
		UPDATE tender
		SET name = $2, description = $3, serviceType = $4, updated_at = NOW(), current_version = current_version + 1
		WHERE id = $1
	`
	_, err := r.DB.Exec(query, updatedTender.ID, updatedTender.Name, updatedTender.Description, updatedTender.ServiceType)
	if err != nil {
		log.Println(op, err)
	}
	return err
}

// limit=-1 returns all published tenders
func (r *postgresTenderRepo) GetPublishedTenders(limit, offset int, serviceType string) ([]*tender.Tender, error) {
	const op = "postgres.tender.GetPublishedTenders:"

	query := `
		SELECT id, name, description, servicetype, organization_id, creator_id, current_version, created_at, updated_at, status
		FROM tender
		WHERE status = 'Published'
	`

	var args []interface{}
	var limitArg interface{}

	if limit == -1 {
		limitArg = "NULL"
	} else {
		limitArg = limit
	}
	if serviceType != "" {
		query += " AND servicetype = $1"
		query += " LIMIT $2 OFFSET $3"
		args = append(args, serviceType)
	} else {
		query += " LIMIT $1 OFFSET $2"
	}
	args = append(args, limitArg, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	defer rows.Close()

	var tenders []*tender.Tender
	for rows.Next() {
		var tender tender.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.OrganizationID, &tender.CreatorID, &tender.CurrentVersion, &tender.CreatedAt, &tender.UpdatedAt, &tender.Status)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

func (r *postgresTenderRepo) GetTendersByResponsibleUser(username string, limit, offset int, serviceType string) ([]*tender.Tender, error) {
	const op = "postgres.tender.GetTendersByResponsibleUser:"

	query := `
		SELECT t.id, t.name, t.description, t.servicetype, t.organization_id, t.creator_id, t.current_version, t.created_at, t.updated_at
		FROM tender t
		JOIN organization_responsible org_res ON t.organization_id = org_res.organization_id
		JOIN employee e ON org_res.user_id = e.id
		WHERE e.username = $1
	`

	args := []interface{}{username, limit, offset}

	if serviceType != "" {
		query += " AND t.servicetype = $4"
		args = append(args, serviceType)
	}

	query += " LIMIT $2 OFFSET $3"
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	defer rows.Close()

	var tenders []*tender.Tender
	for rows.Next() {
		var tender tender.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.OrganizationID, &tender.CreatorID, &tender.CurrentVersion, &tender.CreatedAt, &tender.UpdatedAt)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}

// returns tenders for which user is responsible for
func (r *postgresTenderRepo) GetOwnedTenders(limit, offset int, username string) ([]*tender.Tender, error) {
	const op = "postgres.tender.GetOwnedTenders:"

	query := `
		SELECT t.id, t.name, t.description, t.servicetype, t.organization_id, t.creator_id, t.current_version, t.created_at, t.updated_at, t.status
		FROM tender t
		JOIN organization_responsible org_res ON t.organization_id = org_res.organization_id
		JOIN employee e ON org_res.user_id = e.id
		WHERE e.username = $1
	`

	args := []interface{}{username}

	if limit != -1 {
		query += " LIMIT $2 "
		query += " OFFSET $3 "
		args = append(args, limit, offset)
	} else {
		query += " OFFSET $2 "
		args = append(args, offset)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	defer rows.Close()

	var tenders []*tender.Tender
	for rows.Next() {
		var tender tender.Tender
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.ServiceType, &tender.OrganizationID, &tender.CreatorID, &tender.CurrentVersion, &tender.CreatedAt, &tender.UpdatedAt, &tender.Status)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		tenders = append(tenders, &tender)
	}

	return tenders, nil
}
