package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	Base
	Nickname       string          `gorm:"size:255;not null;unique" json:"nickname"`
	Email          string          `gorm:"size:100;" json:"email"`
	Steam64ID      string          `gorm:"size:255;default:null;unique" json:"steam64Id"`
	Steam32ID      string          `gorm:"size:255;default:null;unique" json:"steam32Id"`
	Avatar         string          `gorm:"size:255;" json:"avatar"`
	Rank           string          `gorm:"size:255;" json:"rank"`
	ReviewRequests []ReviewRequest `gorm:"constraint:OnDelete:CASCADE;foreignkey:author_id" json:"review_requests"`
	Reviews        []Review        `gorm:"constraint:OnDelete:CASCADE;foreignkey:author_id" json:"reviews"`
	Plan           PlanType        `gorm:"default:'basic';" json:"plan"`
}

type PlanType string

const (
	Basic PlanType = "basic"
	Pro   PlanType = "pro"
)

func (u *User) Prepare() {
	u.UUID = uuid.NewV4()
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Required Nickname")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	err := db.Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, nil
}

func (u *User) FindUserByUIID(db *gorm.DB, uid uuid.UUID) (*User, error) {
	err := db.Model(User{}).Where("uuid = ?", uid).Preload("ReviewRequests").Preload("ReviewRequests.Reviews").Preload("Reviews").Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindUserBySteamID(db *gorm.DB, steam64Id string) (*User, error) {
	err := db.Model(User{}).Where("steam64_id = ?", steam64Id).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateUser(db *gorm.DB, uid uuid.UUID) (*User, error) {

	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"nickname":   u.Nickname,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err := db.Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteUser(db *gorm.DB, uid uuid.UUID) (int64, error) {

	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
