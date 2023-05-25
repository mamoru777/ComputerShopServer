package UserRepository

import (
	"context"
	"gorm.io/gorm"
)

type UserStorage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (r *UserStorage) Create(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserStorage) Get(ctx context.Context, id int64) (*User, error) {
	u := new(User)
	err := r.db.WithContext(ctx).First(u, id).Error
	return u, err
}

func (r *UserStorage) Update(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserStorage) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&User{Id: id}).Error
}
