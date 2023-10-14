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
	err := r.db.Preload("Orders").WithContext(ctx).First(u, id).Error
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
	err := r.db.Preload("Orders").WithContext(ctx).Where("login = ?", login).First(u).Error
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
	err := r.db.Preload("Orders").WithContext(ctx).Where("email = ?", email).First(u).Error
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

func (r *UserStorage) GetByEmailUser(ctx context.Context, email string) (*Models.Usr, error) {
	u := new(Models.Usr)
	err := r.db.Preload("Orders").WithContext(ctx).Where("email = ?", email).First(u).Error
	if err != nil {
		log.Println("Не удалось найти запись пользователя по почте или неизвестная ошибка", err)
		return nil, err
	}
	return u, nil
}

func (r *UserStorage) GetByLoginAndPassword(ctx context.Context, login string, password string) (bool, uuid.UUID, string, error) {
	var emptyUUID uuid.UUID
	u := new(Models.Usr)
	err := r.db.Preload("Orders").WithContext(ctx).Where("login = ? AND password = ?", login, password).First(u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Запись пользователя не была найдена")
			return false, emptyUUID, "", nil
		} else {
			log.Println("Ошибка при выполнения запроса на получение пользователя", err)
			return false, emptyUUID, "", err
		}
	}
	log.Println("Запись пользователя была найдена")
	return true, u.ID, u.Role, nil
}

func (r *UserStorage) CreateCode(ctx context.Context, ec *Models.EmailCode) error {
	return r.db.WithContext(ctx).Create(ec).Error
}

func (r *UserStorage) GetCode(ctx context.Context, email string) (bool, *Models.EmailCode, error) {
	ec := new(Models.EmailCode)
	err := r.db.WithContext(ctx).Where("email = ?", email).First(ec).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Запись с почтой и кодом не была найдена")
			return false, nil, nil
		} else {
			log.Println("Ошибка при выполнении запроса на получение почты и кода", err)
			return false, nil, err
		}
	}
	log.Println("Запись почты и кода была найдена")
	return true, ec, err
}

func (r *UserStorage) UpdateCode(ctx context.Context, ec *Models.EmailCode) error {
	return r.db.WithContext(ctx).Save(ec).Error
}
