package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	authDomain "invitation-api/internal/domain/auth"
	userDomain "invitation-api/internal/domain/user"
	authRepo "invitation-api/internal/repository/auth"
	userRepo "invitation-api/internal/repository/user"
	userService "invitation-api/internal/service/user"
	"invitation-api/pkg/config"
	"invitation-api/pkg/utils"

	"github.com/google/uuid"
)

// Service defines the interface for authentication business operations
type Service interface {
	Register(req *authDomain.RegisterRequest) (*authDomain.TokenResponse, error)
	Login(req *authDomain.LoginRequest) (*authDomain.TokenResponse, error)
	RefreshToken(req *authDomain.RefreshTokenRequest) (*authDomain.TokenResponse, error)
	Logout(token string) error
	ValidateToken(token string) (*utils.JWTClaims, error)
	RevokeAllUserTokens(userID uuid.UUID) error
}

// service implements the Service interface
type service struct {
	authRepo    authRepo.Repository
	userRepo    userRepo.Repository
	userService userService.Service
	config      *config.Config
}

// NewService creates a new auth service instance
func NewService(authRepo authRepo.Repository, userRepo userRepo.Repository, config *config.Config) Service {
	userSvc := userService.NewService(userRepo)

	return &service{
		authRepo:    authRepo,
		userRepo:    userRepo,
		userService: userSvc,
		config:      config,
	}
}

// Register handles user registration and returns tokens
func (s *service) Register(req *authDomain.RegisterRequest) (*authDomain.TokenResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Convert role string to UserRole
	var role userDomain.UserRole
	switch req.Role {
	case "admin":
		role = userDomain.RoleAdmin
	case "user":
		role = userDomain.RoleUser
	default:
		return nil, errors.New("invalid role")
	}

	// Create user service and register user
	fullName := req.FirstName + " " + req.LastName
	user, err := s.userService.RegisterUser(req.Username, req.Email, req.Password, fullName, role)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// Generate tokens
	return s.generateTokens(user)
}

// Login handles user login and returns tokens
func (s *service) Login(req *authDomain.LoginRequest) (*authDomain.TokenResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Get identifier (username or email)
	identifier := req.GetIdentifier()

	// Validate credentials - identifier can be username or email
	user, err := s.userService.ValidateCredentials(identifier, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	// Generate tokens
	return s.generateTokens(user)
}

// RefreshToken handles token refresh
func (s *service) RefreshToken(req *authDomain.RefreshTokenRequest) (*authDomain.TokenResponse, error) {
	// Get refresh token from database
	token, err := s.authRepo.GetTokenByValue(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if token is valid
	if !token.IsValid() {
		return nil, errors.New("refresh token is invalid or expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(token.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActiveAndNotDeleted() {
		return nil, errors.New("user account is inactive")
	}

	// Revoke old refresh token
	s.authRepo.RevokeToken(token.ID)

	// Generate new tokens
	return s.generateTokens(user)
}

// Logout handles user logout by revoking tokens
func (s *service) Logout(token string) error {
	// Validate token and get claims
	claims, err := s.ValidateToken(token)
	if err != nil {
		return err
	}

	// Revoke all user tokens
	return s.RevokeAllUserTokens(claims.UserID)
}

// ValidateToken validates a JWT token
func (s *service) ValidateToken(token string) (*utils.JWTClaims, error) {
	claims, err := utils.ValidateToken(token, s.config.JWT.SecretKey)
	if err != nil {
		return nil, err
	}

	// Check if token exists in database and is not revoked
	dbToken, err := s.authRepo.GetTokenByValue(token)
	if err != nil {
		return nil, errors.New("token not found in database")
	}

	if !dbToken.IsValid() {
		return nil, errors.New("token is revoked or expired")
	}

	return claims, nil
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (s *service) RevokeAllUserTokens(userID uuid.UUID) error {
	return s.authRepo.RevokeAllUserTokens(userID)
}

// generateTokens generates access and refresh tokens for a user
func (s *service) generateTokens(user *userDomain.User) (*authDomain.TokenResponse, error) {
	// Get full name from profile if available, fallback to username
	fullName := user.Username
	if user.Profile != nil && user.Profile.FullName != "" {
		fullName = user.Profile.FullName
	}

	// Generate access token
	accessToken, err := utils.GenerateTokenWithName(
		user.ID,
		user.Email,
		string(user.Role),
		fullName,
		s.config.JWT.SecretKey,
		s.config.JWT.AccessTokenTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save tokens to database
	accessTokenModel := authDomain.NewToken(
		user.ID,
		accessToken,
		authDomain.TokenTypeAccess,
		time.Now().Add(s.config.JWT.AccessTokenTTL),
	)
	if err := s.authRepo.CreateToken(accessTokenModel); err != nil {
		return nil, fmt.Errorf("failed to save access token: %w", err)
	}

	refreshTokenModel := authDomain.NewToken(
		user.ID,
		refreshToken,
		authDomain.TokenTypeRefresh,
		time.Now().Add(s.config.JWT.RefreshTokenTTL),
	)
	if err := s.authRepo.CreateToken(refreshTokenModel); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &authDomain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.JWT.AccessTokenTTL.Seconds()),
	}, nil
}

// generateRefreshToken generates a random refresh token
func (s *service) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
