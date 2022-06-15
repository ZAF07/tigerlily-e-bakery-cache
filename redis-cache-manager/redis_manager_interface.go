package redismanager

import (
	"context"

	rpc "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
)

type Redismanager interface {
	Ping(ctx context.Context) (err error)
	GetAllInventories(ctx context.Context, items []*rpc.Sku) (resp *rpc.GetAllInventoriesResp, err error)
	GetOneInventory(ctx context.Context, item string) (resp *rpc.Sku, err error)
	DeductQuantity(ctx context.Context, itemName string, quantity int) (err error)
	DeductQuantities(ctx context.Context, itemName []map[string]interface{}) error
	AddInventories(ctx context.Context, inventories []*rpc.Sku) (err error)
	AddInventory(ctx context.Context, item *rpc.Sku) (er error)
	DeleteOne(ctx context.Context, item string) (err error)
	UpdateOne(ctx context.Context, item string, field string, val interface{}) (err error)
	UpdateMany(ctx context.Context, items []map[string]interface{}) (err error)
}
