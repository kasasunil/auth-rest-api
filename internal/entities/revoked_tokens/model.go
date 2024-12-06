package revoked_tokens

import (
	"github.com/kasasunil/auth-rest-api/internal/database"
	"gorm.io/gorm"
)

type RevokedToken struct {
	db *database.Db
}

type RevokedTokenModel struct {
	gorm.Model
	Token string `gorm:"primaryKey"`
}

func New(db *database.Db) *RevokedToken {
	return &RevokedToken{
		db: db,
	}
}

// CreateRevokedToken creates a new revoked token
func (r *RevokedToken) CreateRevokedToken(revokedToken *RevokedTokenModel) error {
	result := r.db.Db.Create(revokedToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindRevokedTokenByTokenId finds a revoked token by token id
func (r *RevokedToken) FindRevokedTokenByTokenId(id string) (*RevokedTokenModel, error) {
	var revokedToken RevokedTokenModel
	result := r.db.Db.Where("token = ?", id).First(&revokedToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &revokedToken, nil
}
