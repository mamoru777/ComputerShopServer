package UserRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
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

func (r *UserStorage) Get(ctx context.Context, id uuid.UUID) (*Models.Usr, error) {
	u := new(Models.Usr)
	err := r.db.WithContext(ctx).First(u, id).Error
	return u, err
}

func (r *UserStorage) Update(ctx context.Context, u *Models.Usr) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserStorage) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Models.Usr{ID: id}).Error
}

func (r *UserStorage) GetByLogin(ctx context.Context, login string) (bool, error) {
	u := new(Models.Usr)
	err := r.db.WithContext(ctx).Where("login = ?", login).First(&u).Error
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

func (r *UserStorage) GetByEmail(ctx context.Context, email string) (bool, error) {
	u := new(Models.Usr)
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Запись почты не была найдена")
			return false, nil
		} else {
			log.Println("Ошибка при выполнения запроса на получение почты", err)
			return true, err
		}
	}
	log.Println("Запись почты была найдена")
	return true, nil
}
