package seed

import (
	"log"

	"github.com/Parapheen/skillreview-backend/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
	},
	{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
	},
}

var posts = []models.ReviewRequest{
	{
		Description: "Hello world 1",
	},
	{
		Description: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	err := db.DropTableIfExists(&models.ReviewRequest{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.ReviewRequest{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Model(&models.ReviewRequest{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Model(&models.ReviewRequest{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
