package redismanager

import (
	"log"

	"context"

	rpc "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
	"github.com/go-redis/redis/v9"
)

type RedisManager struct {
	Conn *redis.Client
}

func NewRedisManager(conn *redis.Client) *RedisManager {
	return &RedisManager{
		Conn: conn,
	}
}

func (r *RedisManager) GetAllInventories(ctx context.Context, items []*rpc.Sku) (resp *rpc.GetAllInventoriesResp, err error) {
	defer r.Conn.Close()
	pipe := r.Conn.Pipeline()
	listOfInventoryItems := []SingleRedisItem{}

	for _, v := range items {
		singleItem := SingleRedisItem{
			Name:   v.Name,
			Store:  NewSku(),
			Result: pipe.HGetAll(ctx, v.Name),
		}
		listOfInventoryItems = append(listOfInventoryItems, singleItem)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error executing Redis pipeline : %+v", err)
	}

	for _, v := range listOfInventoryItems {
		err = v.Result.Scan(v.Store)
		if err != nil {
			log.Printf("Something went wrong trying to scan redis result: %+v", err)
		}
		resp.Inventories = append(resp.Inventories, v.Store)
	}

	return
}

func (r *RedisManager) DeductQuantity(ctx context.Context, itemName string, quantity int) (err error) {

	defer r.Conn.Close()
	item := &SkuInstance{}

	err = r.Conn.HGetAll(ctx, itemName).Scan(item.Sku)

	item.DeductItemQuantity(quantity)

	// This might not work. *rpc
	r.Conn.HSet(ctx, itemName, "quantity", item.Sku.GetQuantity())

	return
}
