package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/Bank-Thanapat-Developer/basic-redis/internal/domains"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/dto"
	"github.com/Bank-Thanapat-Developer/basic-redis/internal/entities"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepository domains.UserRepository
	jwtSecret      string
	jwtExpiry      time.Duration
}

func NewAuthUsecase(userRepo domains.UserRepository, jwtSecret string, jwtExpiryHours int) domains.AuthUsecase {
	return &authUsecase{
		userRepository: userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiry:      time.Duration(jwtExpiryHours) * time.Hour,
	}
}

func (u *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) error {
	if req.Username == "" || req.Password == "" {
		return errors.New("username and password are required")
	}

	existing, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entities.User{
		Username:  req.Username,
		Password:  string(hashed),
		CreatedAt: time.Now(),
	}
	return u.userRepository.Create(ctx, user)
}

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(u.jwtExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken: signed,
		TokenType:   "Bearer",
	}, nil
}
