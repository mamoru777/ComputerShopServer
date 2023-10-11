package OrderRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type OrderStorage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *OrderStorage {
	return &OrderStorage{
		db: db,
	}
}

func (os *OrderStorage) Create(ctx context.Context, o *Models.Order) error {
	return os.db.WithContext(ctx).Create(o).Error
}

func (os *OrderStorage) Get(ctx context.Context, id uuid.UUID) (*Models.Order, error) {
	o := new(Models.Order)
	err := os.db.Preload("Goods").WithContext(ctx).First(o, id).Error
	return o, err
}

func (os *OrderStorage) GetByUserId(ctx context.Context, userId uuid.UUID) ([]*Models.Order, error) {
	orders := []*Models.Order{}
	err := os.db.Preload("Goods").WithContext(ctx).Where("usr_id = ?", userId).Find(&orders).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Заказы пользователя ", userId, " не были найдены")
			return nil, err
		} else {
			log.Println("Ошибка при выполнения запроса на получение заказов пользователя ", userId, err)
			return nil, err
		}
	}
	return orders, nil
}
