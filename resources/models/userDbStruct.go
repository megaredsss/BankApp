package models

// User struct for database
type UserDb struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Email      string `gorm:"type:varchar(256);unique"`
	FirstName  string `gorm:"type:varchar(256)"`
	SecondName string `gorm:"type:varchar(256)"`
	Balance    uint   `gorm:"type:decimal(10,2)"`
	Password   string `gorm:"type:varchar(256)"`
}
