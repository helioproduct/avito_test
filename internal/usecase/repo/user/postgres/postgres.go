package postgres

import (
	user "avito_api/internal/entities/user"
	uc "avito_api/internal/usecase"
	"database/sql"
	"log"
)

type postgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) uc.UserRepository {
	return &postgresUserRepository{
		DB: db,
	}
}

func (r *postgresUserRepository) CreateUser(user *user.User) error {
	// Prepare the SQL query
	query := `
			INSERT INTO employee (username, first_name, last_name, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			RETURNING id
		`
	err := r.DB.QueryRow(query, user.Username, user.FirstName, user.LastName).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresUserRepository) GetUserByID(id string) (*user.User, error) {
	var user user.User
	query := `
        SELECT id, username, first_name, last_name, created_at, updated_at
        FROM employee
        WHERE id = $1
    `
	row := r.DB.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) GetUserByUsername(username string) (*user.User, error) {
	var user user.User
	query := `
        SELECT id, username, first_name, last_name, created_at, updated_at
        FROM employee
        WHERE username = $1
    `
	row := r.DB.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepository) GetAllUsers() ([]*user.User, error) {
	const op = "postgres.user.GetAllUser:"

	query := `
		SELECT id, username, first_name, last_name
		FROM employee
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var user user.User
		err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName)
		if err != nil {
			log.Println(op, err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
