package UserRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u *Models.Usr) error
	Get(ctx context.Context, id uuid.UUID) (*Models.Usr, error)
	//List (ctx context.Context, )
	Update(ctx context.Context, u *Models.Usr) error
	Delete(ctx context.Context, id uuid.UUID) error
}
