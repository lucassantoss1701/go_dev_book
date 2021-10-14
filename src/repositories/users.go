package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type Users struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *Users {
	return &Users{db}
}

func (repository Users) Create(user models.User) (uint64, error) {

	statment, err := repository.db.Prepare("INSERT INTO users (name, nick, email, password) values (?, ?, ?, ?)")

	if err != nil {
		return 0, nil
	}

	defer statment.Close()

	result, err := statment.Exec(user.Name, user.Nick, user.Email, user.Passoword)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}

func (repository Users) GetByNameOrNick(nameOrNick string) ([]models.User, error) {

	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick)
	rows, err := repository.db.Query(`SELECT id, name, nick, email, created_at FROM users where name LIKE ? OR nick LIKE ?`, nameOrNick, nameOrNick)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Created_at,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository Users) GetUserById(ID uint64) (models.User, error) {

	rows, err := repository.db.Query("SELECT id, name, nick, email, created_at FROM users where id = ?", ID)

	if err != nil {
		return models.User{}, err
	}

	defer rows.Close()

	var user models.User

	if rows.Next() {
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nick,
			&user.Email,
			&user.Created_at,
		); err != nil {
			return models.User{}, err
		}
	}
	return user, nil
}
