package redismanager

import (
	rpc "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
	"github.com/go-redis/redis/v9"
)

type SkuInstance struct {
	Sku *rpc.Sku
}

type SingleRedisItem struct {
	Name   string
	Result *redis.MapStringStringCmd
	Store  *rpc.Sku
}

func NewSku() *rpc.Sku {
	return &rpc.Sku{}
}

func (s *SkuInstance) DeductItemQuantity(i int) {
	s.Sku.Quantity = s.Sku.Quantity - int32(i)
}
