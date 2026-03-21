package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/config"
	"github.com/YannisMuminov/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	UserID      int64    `json:"user_id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`

	jwt.RegisteredClaims
}

type authService struct {
	repo domain.AuthRepository
	cfg  *config.JWTConfig
}

func NewAuthService(repo domain.AuthRepository, cfg *config.JWTConfig) domain.AuthService {
	return &authService{
		repo: repo,
		cfg:  cfg,
	}
}

func (a *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	existing, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperror.ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		RoleID:       3,
	}

	if err := a.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.TokenPair, error) {
	user, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, apperror.ErrInvalidCreds
	}

	if !user.IsActive {
		return nil, apperror.ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperror.ErrInvalidCreds
	}

	role, permission, err := a.repo.GetRolePermissions(ctx, user.RoleID)

	if err != nil {
		return nil, fmt.Errorf("load permissions: %w", err)
	}

	return a.generateTokenPair(user, role, permission)
}

func (a *authService) Refresh(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenPair, error) {
	token, err := jwt.ParseWithClaims(
		req.RefreshToken,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(a.cfg.Secret), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, apperror.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject == "" {
		return nil, apperror.ErrUnauthorized
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return nil, apperror.ErrUnauthorized
	}

	user, err := a.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("refresh get user: %w", err)
	}

	if user == nil || !user.IsActive {
		return nil, apperror.ErrUnauthorized
	}

	role, permission, err := a.repo.GetRolePermissions(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("refresh load permissions: %w", err)
	}

	return a.generateTokenPair(user, role, permission)

}

func (a *authService) Me(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := a.repo.GetUserByID(ctx, userID)

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.ErrInvalidCreds
	}

	role, permission, err := a.repo.GetRolePermissions(ctx, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("load permissions: %w", err)
	}

	user.Role = role
	user.Permissions = permission
	return user, nil
}

func (a *authService) generateTokenPair(user *domain.User, role *domain.Role, permissions []domain.Permission) (*domain.TokenPair, error) {
	perms := make([]string, len(permissions))

	for i, p := range permissions {
		perms[i] = p.Name
	}

	roleName := ""

	if role != nil {
		roleName = role.Name
	}

	now := time.Now()

	accessClaims := JWTClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        roleName,
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(
				time.Duration(a.cfg.ExpireMinutes) * time.Minute,
			)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(a.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject: fmt.Sprintf("%d", user.ID),
		ExpiresAt: jwt.NewNumericDate(now.Add(
			time.Duration(a.cfg.RefreshExpireHours) * time.Minute,
		)),
		IssuedAt: jwt.NewNumericDate(now),
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(a.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
