package OrderRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, o *Models.Order) error
	Get(ctx context.Context, id uuid.UUID) (*Models.Order, error)
	GetAll(ctx context.Context) ([]*Models.Order, error)
	Update(ctx context.Context, o *Models.Order) error
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*Models.Order, error)
}
