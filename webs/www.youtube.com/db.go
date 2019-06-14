package youtube

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Record struct {
	gorm.Model
	Link  string
	Title string
}

func (d *Record) GetDatabaseName() string {
	return "default_db"
}
func (d *Record) TableName() string {
	return "crawler"
}

func Write(link string, title string) {
	//need modify
	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	r := Record{Link: link, Title: title}
	db.Create(&r)
}
