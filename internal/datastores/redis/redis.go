package redis

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/trelore/todoapi/internal"
	"go.uber.org/zap"
)

type r struct {
	redisClient *redis.Client
	sugar       *zap.SugaredLogger
}

func New(sugar *zap.SugaredLogger) *r {
	client := redis.NewClient(&redis.Options{
		Addr:     "db:6379",
		Password: "",
		DB:       0,
	})
	return &r{
		redisClient: client,
		sugar:       sugar,
	}
}

// Insert implements the interface
func (r *r) Insert(description string) (*internal.Item, error) {
	uuid := uuid.New()
	item := internal.Item{Description: description, ID: uuid, Done: false}
	b, err := json.Marshal(&item)
	if err != nil {
		return nil, err
	}
	if status := r.redisClient.Set(uuid.String(), b, 0); status.Err() != nil {
		return nil, status.Err()
	}
	return &item, nil
}

// List implements the interface
func (r *r) List() ([]*internal.Item, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// Get implements the interface
func (r *r) Get(id string) (*internal.Item, error) {
	res, err := r.redisClient.Get(id).Result()
	if err != nil {
		return nil, err
	}
	var item internal.Item
	err = json.Unmarshal([]byte(res), &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Delete implements the interface
func (r *r) Delete(id string) error {
	_, err := r.redisClient.Del(id).Result()
	if err != nil {
		return err
	}
	return nil
}

// Upsert implements the interface
func (r *r) Upsert(id string, item *internal.Item) (*internal.Item, error) {
	b, err := json.Marshal(&item)
	if err != nil {
		return nil, err
	}
	if _, err := r.redisClient.Set(item.ID.String(), b, 0).Result(); err != nil {
		return nil, err
	}
	return item, nil
}