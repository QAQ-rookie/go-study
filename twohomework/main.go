package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"os"
)

var DB *gorm.DB

type Plans struct {
	Id          int    `gorm:"column:id" json:"id,omitempty"`
}

func initDB() error {
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return errors.Wrap(err, "db conn error")
	}
	DB = db
	return nil
}

func (p *Plans) TableName() string {
	return "plan"
}

func (p *Plans) GetAllPlans() (Plans, error) {
	plans := Plans{}
	result := DB.Table(p.TableName()).Where("id = ?", 6).First(&plans)
	return plans, errors.Wrap(result.Error, "get all plan")
}

func main() {
	if err := initDB(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var p Plans
	_, err := p.GetAllPlans()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println(err)
	}
}
