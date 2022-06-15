package main

import (
	"context"
	"fmt"
	"log"
	"time"

	r_manager "github.com/ZAF07/tigerlily-e-bakery-cache/redis-cache-manager"
	// "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
	"github.com/ZAF07/tigerlily-e-bakery-inventories/api/rpc"
	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()

// //  BULK INSERT
var lemon = &rpc.Sku{
	Name:        "lemon tart",
	Description: "Sweet and sour",
	Price:       2.5,
	Quantity:    10,
	SkuId:       "11111",
	ImageUrl:    "lemon_tart.com",
	Type:        "tart",
}
var egg = &rpc.Sku{
	Name:        "egg tart",
	Description: "Eggy",
	Price:       2.5,
	Quantity:    10,
	SkuId:       "11111",
	ImageUrl:    "egg_tart.com",
	Type:        "tart",
}
var cheese = &rpc.Sku{
	Name:        "cheese tart",
	Description: "Cheesy",
	Price:       2.5,
	Quantity:    10,
	SkuId:       "11111",
	ImageUrl:    "cheese_tart.com",
	Type:        "tart",
}

/*
	❌ USE GO TEST MODULE
*/
func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	start := time.Now()
	// HEALTH CHECK
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("ERROR REDIS : %+v", err)
	}
	cache := r_manager.NewRedisManager(rdb)

	itemsToAdd := []*rpc.Sku{}
	itemsToAdd = append(itemsToAdd, lemon)
	itemsToAdd = append(itemsToAdd, egg)
	itemsToAdd = append(itemsToAdd, cheese)

	bulkErr := cache.AddInventories(ctx, itemsToAdd)
	if bulkErr != nil {
		log.Printf("BULD ADD ERROR : %+v", bulkErr)
	}
	fmt.Println("DONE BULK ADD")
	fmt.Println("************************************")
	item := &rpc.Sku{
		Name:        "lemon tart",
		Price:       2.4,
		SkuId:       "001001",
		Type:        "Tart",
		ImageUrl:    "lemon_tart.com",
		Quantity:    10,
		Description: "Sweet & Sour",
	}

	// INSERT ONE ITEM TO CACHE
	err := cache.AddInventory(ctx, item)
	if err != nil {
		log.Fatalf("Failed!! : %+v", err)
	}
	fmt.Println("DONE ADDING TO INVENTORY")
	fmt.Println("************************************")

	// DEDUCT ONE ITEM QUANTITY
	if err := cache.DeductQuantity(ctx, item.Name, 8); err != nil {
		log.Fatalf("ERROR DEDUCT : %+v", err)
	}
	fmt.Println("************************************")

	// // GET ALL INVENTORIES
	// items := []*rpc.Sku{}
	// items = append(items, item)
	// resp, err := cache.GetAllInventories(ctx, items)
	// if err != nil {
	// 	log.Printf("ERROR ALL : %+v\n", err)
	// }
	// fmt.Printf("ALL INVENTORIES : %+v\n", resp)

	// Get one item
	GetOneItem(rdb, cheese, "price")
	updateOne(rdb, "lemon tart", "price", 200)
	updateMany(rdb, []map[string]interface{}{
		{"item": "lemon tart", "key": "price", "value": 500},
		{"item": "egg tart", "key": "price", "value": 500},
		{"item": "cheese tart", "key": "price", "value": 500},
	})
	deductMany(rdb, []map[string]interface{}{})
	// delOne(rdb, "lemon tart")
	// delMany(rdb, []string{"lemon tart", "cheese tart", "egg tart"})

	fmt.Println("DONE MAIN --> ", time.Since(start))
}

func GetOneItem(r *redis.Client, item *rpc.Sku, field string) (resp *rpc.Sku, err error) {
	start := time.Now()

	rc := r_manager.NewAdminRedisManager(r)
	temp := r_manager.NewSku()

	resp, err = rc.GetOneInventory(ctx, lemon.Name)
	if err != nil {
		log.Printf("ERROR GET ONE : %+v\n", err)
	}
	fmt.Printf("HERE WE GO AGAIN : %+v\n", temp)

	fmt.Printf("GOT ONE ITEM : %+v\n", resp)
	fmt.Println("END GetOneItem: ", time.Since(start))

	fmt.Println("************************************")
	return
}
func delOne(r *redis.Client, item string) {
	rc := r_manager.NewAdminRedisManager(r)

	err := rc.DeleteOne(ctx, item)
	if err != nil {
		log.Println("Errr deleteing : ", err)
	}

	fmt.Println("************************************")
}

func delMany(r *redis.Client, items []string) {
	rc := r_manager.NewAdminRedisManager(r)
	err := rc.DeleteMany(ctx, items)
	if err != nil {
		log.Panicf("error del many : %+v", err)
	}

	fmt.Println("************************************")
}

func updateOne(r *redis.Client, item, field string, val interface{}) {
	rc := r_manager.NewAdminRedisManager(r)

	err := rc.UpdateOne(ctx, item, field, val)
	if err != nil {
		log.Panicln("error update one: ", err)
	}
	fmt.Println("************************************")
}

func updateMany(r *redis.Client, item []map[string]interface{}) {
	rc := r_manager.NewAdminRedisManager(r)

	err := rc.UpdateMany(ctx, item)
	if err != nil {
		log.Panicln("error update one: ", err)
	}
	fmt.Println("************************************")
}

func deductMany(r *redis.Client, item []map[string]interface{}) {
	rc := r_manager.NewAdminRedisManager(r)
	toDeductItems := []map[string]interface{}{
		{"item": "cheese tart", "quantity": 500},
		{"item": "lemon tart", "quantity": 500},
		{"item": "egg tart", "quantity": 500},
	}
	err := rc.DeductQuantities(ctx, toDeductItems)
	if err != nil {
		log.Println("ERROR DEDUCT : ", err)
	}
	fmt.Println("************************************")
}
