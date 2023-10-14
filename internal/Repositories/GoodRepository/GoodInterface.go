package GoodRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
)

type GoodRepository interface {
	Create(ctx context.Context, g *Models.Good) error
	Get(ctx context.Context, id uuid.UUID) (*Models.Good, error)
	GetByName(ctx context.Context, name string) (bool, error)
	GetByType(ctx context.Context, gtype string) ([]*Models.Good, error)
	Update(ctx context.Context, g *Models.Good) error
	Delete(ctx context.Context, corsina *Models.Corsina, good *Models.Good) error
}
