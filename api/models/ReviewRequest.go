package models

import (
	"errors"
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
	Position           MatchPos     `gorm:"size:255;" json:"position"`
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

type MatchPos string

const (
	HardSupport MatchPos = "Hard Support"
	Support     MatchPos = "Support"
	Offlane     MatchPos = "Offlane"
	Mid         MatchPos = "Mid"
	Carry       MatchPos = "Carry"
)

type Filters struct {
	State    RequestState
	Position MatchPos
}

func (rr *ReviewRequest) Prepare() {
	rr.UUID = uuid.NewV4()
	rr.Description = strings.TrimSpace(rr.Description)
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

func (rr *ReviewRequest) FindAllReviewRequests(db *gorm.DB, limit int, offset int, filters Filters) (*[]ReviewRequest, error) {
	var err error

	reviewRequests := []ReviewRequest{}
	if filters.Position != "" {
		err = db.Model(&ReviewRequest{}).Where("state = ? AND position = ?", filters.State, filters.Position).Offset(offset).Limit(limit).Order("created_at desc").Preload("Reviews").Find(&reviewRequests).Error
	} else {
		err = db.Model(&ReviewRequest{}).Where("state = ?", filters.State).Offset(offset).Limit(limit).Order("created_at desc").Preload("Reviews").Find(&reviewRequests).Error
	}
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

func (rr *ReviewRequest) CountReviewRequests(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&ReviewRequest{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (rr *ReviewRequest) FindReviewRequestByUUID(db *gorm.DB, pid uuid.UUID) (*ReviewRequest, error) {
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

func (rr *ReviewRequest) UpdateReviewRequest(db *gorm.DB, rrUUID uuid.UUID) (*ReviewRequest, error) {
	err := db.Model(&ReviewRequest{}).Where("uuid = ?", rrUUID).Updates(
		ReviewRequest{
			Description: rr.Description,
			State:       rr.State,
		}).Error
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

func (rr *ReviewRequest) DeleteReviewRequest(db *gorm.DB, uuid uuid.UUID, author_id uint32) (int64, error) {

	err := db.Model(&ReviewRequest{}).Where("uuid = ? and author_id = ?", uuid, author_id).Take(&ReviewRequest{}).Delete(&ReviewRequest{}).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.New("Request not found")
		}
		return 0, db.Error
	}
	return 1, nil
}
