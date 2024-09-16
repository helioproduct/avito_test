package postgres

import (
	user "avito_api/internal/entities/user"
	uc "avito_api/internal/usecase"
	"database/sql"
	"log"
)

type postgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) uc.UserRepo {
	return &postgresUserRepo{
		DB: db,
	}
}

func (r *postgresUserRepo) CreateUser(user *user.User) error {
	const op = "postgres.user.CreateUser:"

	query := `
			INSERT INTO employee (username, first_name, last_name, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			RETURNING id
		`
	err := r.DB.QueryRow(query, user.Username, user.FirstName, user.LastName).Scan(&user.ID)
	if err != nil {
		log.Println(op, err)
		return err
	}
	return nil
}

func (r *postgresUserRepo) GetUserByID(id string) (*user.User, error) {
	const op = "postgres.user.GetUserByID:"

	var user user.User
	query := `
        SELECT id, username, first_name, last_name, created_at, updated_at
        FROM employee
        WHERE id = $1
    `
	row := r.DB.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepo) GetUserByUsername(username string) (*user.User, error) {
	const op = "postgres.user.GetUserByUsername:"
	var user user.User
	query := `
        SELECT id, username, first_name, last_name, created_at, updated_at
        FROM employee
        WHERE username = $1
    `
	row := r.DB.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println(op, err)
		return nil, err
	}
	return &user, nil
}

func (r *postgresUserRepo) GetAllUsers() ([]*user.User, error) {
	const op = "postgres.user.GetAllUser:"

	query := `
		SELECT id, username, first_name, last_name
		FROM employee
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Println(op, err)
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
