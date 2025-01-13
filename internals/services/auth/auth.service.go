package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/omniflare/go_starter/internals/config"
	"github.com/omniflare/go_starter/internals/http/email"
	types "github.com/omniflare/go_starter/internals/http/types/auth"
	models "github.com/omniflare/go_starter/internals/models/user"
	utils "github.com/omniflare/go_starter/internals/utils/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repository models.AuthRepository
}

func (a *AuthService) Login(c *gin.Context, loginData *types.AuthLogin) (string, *models.User, error) {
	user, err := a.repository.GetUser(c, "email = ?", loginData.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("invalid credentials")
		}
		return "", nil, err
	}

	if !utils.MatchHash(loginData.Password, user.Password) {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	jwt_val := config.NewEnvConfig().JWT_SECRET
	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, jwt_val)
	if err != nil {
		return "", nil, err
	}
	c.SetCookie(
		"token",
		token,
		3600*24, // 24 hours in seconds
		"/",
		"",
		false, // PROD : set = true
		true,  // HTTP only
	)

	return token, user, nil
}

func (a *AuthService) Register(c *gin.Context, registerData *types.AuthRegister) (string, *models.User, error) {
	if !utils.ValidateEmail(registerData.Email) {
		return "", nil, fmt.Errorf("please, provide a valid email to register")
	}

	if err := utils.ValidatePassword(registerData.Password); err != nil {
		return "", nil, err
	}

	if _, err := a.repository.GetUser(c, "email = ?", registerData.Email); !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil, fmt.Errorf("the user email is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	registerData.Password = string(hashedPassword)

	user, err := a.repository.Register(c, registerData)
	if err != nil {
		return "", nil, err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", nil, fmt.Errorf("failed to generate verification token: %w", err)
	}
	verificationCode := utils.GenerateOTP()
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := a.repository.CreateVerificationToken(c, uint(user.ID), verificationCode, expiresAt); err != nil {
		return "", nil, fmt.Errorf("failed to create verification token: %w", err)
	}
	if err := email.SendVerificationEmail(user.Email, verificationCode); err != nil {
		c.Error(fmt.Errorf("failed to send verification email: %w", err))
	}
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 200).Unix(),
	}
	jwt_val := config.NewEnvConfig().JWT_SECRET
	token, err := utils.GenerateJWT(claims, jwt.SigningMethodHS256, jwt_val)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (a *AuthService) Logout(c *gin.Context) error {
	c.SetCookie("token", "", -1, "/", "", false, true)
	return nil
}

func (a *AuthService) VerifyEmail(c *gin.Context, code string) error {
	if len(code) != 6 {
		return fmt.Errorf("invalid verification code format")
	}
	user, err := a.repository.GetUser(c, "verification_token = ? AND verification_token_expires_at > ?", code, time.Now())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invalid or expired verification code")
		}
		return err
	}

	updates := map[string]interface{}{
		"is_verified":                   true,
		"verification_token":            nil,
		"verification_token_expires_at": nil,
	}

	if err := a.repository.UpdateUser(c, uint(user.ID), updates); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if err := email.SendWelcomeEmail(user.Email, user.Name); err != nil {
		c.Error(fmt.Errorf("failed to send welcome email: %w", err))
	}

	return nil
}

func (a *AuthService) ForgotPassword(c *gin.Context, userEmail string) error {
	user, err := a.repository.GetUser(c, "email = ?", userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}
	resetToken := hex.EncodeToString(tokenBytes)
	resetExpires := time.Now().Add(time.Hour)

	updates := map[string]interface{}{
		"reset_password_token":      resetToken,
		"reset_password_expires_at": resetExpires,
	}

	if err := a.repository.UpdateUser(c, uint(user.ID), updates); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	clientUrl := config.NewEnvConfig().CLIENT_URL
	resetURL := fmt.Sprintf("%s/api/auth/reset-password/v1/%s", clientUrl, resetToken)
	if err := email.SendPasswordResetEmail(user.Email, resetURL); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

func (a *AuthService) ResetPassword(c *gin.Context, token, newPassword string) error {
	user, err := a.repository.GetUser(c, "reset_password_token = ? AND reset_password_expires_at > ?", token, time.Now())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invalid or expired reset token")
		}
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	updates := map[string]interface{}{
		"password":                  string(hashedPassword),
		"reset_password_token":      nil,
		"reset_password_expires_at": nil,
	}

	if err := a.repository.UpdateUser(c, uint(user.ID), updates); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if err := email.SendResetSuccessEmail(user.Email); err != nil {
		c.Error(fmt.Errorf("failed to send reset success email: %w", err))
	}

	return nil
}

func (a *AuthService) GetMe(c *gin.Context, userID uint) (*models.User, error) {
	user, err := a.repository.GetUser(c, "id = ?", userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (a *AuthService) ResendVerification(c *gin.Context, userEmail string) error {
	user, err := a.repository.GetUser(c, "email = ?", userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	if user.IsVerified {
		return fmt.Errorf("email is already verified")
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}
	verificationToken := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := a.repository.CreateVerificationToken(c, uint(user.ID), verificationToken, expiresAt); err != nil {
		return fmt.Errorf("failed to create verification token: %w", err)
	}

	return email.SendVerificationEmail(user.Email, verificationToken)
}

func NewAuthService(repository models.AuthRepository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}
