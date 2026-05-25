package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jellydator/validation"
	"github.com/wreckitral/production-backend-go/internal/apperr"
	"github.com/wreckitral/production-backend-go/internal/model"
	"github.com/wreckitral/production-backend-go/internal/platform/config"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repo
	jwt  config.JWT
}

func NewService(repo *Repo, jwt config.JWT) *Service {
	return &Service{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (model.User, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required),
		validation.Field(&req.Password, validation.Required, validation.Length(8, 128)),
	); err != nil {
		return model.User{}, fmt.Errorf("validate: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user := model.User{
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return created, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required),
		validation.Field(&req.Password, validation.Required),
	); err != nil {
		return LoginResponse{}, fmt.Errorf("validate: %w", err)
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return LoginResponse{}, apperr.ErrUnauthorized
		}
		return LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return LoginResponse{}, apperr.ErrUnauthorized
	}

	expiresAt := time.Now().Add(s.jwt.TTL)
	token, err := s.signToken(user.ID, expiresAt)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("sign token: %w", err)
	}

	return LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

type Claims struct {
	jwt.RegisteredClaims
}

func (s *Service) signToken(userID uuid.UUID, expiresAt time.Time) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    s.jwt.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwt.Secret))
}
