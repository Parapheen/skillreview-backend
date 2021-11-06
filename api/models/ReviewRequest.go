package models

import (
	"errors"
	"html"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ReviewRequest struct {
	Base
	MatchId            string       `gorm:"not null" json:"match_id"`
	Description        string       `gorm:"size:255;not null;" json:"description"`
	State              RequestState `gorm:"not null; default:'open'" json:"state"`
	HeroPlayed         int          `gorm:"size:255;not null;" json:"hero_played"`
	AuthorRank         string       `gorm:"not null;" json:"author_rank"`
	SelfRateLaning     int          `gorm:"not null" json:"self_rate_laning"`
	SelfRateTeamfights int          `gorm:"not null" json:"self_rate_teamfights"`
	SelfRateOverall    int          `gorm:"not null" json:"self_rate_overall"`
	AuthorID           uint32       `gorm:"not null" json:"-"`
	Author             User         `json:"author"`
	Reviews            []Review     `gorm:"constraint:OnDelete:CASCADE;foreignkey:review_request_id" json:"reviews"`
}

type RequestState string

const (
	Open   RequestState = "open"
	Closed RequestState = "closed"
)

func (rr *ReviewRequest) Prepare() {
	rr.UUID = uuid.NewV4()
	rr.Description = html.EscapeString(strings.TrimSpace(rr.Description))
	rr.Author = User{}
	rr.CreatedAt = time.Now()
	rr.UpdatedAt = time.Now()
}

func (rr *ReviewRequest) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if rr.Description == "" {
			return errors.New("Required Content")
		}
		if rr.Author.UUID.String() == "" {
			return errors.New("Required Author")
		}
		return nil
	default:
		if rr.Description == "" {
			return errors.New("Required Content")
		}
		if rr.Author.UUID.String() == "" {
			return errors.New("Required Author")
		}
		if rr.SelfRateLaning < 1 {
			return errors.New("Required Self Laning Rating")
		}
		if rr.SelfRateTeamfights < 1 {
			return errors.New("Required Self TeamFighting Rating")
		}
		if rr.SelfRateOverall < 1 {
			return errors.New("Required Self Overall Rating")
		}
		return nil
	}
}

func (rr *ReviewRequest) SaveReviewRequest(db *gorm.DB) (*ReviewRequest, error) {
	var err error
	err = db.Model(&ReviewRequest{}).Create(&rr).Error
	if err != nil {
		return &ReviewRequest{}, err
	}
	if rr.ID != 0 {
		err = db.Model(&User{}).Where("uuid = ?", rr.Author.UUID).Take(&rr.Author).Error
		if err != nil {
			return &ReviewRequest{}, err
		}
	}
	return rr, nil
}

func (rr *ReviewRequest) FindAllReviewRequests(db *gorm.DB) (*[]ReviewRequest, error) {
	reviewRequests := []ReviewRequest{}
	err := db.Model(&ReviewRequest{}).Limit(100).Preload("Reviews").Find(&reviewRequests).Error
	if err != nil {
		return &[]ReviewRequest{}, err
	}
	if len(reviewRequests) > 0 {
		for i := range reviewRequests {
			err := db.Model(&User{}).Where("id = ?", reviewRequests[i].AuthorID).Take(&reviewRequests[i].Author).Error
			if err != nil {
				return &[]ReviewRequest{}, err
			}
		}
	}
	return &reviewRequests, nil
}

func (rr *ReviewRequest) FindReviewRequestByUIID(db *gorm.DB, pid uuid.UUID) (*ReviewRequest, error) {
	var err error
	err = db.Model(&ReviewRequest{}).Where("uuid = ?", pid).Preload("Reviews").Preload("Reviews.Author").Take(&rr).Error
	if err != nil {
		return &ReviewRequest{}, err
	}
	if rr.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", rr.AuthorID).Take(&rr.Author).Error
		if err != nil {
			return &ReviewRequest{}, err
		}
	}
	return rr, nil
}

func (rr *ReviewRequest) UpdateReviewRequest(db *gorm.DB) (*ReviewRequest, error) {
	err := db.Model(&ReviewRequest{}).Where("id = ?", rr.ID).Updates(
		ReviewRequest{
			Description: rr.Description,
			State:       rr.State,
		})
	if err != nil {
		return &ReviewRequest{}, db.Error
	}
	if rr.ID != 0 {
		err := db.Model(&User{}).Where("id = ?", rr.AuthorID).Take(&rr.Author).Error
		if err != nil {
			return &ReviewRequest{}, err
		}
	}
	return rr, nil
}

func (rr *ReviewRequest) DeleteReviewRequest(db *gorm.DB, pid uuid.UUID, uid uuid.UUID) (int64, error) {

	db = db.Model(&ReviewRequest{}).Where("id = ? and author_id = ?", pid, uid).Take(&ReviewRequest{}).Delete(&ReviewRequest{})

	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
