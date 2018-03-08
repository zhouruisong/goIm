package model

import (
	"time"

	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	uuid2 "github.com/satori/go.uuid"
)

var DB *gorm.DB

func InitDb() (*gorm.DB, error) {

	db, err := gorm.Open("mysql", "root:@tcp(111.230.235.49:3306)/goim?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	DB = db
	db.AutoMigrate(&User{})
	return db, err
}

type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	BaseModel
	Uuid     string
	Nickname string
	Password string
}

//根据用户名查找用户信息
func (this *User) GetUserByName(name string) *User {
	DB.Where("nickname = ?", name).Find(&this)
	return this
}

//根据uuid  查找用户信息
func (this *User) GetUserByUuid(uuid string) *User {
	DB.Where("uuid = ?", uuid).Find(&this)
	return this
}

//创建用户
func (this *User) CreateUser(name string, password string) bool {
	uuid, _ := uuid2.NewV4()
	this.Nickname = name
	this.Password = password
	this.Uuid = uuid.String()
	err := DB.Where(&User{}).Create(&this).Error

	if err != nil {
		log.Fatalf("创建用户发生异常：%v", err)
		return false
	}
	return true
}
