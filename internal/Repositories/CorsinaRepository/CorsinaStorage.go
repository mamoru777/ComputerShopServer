package CorsinaRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CorsinaStorage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *CorsinaStorage {
	return &CorsinaStorage{
		db: db,
	}
}

func (cs *CorsinaStorage) Create(ctx context.Context, c *Models.Corsina) error {
	return cs.db.WithContext(ctx).Create(c).Error
}

func (cs *CorsinaStorage) Get(ctx context.Context, id uuid.UUID) (*Models.Corsina, error) {
	c := new(Models.Corsina)
	err := cs.db.Preload("Goods").WithContext(ctx).First(c, id).Error
	return c, err
}

func (cs *CorsinaStorage) Update(ctx context.Context, c *Models.Corsina) error {
	return cs.db.WithContext(ctx).Save(c).Error
}
func (cs *CorsinaStorage) GetByUser(ctx context.Context, userId uuid.UUID) (*Models.Corsina, error) {
	c := new(Models.Corsina)
	err := cs.db.Preload("Goods").WithContext(ctx).Where("usr_id = ?", userId).First(&c).Error
	return c, err
}
