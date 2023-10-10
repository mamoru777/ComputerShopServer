package GoodRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type GoodStorage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GoodStorage {
	return &GoodStorage{
		db: db,
	}
}

func (gs *GoodStorage) Create(ctx context.Context, g *Models.Good) error {
	return gs.db.WithContext(ctx).Create(g).Error
}

func (gs *GoodStorage) Get(ctx context.Context, id uuid.UUID) (*Models.Good, error) {
	g := new(Models.Good)
	err := gs.db.WithContext(ctx).First(g, id).Error
	return g, err
}

func (gs *GoodStorage) Update(ctx context.Context, g *Models.Good) error {
	return gs.db.WithContext(ctx).Save(g).Error
}

func (gs *GoodStorage) Delete(ctx context.Context, id uuid.UUID) error {
	return gs.db.WithContext(ctx).Delete(&Models.Good{Id: id}).Error
}

func (gs *GoodStorage) GetByName(ctx context.Context, name string) (bool, error) {
	g := new(Models.Good)
	err := gs.db.WithContext(ctx).Where("name = ?", name).First(&g).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Запись логина не была найдена")
			return false, nil
		} else {
			log.Println("Ошибка при выполнения запроса на получения логина", err)
			return true, err
		}
	}
	log.Println("Запись логина была найдена")
	return true, nil
}

func (gs *GoodStorage) GetByType(ctx context.Context, gtype string) ([]*Models.Good, error) {
	goods := []*Models.Good{}
	err := gs.db.WithContext(ctx).Where("good_type = ?", gtype).Find(&goods).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Товары по типу", gtype, "не были найдены")
			return nil, err
		} else {
			log.Println("Ошибка при выполнения запроса на получение записей товаров", err)
			return nil, err
		}
	}
	return goods, nil
}
