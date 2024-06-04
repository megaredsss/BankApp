package models

type User struct {
	ID         uint   `gorm:"primaryKey"`
	FirstName  string `gorm:"<-"`
	SecondName string `gorm:"<-"`
	ThirdName  string `gorm:"<-"`
	Balance    uint   `gorm:"<-"`
}
