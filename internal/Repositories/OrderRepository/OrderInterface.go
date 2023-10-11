package OrderRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, g *Models.Order) error
	Get(ctx context.Context, id uuid.UUID) (*Models.Order, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]*Models.Order, error)
}
