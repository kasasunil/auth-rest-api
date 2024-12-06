package user

import (
	"errors"
	"github.com/kasasunil/auth-rest-api/internal/database"
	"gorm.io/gorm"
	"log"
)

type User struct {
	db *database.Db
}

type UserModel struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	JwtToken string `gorm:"not null"`
}

func New(db *database.Db) *User {
	return &User{
		db: db,
	}
}

func (u *User) CreateUser(user *UserModel) error {
	// TODO: Try to hash password if possible
	result := u.db.Db.Create(user)
	if result.Error != nil {
		log.Println("User creation failed: ", result.Error)
	} else {
		log.Println("User created successfully : ")
	}
	return result.Error
}

func (u *User) FindUserByEmailId(email string) (*UserModel, error) {
	var user UserModel
	result := u.db.Db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("User not found: ", email)
		} else {
			log.Println("Error finding user: ", result.Error.Error())
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *User) UpdateUserByEmailId(email string, updatedData *UserModel) error {
	result := u.db.Db.Model(&UserModel{}).Where("email = ?", email).Updates(updatedData)
	if result.Error != nil {
		log.Println("Error updating user: ", result.Error.Error())
		return result.Error
	}
	if result.RowsAffected == 0 {
		log.Println("No user found with email: ", email)
		return gorm.ErrRecordNotFound
	}
	log.Println("User updated successfully: ", email)
	return nil
}
