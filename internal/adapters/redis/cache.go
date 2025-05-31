package myredis

import (
	"LostAndFound/internal/domain/entity"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
}

func (c *CacheRepository) SetUserData(ctx context.Context, user *entity.User) error {
	fields := map[string]any{
		"email":   user.Email,
		"name":    user.Name,
		"surname": user.Surname,
	}

	fields["phone"] = user.Phone

	key := "user:" + user.ID
	return c.client.HSet(ctx, key, fields, 30*time.Minute).Err()
}

func (c *CacheRepository) GetUserData(ctx context.Context, userID string) (entity.User, error) {
	key := "user:" + userID
	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return entity.User{}, err
	}

	if len(data) == 0 {
		return entity.User{}, redis.Nil
	}

	return entity.User{
		Email:   data["email"],
		Name:    data["name"],
		Surname: data["surname"],
		Phone:   data["phone"],
	}, nil
}

func (c *CacheRepository) DeleteUserData(ctx context.Context, userID string) error {
	return c.client.Del(ctx, "user:"+userID).Err()
}

func (c *CacheRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return c.client.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}

func (c *CacheRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := c.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (c *CacheRepository) SaveCard(ctx context.Context, card *entity.Card) error {
	data, err := json.Marshal(card)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "card:"+card.ID, data, 10*time.Minute).Err()
}

func (c *CacheRepository) GetCardByID(ctx context.Context, id string) (*entity.Card, error) {
	data, err := c.client.Get(ctx, "card:"+id).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var card entity.Card
	if err = json.Unmarshal([]byte(data), &card); err != nil {
		return nil, err
	}
	return &card, nil
}

func (c *CacheRepository) DeleteCard(ctx context.Context, id string) error {
	return c.client.Del(ctx, "card:"+id).Err()
}

func NewCacheRepo(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}
