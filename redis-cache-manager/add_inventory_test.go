package redismanager

import (
	"fmt"
	"testing"

	rpc "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

/*
	❌ Look into this. Does not work now
*/

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
var ctx = context.Background()

func TestAddInventory(t *testing.T) {
	inventory := &rpc.Sku{
		Name:        "cheese tart",
		Description: "Cheesy",
		Price:       2.5,
		Quantity:    10,
		SkuId:       "11111",
		ImageUrl:    "cheese_tart.com",
		Type:        "tart",
	}
	manager := NewRedisManager(rdb)
	passed := true
	if err := manager.AddInventory(ctx, inventory); err != nil {
		passed = false
	}
	assert.True(t, passed, "AddInventory redis client passed")
}

func TestAddInventories(t *testing.T) {
	lemon := &rpc.Sku{
		Name:        "lemon tart",
		Description: "Sweet and sour",
		Price:       2.5,
		Quantity:    10,
		SkuId:       "11111",
		ImageUrl:    "lemon_tart.com",
		Type:        "tart",
	}
	egg := &rpc.Sku{
		Name:        "egg tart",
		Description: "Eggy",
		Price:       2.5,
		Quantity:    10,
		SkuId:       "11111",
		ImageUrl:    "egg_tart.com",
		Type:        "tart",
	}
	cheese := &rpc.Sku{
		Name:        "cheese tart",
		Description: "Cheesy",
		Price:       2.5,
		Quantity:    10,
		SkuId:       "11111",
		ImageUrl:    "cheese_tart.com",
		Type:        "tart",
	}
	itemsToAdd := []*rpc.Sku{}
	itemsToAdd = append(itemsToAdd, lemon)
	itemsToAdd = append(itemsToAdd, egg)
	itemsToAdd = append(itemsToAdd, cheese)
	manager := NewRedisManager(rdb)
	var err error
	if err = manager.AddInventories(ctx, itemsToAdd); err != nil {
		fmt.Println("RUNNING")
	}
	assert.Error(t, err)
	// (t, passed, "AddInventory redis client passed")
}
