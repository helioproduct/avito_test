package repositories

import (
	"avito_api/internal/entities/organization"
	"avito_api/internal/entities/user"
	uc "avito_api/internal/usecase"
	"database/sql"
	"log"
	"time"
)

type postgresOrganizationRepo struct {
	DB *sql.DB
}

func NewPostgresOrganizationRepo(db *sql.DB) uc.OrganizationRepo {
	return &postgresOrganizationRepo{
		DB: db,
	}
}

// Create inserts a new organization into the database
func (r *postgresOrganizationRepo) Create(organization *organization.Organization) (string, error) {
	const op = "postgres.organizatiom.Create:"
	query := `
		INSERT INTO organization (name, description, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	currentTime := time.Now()
	organization.CreatedAt = currentTime
	organization.UpdatedAt = currentTime

	// Let the database generate the UUID and return it
	err := r.DB.QueryRow(query, organization.Name, organization.Description, organization.Type, organization.CreatedAt, organization.UpdatedAt).Scan(&organization.ID)
	if err != nil {
		log.Println(op, err)
		return "", err
	}

	return organization.ID, nil
}

// GetResponsibleUsers returns the list of users responsible for the organization
func (r *postgresOrganizationRepo) GetResponsibleUsers(organizationID string) ([]user.User, error) {
	const op = "postgres.organization.GetResponsibleUsers"

	query := `
		SELECT e.id, e.username, e.first_name, e.last_name, e.created_at, e.updated_at
		FROM employee e
		JOIN organization_responsible org_rep ON e.id = org_rep.user_id
		WHERE org_rep.organization_id = $1
	`

	rows, err := r.DB.Query(query, organizationID)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var user user.User
		err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// SetResponsibleUsers assigns a list of responsible users to an organization
func (r *postgresOrganizationRepo) SetResponsibleUsers(organizationID string, userIDs []string) error {
	const op = "postgtes.organization.SetResponsibleUsers:"

	deleteQuery := `DELETE FROM organization_responsible WHERE organization_id = $1`
	_, err := r.DB.Exec(deleteQuery, organizationID)
	if err != nil {
		log.Println(op, err)
		return err
	}

	insertQuery := `
		INSERT INTO organization_responsible (organization_id, user_id)
		VALUES ($1, $2)
	`
	for _, userID := range userIDs {
		_, err := r.DB.Exec(insertQuery, organizationID, userID)
		if err != nil {
			log.Println(op, err)
			return err
		}
	}

	return nil
}

func (r *postgresOrganizationRepo) GetByName(name string) (*organization.Organization, error) {
	const op = "postgtes.organization.GetByName:"

	query := `
		SELECT id, name, description, type, created_at, updated_at
		FROM organization
		WHERE name = $1
	`

	var org organization.Organization
	err := r.DB.QueryRow(query, name).Scan(&org.ID, &org.Name, &org.Description, &org.Type, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		log.Println(op, err)
		if err == sql.ErrNoRows {
			log.Println(op, err)
			return nil, nil
		}
		return nil, err
	}

	return &org, nil
}

func (r *postgresOrganizationRepo) IsUserResponsibleForOrganization(username, organizationID string) (bool, error) {
	const op = "postgtes.organization.IsUserResponsibleForOrganization:"

	query := `
	SELECT EXISTS(
		SELECT organization_responsible.user_id
		FROM organization_responsible
		JOIN (SELECT id FROM employee WHERE username = $1) AS sb
		ON sb.id = organization_responsible.user_id
		WHERE organization_responsible.organization_id = $2
	)`

	var exists bool
	err := r.DB.QueryRow(query, username, organizationID).Scan(&exists)
	if err != nil {
		log.Println(op, err)
		return false, err
	}
	log.Println(op, username, organizationID, exists)

	return exists, nil
}

func (r *postgresOrganizationRepo) GetAllOrganizations() ([]*organization.Organization, error) {
	const op = "postgtes.organization.GetAllOrganizations:"

	query := `
		SELECT id, name, description, type, created_at, updated_at
		FROM organization
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	defer rows.Close()

	var orgs []*organization.Organization
	for rows.Next() {
		var org organization.Organization
		err := rows.Scan(&org.ID, &org.Name, &org.Description, &org.Type, &org.CreatedAt, &org.UpdatedAt)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		orgs = append(orgs, &org)
	}

	return orgs, nil
}

func (r *postgresOrganizationRepo) GetUserOrganizationID(userID string) (string, error) {
	const op = "postgtes.organization.GetUserOrganizationID:"
	query := `
		SELECT organization_id
		FROM organization_responsible
		WHERE user_id = $1
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		log.Println(op, err)
		return "", err
	}
	defer rows.Close()

	var orgsID []string
	for rows.Next() {
		var organizationID string
		err := rows.Scan(&organizationID)
		if err != nil {
			log.Println(op, err)
			return "", err
		}
		orgsID = append(orgsID, organizationID)
	}
	return orgsID[0], nil
}
