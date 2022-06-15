package redismanager

import (
	"context"
	"fmt"
	"log"
	"time"

	rpc "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
	"github.com/go-redis/redis/v9"
)

type AdminRedisManager struct {
	Conn *redis.Client
}

func NewAdminRedisManager(conn *redis.Client) *AdminRedisManager {
	return &AdminRedisManager{
		Conn: conn,
	}
}

func (r *AdminRedisManager) Ping(ctx context.Context) (err error) {
	if err = r.Conn.Ping(ctx).Err(); err != nil {
		log.Printf("ERROR : %+v", err)
		return err
	}
	return nil
}

// ✅ GetAllInventories returns all inventory items from the cache
func (r *AdminRedisManager) GetAllInventories(ctx context.Context, items []*rpc.Sku) (resp *rpc.GetAllInventoriesResp, err error) {
	start := time.Now()
	// defer r.Conn.Close()
	resp = &rpc.GetAllInventoriesResp{}
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

	for _, item := range listOfInventoryItems {
		err = item.Result.Scan(item.Store)
		if err != nil {
			log.Printf("Something went wrong trying to scan redis result: %+v", err)
		}
		fmt.Printf("WAHT THE WAHT : %+v\n", resp)
		sku := &rpc.Sku{
			Name:        item.Store.Name,
			Price:       item.Store.Price,
			Quantity:    int32(item.Store.Quantity),
			SkuId:       item.Store.SkuID,
			Description: item.Store.Description,
			Type:        item.Store.Type,
			ImageUrl:    item.Store.ImageUrl,
		}

		resp.Inventories = append(resp.Inventories, sku)
	}
	fmt.Println("DONE GetAllInventories : ", time.Since(start))
	return
}

func (r *AdminRedisManager) GetOneInventory(ctx context.Context, item string) (resp *rpc.Sku, err error) {
	start := time.Now()
	temp := &Sku{}
	err = r.Conn.HGetAll(ctx, item).Scan(temp)
	resp = &rpc.Sku{
		Name:        temp.Name,
		SkuId:       temp.SkuID,
		Description: temp.Description,
		Type:        temp.Type,
		Price:       temp.Price,
		ImageUrl:    temp.ImageUrl,
		Quantity:    int32(temp.Quantity),
	}
	fmt.Println("END GetOneItem: ", time.Since(start))
	return
}

// ✅ DeductQuantity removes one item quantity from the cache
func (r *AdminRedisManager) DeductQuantity(ctx context.Context, itemName string, quantity int) (err error) {
	start := time.Now()
	// defer r.Conn.Close()
	item := Sku{}

	/*
		❌ READ UP ON HSCAN VS HGET
	*/
	// i, _, err := r.Conn.HScan(ctx, "inventories", 0, itemName, 1).Result()
	// if err != nil {
	// 	return err
	// }

	// unmarshalErr := json.Unmarshal([]byte(i[1]), item.Sku)
	// if unmarshalErr != nil {
	// 	return unmarshalErr
	// }
	// item.DeductItemQuantity(quantity)

	// // // Insert the item back to cache
	// r.Conn.HSet(ctx, "inventories", itemName, item)

	err = r.Conn.HGetAll(ctx, itemName).Scan(&item)
	if err != nil {
		log.Printf("ERROR GET ALL : %+v", err)
	}

	fmt.Println("Got :", item)

	item.DeductItemQuantity(quantity)

	r.Conn.HSet(ctx, itemName, "quantity", item.Quantity)

	fmt.Println("DONE DeductQuantity : ", time.Since(start))
	return nil
}

// ❌
// DeductQuantities removes multiple item quantities from the cache
// func (r *RedisManager) DeductQuantities(ctx context.Context, itemName map[string]interface{}, quantity int) (err error) {
// 	for _ v := range itemName {
// 		item := &SkuInstance{}
// 		err = r.Conn.HGetAll(ctx, v).Scan(item.Sku)
// 		item.DeductItemQuantity(quantity)
// 	}
// }

// AddInventories acts like bulk insert. It add multiple items to the cache
func (r *AdminRedisManager) AddInventories(ctx context.Context, inventories []*rpc.Sku) (err error) {
	start := time.Now()
	if _, err = r.Conn.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		for _, item := range inventories {
			// 		/*
			// 			💡 Set the redis key as the item name. It's fields are tje Sku Struct
			// 		*/
			// 		r.Conn.HSet(ctx, item.Name, item)

			// Previous implementation
			r.Conn.HSet(ctx, item.Name, "name", item.Name, "description", item.Description, "price", item.Price, "quantity", item.Quantity, "sku_id", item.SkuId, "image_url", item.ImageUrl, "type", item.Type)
		}
		return nil
	}); err != nil {
		log.Printf("Error trying to set pipeline for adding inventories : %+v", err)
		return err
	}

	fmt.Println("DONE AddInventories : ", time.Since(start))
	return
}

// AddInventory adds one item to the inventory
func (r *AdminRedisManager) AddInventory(ctx context.Context, item *rpc.Sku) (er error) {
	start := time.Now()
	if _, err := r.Conn.Pipelined(ctx, func(rdb redis.Pipeliner) error {

		r.Conn.Conn(ctx).HSet(ctx, item.Name, "name", item.Name, "description", item.Description, "price", item.Price, "quantity", item.Quantity, "sku_id", item.SkuId, "image_url", item.ImageUrl, "type", item.Type)

		return nil
	}); err != nil {
		log.Printf("Error trying to set pipeline for adding inventories : %+v", err)
		return err
	}
	fmt.Println("DONE AddInventory : ", time.Since(start))
	return nil
}

// func (r *AdminRedisManager) DeleteOne(ctx context.Context, item string) (err error)     {}
// func (r *AdminRedisManager) DeleteMany(ctx context.Context, items []string) (err error) {}
// func (r *AdminRedisManager) UpdateOne(ctx context.Context, item string) (err error)     {}
// func (r *AdminRedisManager) UpdateMany(ctx context.Context, item []string) (err error)  {}
