package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	SKUID       uuid.UUID
	Name        string
	Description string
	Price       float64
	MerchantId  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
