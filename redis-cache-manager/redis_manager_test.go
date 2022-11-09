package redismanager

import (
	"context"
	"fmt"
	"testing"

	// rpc "github.com/ZAF07/tigerlily-e-bakery-cache/rpc"
	"github.com/ZAF07/tigerlily-e-bakery-inventories/api/rpc"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
var ctx = context.Background()

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

func TestAddInventory(t *testing.T) {
	manager := NewRedisManager(rdb)
	passed := true
	if err := manager.AddInventory(ctx, cheese); err != nil {
		passed = false
	}
	assert.True(t, passed, "AddInventory redis client passed")
}

func TestAddInventories(t *testing.T) {
	itemsToAdd := []*rpc.Sku{}
	itemsToAdd = append(itemsToAdd, lemon)
	itemsToAdd = append(itemsToAdd, egg)
	itemsToAdd = append(itemsToAdd, cheese)
	manager := NewRedisManager(rdb)
	passed := true
	if err := manager.AddInventories(ctx, itemsToAdd); err != nil {
		passed = false
	}
	assert.True(t, passed, "AddInventory redis client passed")
}

func TestDeductQuantity(t *testing.T) {
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
	if err := manager.DeductQuantity(ctx, inventory.Name, 1); err != nil {
		passed = false
	}
	assert.True(t, passed, "AddInventory redis client passed")
}

func TestGetAllInventories(t *testing.T) {
	itemsToGet := []*rpc.Sku{}
	itemsToGet = append(itemsToGet, lemon)
	itemsToGet = append(itemsToGet, egg)
	itemsToGet = append(itemsToGet, cheese)
	manager := NewRedisManager(rdb)
	passed := true
	resp, err := manager.GetAllInventories(ctx, itemsToGet)
	if err != nil {
		passed = false
	}
	if resp == nil || len(resp.Inventories) < 1 {
		passed = false
	}
	msg := fmt.Sprintf("Test for GetAllInventories passed: %t", passed)
	assert.True(t, passed, msg)
}

func TestGetOneItem(t *testing.T) {

	manager := NewAdminRedisManager(rdb)
	passed := true
	resp, err := manager.GetOneInventory(ctx, cheese.Name)
	if err != nil {
		passed = false
	}
	if resp == nil {
		passed = false
	}
	msg := fmt.Sprintf("Test for GetOneItem passed: %t", passed)
	assert.True(t, passed, msg)
}

func TestDeleteOneItem(t *testing.T) {

	manager := NewAdminRedisManager(rdb)
	passed := true
	err := manager.DeleteOne(ctx, cheese.Name)
	if err != nil {
		passed = false
	}

	msg := fmt.Sprintf("Test for DeleteOne passed: %t", passed)
	assert.True(t, passed, msg)
}

func TestDeleteMany(t *testing.T) {

	manager := NewAdminRedisManager(rdb)
	passed := true
	err := manager.DeleteMany(ctx, []string{"lemon tart", "cheese tart", "egg tart"})
	if err != nil {
		passed = false
	}

	msg := fmt.Sprintf("Test for DeleteMany passed: %t", passed)
	assert.True(t, passed, msg)
}

func TestUpdateone(t *testing.T) {

	manager := NewAdminRedisManager(rdb)
	passed := true
	err := manager.UpdateOne(ctx, "lemon tart", "price", 200)
	if err != nil {
		passed = false
	}

	msg := fmt.Sprintf("Test for UpdateOne passed: %t", passed)
	assert.True(t, passed, msg)
}

func TestUpdateMany(t *testing.T) {

	manager := NewAdminRedisManager(rdb)
	passed := true

	toUpdate := []map[string]interface{}{
		{"item": "lemon tart", "key": "price", "value": 500},
		{"item": "egg tart", "key": "price", "value": 500},
		{"item": "cheese tart", "key": "price", "value": 500},
	}

	err := manager.UpdateMany(ctx, toUpdate)
	if err != nil {
		passed = false
	}

	msg := fmt.Sprintf("Test for UpdateMany passed: %t", passed)
	assert.True(t, passed, msg)
}

func TestDeductMany(t *testing.T) {

	manager := NewAdminRedisManager(rdb)
	passed := true

	toDeductItems := []map[string]interface{}{
		{"item": "cheese tart", "quantity": 500},
		{"item": "lemon tart", "quantity": 500},
		{"item": "egg tart", "quantity": 500},
	}

	err := manager.DeductQuantities(ctx, toDeductItems)
	if err != nil {
		passed = false
	}

	msg := fmt.Sprintf("Test for UpdateMany passed: %t", passed)
	assert.True(t, passed, msg)
}
