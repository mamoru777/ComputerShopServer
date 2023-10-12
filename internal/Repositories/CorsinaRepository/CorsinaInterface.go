package CorsinaRepository

import (
	"ComputerShopServer/internal/Repositories/Models"
	"context"
	"github.com/google/uuid"
)

type CorsinaRepository interface {
	Create(ctx context.Context, g *Models.Corsina) error
	Get(ctx context.Context, id uuid.UUID) (*Models.Corsina, error)
	Update(ctx context.Context, g *Models.Corsina) error
	GetByUser(ctx context.Context, userId uuid.UUID) (*Models.Corsina, error)
}
