package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"time"
)

var ErrNotFound = errors.New("not found")

type AuthSQLite struct {
	db *sql.DB
}

func NewAuthSQLite(db *sql.DB) *AuthSQLite {
	return &AuthSQLite{db: db}
}

func (a *AuthSQLite) RegisterUser(user *models.User) error {
	query := "INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)"
	result, err := a.db.Exec(query, user.Email, user.Username, user.PasswordHash)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(userID)

	return nil
}

func (a *AuthSQLite) GetUserByEmail(email string) (*models.User, error) {
	query := "SELECT id, email, username, password_hash FROM users WHERE email = ?"
	row := a.db.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (a *AuthSQLite) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT id, email, username, password_hash FROM users WHERE username = ?"
	row := a.db.QueryRow(query, username)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (a *AuthSQLite) GetUserByID(userID int) (*models.User, error) {
	query := "SELECT id, email, username, password_hash FROM users WHERE id = ?"
	row := a.db.QueryRow(query, userID)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (a *AuthSQLite) CreateSession(session *models.Session) error {
	query := "INSERT INTO sessions (user_id, token, expiry) VALUES (?, ?, ?)"
	result, err := a.db.Exec(query, session.UserID, session.Token, session.ExpiresAt)
	if err != nil {
		return err
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	session.ID = int(sessionID)

	return nil
}

func (a *AuthSQLite) GetSessionByToken(token string) (*models.Session, error) {
	query := "SELECT id, user_id, token, expiry FROM sessions WHERE token = ?"
	row := a.db.QueryRow(query, token)

	session := &models.Session{}
	err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return session, nil
}

func (a *AuthSQLite) DeleteAllSessionsByUserID(userID int) error {
	query := "DELETE FROM sessions WHERE user_id = ?"
	_, err := a.db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthSQLite) ClearExpiredSessions() error {
	query := "DELETE FROM sessions WHERE expiry < ?"
	result, err := a.db.Exec(query, time.Now())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("Cleared %d expired sessions\n", rowsAffected)
	return nil
}
