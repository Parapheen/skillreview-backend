package models

import (
	"errors"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ReviewerApplication struct {
	Base
	Description string           `gorm:"size:255;not null;" json:"description"`
	State       ApplicationState `gorm:"not null; default:'recieved'" json:"state"`
	Rating      int              `gorm:"not null" json:"rating"`
	AuthorID    uint32           `gorm:"not null" json:"-"`
	Author      User             `json:"author"`
}

type ApplicationState string

const (
	Recieved ApplicationState = "recieved"
	Progress ApplicationState = "in_progress"
	Approved ApplicationState = "approved"
	Declined ApplicationState = "declined"
)

func (application *ReviewerApplication) Prepare() {
	application.UUID = uuid.NewV4()
	application.Description = strings.TrimSpace(application.Description)
	application.Author = User{}
	application.CreatedAt = time.Now()
	application.UpdatedAt = time.Now()
}

func (application *ReviewerApplication) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if application.Description == "" {
			return errors.New("Required Content")
		}
		if application.State == "" {
			return errors.New("Required State")
		}
		return nil
	default:
		if application.Description == "" {
			return errors.New("Required Content")
		}
		if application.Rating < 1 {
			return errors.New("Required Rating")
		}
		return nil
	}
}

func (application *ReviewerApplication) SaveApplication(db *gorm.DB) (*ReviewerApplication, error) {
	var err error
	err = db.Model(&ReviewerApplication{}).Create(&application).Error
	if err != nil {
		return &ReviewerApplication{}, err
	}
	if application.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", application.AuthorID).Take(&application.Author).Error
		if err != nil {
			return &ReviewerApplication{}, err
		}
	}
	return application, nil
}

func (application *ReviewerApplication) FindAllApplications(db *gorm.DB) (*[]ReviewerApplication, error) {
	var err error
	applications := []ReviewerApplication{}
	err = db.Model(&ReviewerApplication{}).Limit(100).Find(&applications).Error
	if err != nil {
		return &[]ReviewerApplication{}, err
	}
	if len(applications) > 0 {
		for i := range applications {
			err := db.Model(&User{}).Where("id = ?", applications[i].AuthorID).Take(&applications[i].Author).Error
			if err != nil {
				return &[]ReviewerApplication{}, err
			}
		}
	}
	return &applications, nil
}

func (application *ReviewerApplication) FindApplicationByUUID(db *gorm.DB, uid uuid.UUID) (*ReviewerApplication, error) {
	var err error
	err = db.Model(&ReviewerApplication{}).Where("uuid = ?", uid).Take(&application).Error
	if err != nil {
		return &ReviewerApplication{}, err
	}
	if application.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", application.AuthorID).Take(&application.Author).Error
		if err != nil {
			return &ReviewerApplication{}, err
		}
	}
	return application, nil
}

func (application *ReviewerApplication) UpdateApplication(db *gorm.DB, uid uuid.UUID) (*ReviewerApplication, error) {

	err := db.Model(&ReviewerApplication{}).Where("uuid = ?", uid).Updates(
		ReviewerApplication{
			State: application.State,
		}).Error
	if err != nil {
		return &ReviewerApplication{}, db.Error
	}
	if application.ID != 0 {
		err = db.Model(&User{}).Where("id = ?", application.AuthorID).Take(&application.Author).Error
		if err != nil {
			return &ReviewerApplication{}, err
		}
	}
	return application, nil
}
