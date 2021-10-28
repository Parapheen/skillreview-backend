package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type ReviewRequest struct {
	Base
	MatchId               uint32    `gorm:"not null" json:"match_id"`
	Description           string    `gorm:"size:255;not null;" json:"description"`
	Author                User      `json:"author"`
	AuthorID              uint32    `gorm:"not null" json:"-"`
	SelfRateLaning        int       `gorm:"not null" json:"self_rate_laning"`
	SelfRateTeamFights    int       `gorm:"not null" json:"self_rate_teamfights"`
	SelfRateOverall       int       `gorm:"not null" json:"self_rate_overall"`
}

func (rr *ReviewRequest) Prepare() {
	rr.UUID = uuid.NewV4()
	rr.Description = html.EscapeString(strings.TrimSpace(rr.Description))
	rr.Author = User{}
	rr.CreatedAt = time.Now()
	rr.UpdatedAt = time.Now()
}

func (rr *ReviewRequest) Validate() error {
	if rr.Description == "" {
		return errors.New("Required Content")
	}
	if rr.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (rr *ReviewRequest) SaveReviewRequest(db *gorm.DB) (*ReviewRequest, error) {
	var err error
	err = db.Model(&ReviewRequest{}).Create(&rr).Error
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

func (rr *ReviewRequest) FindAllReviewRequests(db *gorm.DB) (*[]ReviewRequest, error) {
	var err error
	reviewRequests := []ReviewRequest{}
	err = db.Model(&ReviewRequest{}).Limit(100).Find(&rr).Error
	if err != nil {
		return &[]ReviewRequest{}, err
	}
	if len(reviewRequests) > 0 {
		for i, _ := range reviewRequests {
			err := db.Model(&User{}).Where("id = ?", reviewRequests[i].AuthorID).Take(&reviewRequests[i].Author).Error
			if err != nil {
				return &[]ReviewRequest{}, err
			}
		}
	}
	return &reviewRequests, nil
}

func (rr *ReviewRequest) FindReviewRequestByID(db *gorm.DB, pid uuid.UUID) (*ReviewRequest, error) {
	var err error
	err = db.Model(&ReviewRequest{}).Where("id = ?", pid).Take(&rr).Error
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

	var err error

	db = db.Model(&ReviewRequest{}).Where("id = ?", rr.ID).Take(&ReviewRequest{}).UpdateColumns(
		map[string]interface{}{
			"description":  rr.Description,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &ReviewRequest{}, db.Error
	}
	if rr.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", rr.AuthorID).Take(&rr.Author).Error
		if err != nil {
			return &ReviewRequest{}, err
		}
	}
	return rr, nil
}

func (rr *ReviewRequest) DeleteReviewRequest(db *gorm.DB, pid uuid.UUID, uid uuid.UUID) (int64, error) {

	db = db.Model(&ReviewRequest{}).Where("id = ? and author_id = ?", pid, uid).Take(&ReviewRequest{}).Delete(&ReviewRequest{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}