package Models

import (
	"errors"
	"github.com/google/uuid"
)

type Usr struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"` //`db:"id" gorm:"id;primaryKey;type:serial"`
	Login    string    `json:"login" gorm:"name"`                                         //`db:"login" gorm:"login"`
	Password string    `json:"password" gorm:"password"`                                  //`db:"password" gorm:"password"`
	Name     string    `json:"name" gorm:"name"`                                          //`db:"id" gorm:"name"`
	LastName string    `json:"lastName" gorm:"last_name"`                                 //`db:"lastname" gorm:"lastname"`
	SurName  string    `json:"surName" gorm:"sur_name"`                                   //`db:"surname" gorm:"surname"`
	Email    string    `json:"email" gorm:"email"`                                        //`db:"email" gorm:"email"`
	Avatar   []byte    `json:"avatar" gorm:"avatar"`                                      //`db:"avatar" gorm:"avatar"`
	Role     string    `json:"role" gorm:"role"`
}

type EmailCode struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email string    `json:"email" gorm:"email"`
	Code  string    `json:"code" gorm:"gorm"`
}

type Good struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	GoodType    string    `json:"good_type" gorm:"good_type"`
	Name        string    `json:"name" gorm:"name"`
	Description string    `json:"description" gorm:"description"`
	Price       float64   `json:"price" gorm:"price"`
	Avatar      []byte    `json:"avatar" gorm:"avatar"`
}

type Order struct {
	Id uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
}

func (u *Usr) Validate() error {
	if u.Login == "" {
		return errors.New("Login required value")
	}
	if u.Password == "" {
		return errors.New("Password required value")
	}
	return nil
}
