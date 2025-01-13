package models

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	types "github.com/omniflare/go_starter/internals/http/types/auth"
)

type User struct {
	ID                         int
	Name                       string
	Email                      string
	Password                   string
	IsVerified                 bool
	ResetPasswordToken         string
	ResetPasswordExpiresAt     time.Time
	VerificationToken          string
	VerificationTokenExpiresAt time.Time
}

type AuthRepository interface {
	Register(ctx context.Context, registerData *types.AuthRegister) (*User, error)
	GetUser(ctx context.Context, query interface{}, args ...interface{}) (*User, error)
	UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error
	VerifyEmail(ctx context.Context, token string) (*User, error)
	CreateVerificationToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
}
type AuthService interface {
	Login(ctx context.Context, loginData *types.AuthLogin) (string, *User, error)
	Register(ctx context.Context, registerData *types.AuthRegister) (string, *User, error)
	Logout(ctx context.Context) error
	VerifyEmail(c *gin.Context, code string) error
	ForgotPassword(c *gin.Context, userEmail string) error
	ResetPassword(c *gin.Context, token, newPassword string)
	GetMe(c *gin.Context, userID uint) (*User, error)
	ResendVerification(c *gin.Context, userEmail string) error
}
