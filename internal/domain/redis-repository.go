package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

func InsertRedisPurchaseSummary(key string, value PurchaseSummary) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	txn := c.TxPipeline()

	err = txn.SetNX(context.Background(), key, data, 0).Err()
	if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set value: %v", err)
	}

	_, err = txn.Exec(context.Background())
	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return err
}

func FindByIDRedisPurchaseSumary(client *redis.Client, id uuid.UUID) (PurchaseSummary, error) {

	value, err := client.Get(context.Background(), id.String()).Result()
	if err == redis.Nil {
		return PurchaseSummary{}, fmt.Errorf("value not found: %v", err)
	} else if err != nil {
		return PurchaseSummary{}, fmt.Errorf("failed to get value: %v", err)
	}

	var purchaseSummary PurchaseSummary
	err = json.Unmarshal([]byte(value), &purchaseSummary)
	if err != nil {
		return PurchaseSummary{}, fmt.Errorf("failed to unmarshal value: %v", err)
	}

	return purchaseSummary, nil

}
