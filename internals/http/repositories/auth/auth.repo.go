package repositories

import (
	"context"
	"time"

	types "github.com/omniflare/go_starter/internals/http/types/auth"
	models "github.com/omniflare/go_starter/internals/models/user"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func (r *AuthRepository) Register(ctx context.Context, registerData *types.AuthRegister) (*models.User, error) {
	user := &models.User{
		Name:     registerData.Name,
		Email:    registerData.Email,
		Password: registerData.Password,
	}

	res := r.db.Model(&models.User{}).Create(user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, query interface{}, args ...interface{}) (*models.User, error) {
	user := &models.User{}

	if res := r.db.Model(user).Where(query, args...).First(user); res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, userID uint, updates map[string]interface{}) error {
    result := r.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}
func (r *AuthRepository) CreateVerificationToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
    updates := map[string]interface{}{
        "verification_token": token,
        "verification_token_expires_at": expiresAt,
    }
    return r.UpdateUser(ctx, userID, updates)
}

func (r *AuthRepository) VerifyEmail(ctx context.Context, token string) (*models.User, error) {
    return r.GetUser(ctx, "verification_token = ? AND verification_token_expires_at > ?", token, time.Now())
}

func NewAuthRepository(db *gorm.DB) *AuthRepository{
	return &AuthRepository{db: db}
}