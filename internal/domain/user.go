package domain

import (
	"api/config"
	"api/pkg/auth"
	"api/pkg/middlewares"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	Firstname      string `json:"first_name"`
	Lastname       string `json:"last_name"`
	Username       string `json:"username" gorm:"size:255;not null;unique" validate:"required"`
	Email          string `json:"email" gorm:"size:255;not null;unique" validate:"required"`
	Password       string `json:"password" gorm:"size:255;not null;" validate:"required, max=20"`
	IsVerified     bool   `json:"is_verified"`
	Timezone       string `json:"timezone"`
	CustomerId     string `json:"customer_id"`
	SubscriptionId string `json:"subscription_id"`
	PriceId        string `json:"price_id"`

	Token   string `json:"token"`
	Secret  string `json:"secret"`
	TwoFa   bool   `json:"two_fa"`
	Country string `json:"country"`
	Active  *bool  `json:"active" gorm:"default:false"`
}

type Users []User

func (u *User) BeforeCreate(db *gorm.DB) error {
	fmt.Println("before create")
	hashed, err := auth.Hash(u.Password)
	if err != nil {
		return err
	}
	db.Model(&u).UpdateColumn("Password", hashed)
	db.Model(&u).UpdateColumn("IsVerified", false)
	db.Model(&u).UpdateColumn("Username", strings.Join(strings.Fields(u.Username), "_"))

	return nil

}

func (u *User) PassHash() error {

	r, _ := regexp.Compile("\\A\\$2a?\\$\\d\\d\\$[./0-9A-Za-z]{53}")
	if !r.MatchString(u.Password) && u.Password != "" {

		hashed, err := auth.Hash(u.Password)
		if err != nil {
			return err
		}
		fmt.Println("password hashed",u.ID,hashed, u.Password)
		u.Password = hashed
		return nil
	}

	return nil

}

func (u *User) LoginCheck(p string) (string, error) {

	err := auth.ComparePassword(u.Password, p)

	fmt.Print("pass compare", err, u.Password, p)
	if err != nil {
		return "", err
	}
	tokenMaker, _ := middlewares.NewPasetoMaker(config.C.Token)
	token, err := tokenMaker.CreateToken(u.ID, u.Email, 24*3*time.Hour)
	return token, err
}

func (User) TableName() string {
	return "users"
}
