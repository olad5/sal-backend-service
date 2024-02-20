package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	SKUID       uuid.UUID
	Name        string
	Description string
	Price       int64
	MerchantId  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
