package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	r_manager "github.com/ZAF07/tigerlily-e-bakery-cache/redis-cache-manager"
	"github.com/nitishm/go-rejson/v4"

	// Dont need cache/proto & submodules. Can delete
	// "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
	"github.com/ZAF07/tigerlily-e-bakery-inventories/api/rpc"
	"github.com/go-redis/redis/v8"
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

	rh := rejson.NewReJSONHandler()

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// p := Person{"Zafffere", 10}
	rh.SetGoRedisClient(rdb)
	// toAdd, merr := json.Marshal(p)
	// if merr != nil {
	// 	log.Fatalf("Marshal err : %+v\n", merr)
	// }
	// rh.JSONSet("person", ".", p)

	// 💡 TODO: Migrate redis to redisJSON.

	// GET BY KEY
	// g, gerr := rh.JSONGet("person", "name")

	// GET ENTIRE OBJECT
	g, gerr := rh.JSONGet("person", "name")
	if gerr != nil {
		log.Fatalf("Error getting from redis : %+v\n", gerr)
	}
	// jj := &Person{}
	// rr := json.Unmarshal(g.([]byte), jj)
	// if rr != nil {
	// 	log.Fatalf("Eror nmarshal : %+v\n", rr)
	// }
	fmt.Printf("HERE IS THE PERSON : %s\n", g.([]byte))
	// cache := r_manager.NewRedisManager(rdb)

	// itemsToAdd := []*rpc.Sku{}
	// itemsToAdd = append(itemsToAdd, lemon)
	// itemsToAdd = append(itemsToAdd, egg)
	// itemsToAdd = append(itemsToAdd, cheese)

	// bulkErr := cache.AddInventories(ctx, GetItems())
	// if bulkErr != nil {
	// 	log.Printf("BULD ADD ERROR : %+v", bulkErr)
	// }
	// fmt.Println("DONE BULK ADD")
	// fmt.Println("************************************")
	// item := &rpc.Sku{
	// 	Name:        "lemon tart",
	// 	Price:       2.4,
	// 	SkuId:       "001001",
	// 	Type:        "Tart",
	// 	ImageUrl:    "lemon_tart.com",
	// 	Quantity:    10,
	// 	Description: "Sweet & Sour",
	// }

	// // INSERT ONE ITEM TO CACHE
	// err := cache.AddInventory(ctx, item)
	// if err != nil {
	// 	log.Fatalf("Failed!! : %+v", err)
	// }
	// fmt.Println("DONE ADDING TO INVENTORY")
	// fmt.Println("************************************")

	// // DEDUCT ONE ITEM QUANTITY
	// if err := cache.DeductQuantity(ctx, item.Name, 8); err != nil {
	// 	log.Fatalf("ERROR DEDUCT : %+v", err)
	// }
	// fmt.Println("************************************")

	// // GET ALL INVENTORIES
	// items := []*rpc.Sku{}
	// items = append(items, item)
	// resp, err := cache.GetAllInventories(ctx, items)
	// if err != nil {
	// 	log.Printf("ERROR ALL : %+v\n", err)
	// }
	// fmt.Printf("ALL INVENTORIES : %+v\n", resp)

	// Get one item
	// GetOneItem(rdb, cheese, "price")
	// updateOne(rdb, "lemon tart", "price", 200)
	// updateMany(rdb, []map[string]interface{}{
	// 	{"item": "lemon tart", "key": "price", "value": 500},
	// 	{"item": "egg tart", "key": "price", "value": 500},
	// 	{"item": "cheese tart", "key": "price", "value": 500},
	// })
	// deductMany(rdb, []map[string]interface{}{})
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

func GetItems() []*rpc.Sku {
	var lemon = &rpc.Sku{
		Name:        "lemon tart",
		Description: "Sweet and sour",
		Price:       2.5,
		Quantity:    10,
		SkuId:       "12222",
		ImageUrl:    "lemon",
		Type:        "tart",
	}
	var oreoCake = &rpc.Sku{
		Name:        "orea cake",
		Description: "Eggy",
		Price:       8.5,
		Quantity:    10,
		SkuId:       "11222",
		ImageUrl:    "egg",
		Type:        "tart",
	}
	var cheeseBun = &rpc.Sku{
		Name:        "cheese bun",
		Description: "Cheesy",
		Price:       2.90,
		Quantity:    10,
		SkuId:       "10001",
		ImageUrl:    "cheese",
		Type:        "tart",
	}
	var lemonIceCream = &rpc.Sku{
		Name:        "lemon Ice Cream",
		Description: "Sweet and sour",
		Price:       9.5,
		Quantity:    87,
		SkuId:       "11000",
		ImageUrl:    "lemon",
		Type:        "tart",
	}
	var eggPie = &rpc.Sku{
		Name:        "egg pie",
		Description: "Eggy pie",
		Price:       21.5,
		Quantity:    80,
		SkuId:       "11100",
		ImageUrl:    "egg",
		Type:        "tart",
	}
	var cheesePie = &rpc.Sku{
		Name:        "cheese pie",
		Description: "Cheesy pie",
		Price:       7.5,
		Quantity:    190,
		SkuId:       "11110",
		ImageUrl:    "cheese",
		Type:        "tart",
	}
	var applePie = &rpc.Sku{
		Name:        "apple pie",
		Description: "Sweet and sour",
		Price:       19.5,
		Quantity:    87,
		SkuId:       "22000",
		ImageUrl:    "lemon",
		Type:        "tart",
	}
	var lemonSorbet = &rpc.Sku{
		Name:        "lemon sorbet",
		Description: "Iceee",
		Price:       11.5,
		Quantity:    80,
		SkuId:       "333300",
		ImageUrl:    "egg",
		Type:        "tart",
	}
	var roll = &rpc.Sku{
		Name:        "cheese roll",
		Description: "Cheesy roll",
		Price:       17.5,
		Quantity:    190,
		SkuId:       "13333",
		ImageUrl:    "lemon",
		Type:        "tart",
	}
	items := []*rpc.Sku{}
	items = append(items, lemon, roll, lemonSorbet, lemonIceCream, oreoCake, eggPie, egg, lemon, applePie, cheesePie, cheeseBun)
	b, e := json.Marshal(items)
	if e != nil {
		fmt.Println(":Error unmarslaing : ", e)
	}
	os.WriteFile("inventories.yaml", b, 0644)
	return items
}

// func TestProtoRedis() {
// 	// a, err := proto.Marshal()
// 	if err != nil {
// 		fmt.Println("Err : ", err)
// 	}
// }
