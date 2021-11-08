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

	result, err := statment.Exec(user.Name, user.Nick, user.Email, user.Password)
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

func (repository Users) Update(userID uint64, user models.User) error {

	statament, err := repository.db.Prepare("UPDATE users set name = ?, nick = ?, email = ? where id = ?")
	if err != nil {
		return err
	}
	defer statament.Close()

	if _, err = statament.Exec(user.Name, user.Nick, user.Email, userID); err != nil {
		return err
	}

	return nil
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

func (repository Users) Delete(ID uint64) error {

	statament, err := repository.db.Prepare("DELETE FROM users WHERE id = ?")

	if err != nil {
		return err
	}
	defer statament.Close()

	if _, err = statament.Exec(ID); err != nil {
		return err
	}

	return nil
}

func (repository Users) GetByEmail(email string) (models.User, error) {

	rows, err := repository.db.Query("SELECT id, password, created_at FROM users where email = ?", email)

	if err != nil {
		return models.User{}, err
	}

	defer rows.Close()

	var user models.User

	if rows.Next() {
		if err = rows.Scan(
			&user.ID,
			&user.Password,
		); err != nil {
			return models.User{}, err
		}
	}
	return user, nil
}

func (repository Users) FollowUser(userID uint64, followerID uint64) error {

	statement, err := repository.db.Prepare("INSERT IGNORE INTO followers(user_id, follower_id) VALUES (?, ?)")

	if err != nil {
		return err
	}

	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return err
	}

	return nil
}

func (repository Users) UnfollowUser(userID uint64, followerID uint64) error {

	statement, err := repository.db.Prepare("DELETE FROM followers WHERE user_id = ? AND follower_id = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return err
	}

	return nil
}

func (repository Users) GetFollowers(userID uint64) ([]models.User, error) {
	rows, err := repository.db.Query(`
	    SELECT u.id, u.name, u.nick, u.email, u.created_at 
		FROM users u inner join followers s on u.id = s.follower_id 
		WHERE s.user_id = ?
	`, userID)

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

func (repository Users) GetFollowing(userID uint64) ([]models.User, error) {
	rows, err := repository.db.Query(`
	    SELECT u.id, u.name, u.nick, u.email, u.created_at 
		FROM users u inner join followers s on u.id = s.follower_id 
		WHERE s.user_id = ?
	`, userID)

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

func (repository Users) GetPasswordByUserID(userID uint64) (string, error) {
	rows, err := repository.db.Query("SELECT password from users where id = ?")
	if err != nil {
		return "", err
	}

	defer rows.Close()

	var user models.User

	if rows.Next() {
		if err = rows.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil

}

func (repository Users) UpdatePassword(userID uint64, passoword string) error {
	statement, err := repository.db.Prepare("UPDATE users set password = ? where id = ? ")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(passoword, userID); err != nil {
		return err
	}
	return nil
}
