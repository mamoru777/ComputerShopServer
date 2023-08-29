package UserRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
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

func (r *UserStorage) Create(ctx context.Context, u *Models.Usr) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserStorage) Get(ctx context.Context, id int64) (*Models.Usr, error) {
	u := new(Models.Usr)
	err := r.db.WithContext(ctx).First(u, id).Error
	return u, err
}

func (r *UserStorage) Update(ctx context.Context, u *Models.Usr) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserStorage) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&Models.Usr{ID: id}).Error
}
