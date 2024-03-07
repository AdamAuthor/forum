package service

import (
	"errors"
	"forum/internal/models"
	"forum/internal/repository"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	authService := &AuthService{repo: repo}
	go authService.ClearExpiredSessionsPeriodically(24 * time.Hour)
	return authService
}

func (a *AuthService) CreateUser(user *models.User) error {
	existingUser, err := a.repo.GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	if existingUser != nil {
		return errors.New("email already taken")
	}

	existingUser, err = a.repo.GetUserByUsername(user.Username)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	if existingUser != nil {
		return errors.New("username already taken")
	}

	hash, err := a.generatePasswordHash(user.PasswordHash)
	if err != nil {
		return errors.New("error with create password hash")
	}
	user.PasswordHash = hash

	err = a.repo.RegisterUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) GetUserByID(userID int) (*models.User, error) {
	return a.repo.GetUserByID(userID)
}

func (a *AuthService) CreateSession(username, password string) (*models.Session, int, error) {
	user, err := a.repo.GetUserByUsername(username)
	if err != nil {
		return nil, 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, 0, err
	}

	a.repo.DeleteAllSessionsByUserID(user.ID)

	token, err := uuid.NewV4()
	if err != nil {
		return nil, 0, err
	}
	expiry := time.Now().Add(time.Hour * 24)

	session := &models.Session{
		UserID:    user.ID,
		Token:     token.String(),
		ExpiresAt: expiry,
	}

	err = a.repo.CreateSession(session)
	if err != nil {
		return nil, 0, err
	}

	return session, session.UserID, err
}

func (a *AuthService) GetSessionByToken(token string) (*models.Session, error) {
	session, err := a.repo.GetSessionByToken(token)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (a *AuthService) ClearExpiredSessionsPeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		err := a.repo.ClearExpiredSessions()
		if err != nil {
			log.Printf("Failed to clear expired sessions: %v", err)
		}
	}
}

func (a *AuthService) generatePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
