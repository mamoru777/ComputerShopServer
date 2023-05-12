package UserRepository

import "context"

type User struct {
	Id       int64  `postgres:"id" gorm:"id;primaryKey"`
	Login    string `postgres:"login" gorm:"login"`
	Password string `postgres:"password" gorm:"password"`
	Name     string `postgres:"id" gorm:"name"`
	LastName string `postgres:"lastname" gorm:"lastname"`
	SurName  string `postgres:"surname" gorm:"surname"`
	Email    string `postgres:"email" gorm:"email"`
	Avatar   byte   `postgres:"avatar" gorm:"avatar"`
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	Get(ctx context.Context, id int64) (*User, error)
	//List (ctx context.Context, )
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id int64) error
}
