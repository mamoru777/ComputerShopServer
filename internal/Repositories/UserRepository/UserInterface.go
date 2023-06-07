package UserRepository

import (
	"context"
	"errors"
)

type Usr struct {
	ID       int64  `db:"id" gorm:"id;primaryKey;type:serial"`
	Login    string `db:"login" gorm:"login"`
	Password string `db:"password" gorm:"password"`
	Name     string `db:"id" gorm:"name"`
	LastName string `db:"lastname" gorm:"lastname"`
	SurName  string `db:"surname" gorm:"surname"`
	Email    string `db:"email" gorm:"email"`
	Avatar   []byte `db:"avatar" gorm:"avatar"`
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

type UserRepository interface {
	Create(ctx context.Context, u *Usr) error
	Get(ctx context.Context, id int64) (*Usr, error)
	//List (ctx context.Context, )
	Update(ctx context.Context, u *Usr) error
	Delete(ctx context.Context, id int64) error
}
