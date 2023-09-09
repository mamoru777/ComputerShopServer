package UserRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u *Models.Usr) error
	Get(ctx context.Context, id uuid.UUID) (*Models.Usr, error)
	GetByLogin(ctx context.Context, login string) (bool, error)
	GetByLoginAndPassword(ctx context.Context, login string, password string) (bool, uuid.UUID, error)
	GetByEmail(ctx context.Context, email string) (bool, error)
	//List (ctx context.Context, )
	Update(ctx context.Context, u *Models.Usr) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateCode(ctx context.Context, ec *Models.EmailCode) error
	GetCode(ctx context.Context, email string) (bool, *Models.EmailCode, error)
	UpdateCode(ctx context.Context, ec *Models.EmailCode) error
}
