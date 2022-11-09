package redismanager

import (
	"github.com/go-redis/redis/v8"
)

type Sku struct {
	SkuID       string  `redis:"sku_id"`
	Name        string  `redis:"name"`
	Price       float64 `redis:"price"`
	Type        string  `redis:"type"`
	Description string  `redis:"description"`
	ImageUrl    string  `redis:"image_url"`
	Quantity    int     `redis:"quantity"`
}

type ItemToUpDate struct {
	Item  string
	Key   string
	Value interface{}
}

type ItemToDeductQty struct {
	Item     string
	Quantity int
}

type SingleRedisItem struct {
	Name string
	// Result *redis.MapStringStringCmd // was for go-redis V9
	Result *redis.StringStringMapCmd
	Store  *Sku
}

func NewSku() *Sku {
	return &Sku{}
}

func (s *Sku) DeductItemQuantity(i int) {
	s.Quantity = s.Quantity - i
}
