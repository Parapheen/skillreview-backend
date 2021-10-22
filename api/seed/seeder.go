package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/Parapheen/skillreview-backend/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
	},
	models.User{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
	},
}

var posts = []models.ReviewRequest{
	models.ReviewRequest{
		Content: "Hello world 1",
	},
	models.ReviewRequest{
		Content: "Hello world 2",
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

	for i, _ := range users {
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