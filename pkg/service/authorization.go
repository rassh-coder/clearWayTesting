package service

import (
	"clearWayTest/pkg/models"
	"clearWayTest/pkg/repository"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthorization(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func (s *AuthService) SignIn(login, password, ip string) (string, error) {
	user, err := s.repos.GetUserByLogin(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("unauthorized")
		}
		log.Printf("Can't get user: %s", err)
		return "", err
	}

	if user == nil {
		return "", errors.New("unauthorized")
	}

	passHash := hashPassword(password)

	if user.PasswordHash != passHash {
		return "", errors.New("unauthorized")
	}

	token, err := s.generateToken(user.ID, ip)
	if err != nil {
		log.Printf("Can't generate token: %s", err)
		return "", err
	}

	return token, nil
}

// Hashing password
func hashPassword(password string) string {
	data := []byte(password)
	hash := md5.Sum(data)

	return hex.EncodeToString(hash[:])
}

// Create a new session record with returning session id(token)
func (s *AuthService) generateToken(uid uint, ip string) (string, error) {
	session := models.Session{
		UID: uid,
		Ip:  ip,
	}
	token, err := s.repos.SaveSession(&session)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetSessionByToken(token string) (*models.Session, error) {
	session, err := s.repos.GetSession(token)
	if err != nil {
		log.Printf("Can't get session: %s", err)
		return nil, err
	}

	return session, nil
}
