package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Review struct {
	Base
	ReviewRequestUUID string      `gorm:"not null" json:"review_request_uuid"`
	Description       string      `gorm:"size:255;not null;" json:"description"`
	AuthorID          uint32      `gorm:"not null" json:"-"`
	ReviewRequestID   uint32      `gorm:"not null" json:"-"`
	State             ReviewState `gorm:"not null; default:'submitted'" json:"state"`
	RateLaning        int         `gorm:"not null" json:"rate_laning"`
	RateTeamFights    int         `gorm:"not null" json:"rate_teamfights"`
	RateOverall       int         `gorm:"not null" json:"rate_overall"`
	Author            User        `gorm:"constraint:OnDelete:CASCADE;foreignkey:author_id" json:"author"`
}

type ReviewState string

const (
	Submitted RequestState = "submitted"
	Accepted  RequestState = "accepted"
	Reviewed  RequestState = "reviewed"
)

func (review *Review) Prepare() {
	review.UUID = uuid.NewV4()
	review.Description = html.EscapeString(strings.TrimSpace(review.Description))
	review.Author = User{}
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
}

func (review *Review) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if review.Description == "" {
			return errors.New("Required Content")
		}
		if review.State == "" {
			return errors.New("Required State")
		}
		return nil
	default:
		if review.Description == "" {
			return errors.New("Required Content")
		}
		if review.RateLaning < 1 {
			return errors.New("Required Laning Rating")
		}
		if review.RateTeamFights < 1 {
			return errors.New("Required TeamFighting Rating")
		}
		if review.RateOverall < 1 {
			return errors.New("Required Overall Rating")
		}
		return nil
	}
}

func (review *Review) SaveReview(db *gorm.DB) (*Review, error) {
	var err error
	err = db.Model(&Review{}).Create(&review).Error
	if err != nil {
		return &Review{}, err
	}
	if review.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", review.AuthorID).Take(&review.Author).Error
		if err != nil {
			return &Review{}, err
		}
	}
	return review, nil
}

func (review *Review) FindAllReviews(db *gorm.DB) (*[]Review, error) {
	var err error
	reviews := []Review{}
	err = db.Model(&Review{}).Limit(100).Find(&reviews).Error
	if err != nil {
		return &[]Review{}, err
	}
	if len(reviews) > 0 {
		for i := range reviews {
			err := db.Model(&User{}).Where("id = ?", reviews[i].AuthorID).Take(&reviews[i].Author).Error
			if err != nil {
				return &[]Review{}, err
			}
		}
	}
	return &reviews, nil
}

func (review *Review) FindReviewByUIID(db *gorm.DB, pid uuid.UUID) (*Review, error) {
	var err error
	err = db.Model(&Review{}).Where("uuid = ?", pid).Take(&review).Error
	if err != nil {
		return &Review{}, err
	}
	if review.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", review.AuthorID).Take(&review.Author).Error
		if err != nil {
			return &Review{}, err
		}
	}
	return review, nil
}

func (review *Review) UpdateReview(db *gorm.DB, reviewID uuid.UUID) (*Review, error) {

	err := db.Model(&Review{}).Where("uuid = ?", reviewID).Updates(
		Review{
			Description: review.Description,
			State:       review.State,
		}).Error
	if err != nil {
		return &Review{}, db.Error
	}
	if review.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", review.AuthorID).Take(&review.Author).Error
		if err != nil {
			return &Review{}, err
		}
	}
	return review, nil
}

func (review *Review) DeleteReview(db *gorm.DB, pid uuid.UUID, uid uuid.UUID) (int64, error) {

	db = db.Model(&Review{}).Where("id = ? and author_id = ?", pid, uid).Take(&Review{}).Delete(&Review{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Review not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
