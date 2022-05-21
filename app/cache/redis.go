package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/service/user/app/config"
	"time"
)

type RedisCache interface {
	StartResetUserPassword(
		ctx context.Context,
		passwordActivationKey string,
		userId string,
	) error
	UserInPasswordReset(
		ctx context.Context,
		resetPasswordKey string,
	) (string, error)
	RegisterNewCompany(
		ctx context.Context,
		verificationKey string,
		companyId int64,
	) error
	CompanyInRegistration(
		ctx context.Context,
		verificationKey string,
	) (int64, error)
	CleanUpCompanyRegistration(
		ctx context.Context,
		verificationKey string,
	) error
}

type redisCache struct {
	redis  *redis.Client
	config *config.Config
	log    logger.Logger
}

func NewRedisCache(
	redis *redis.Client,
	config *config.Config,
	log logger.Logger,
) RedisCache {
	return &redisCache{
		redis:  redis,
		config: config,
		log:    log,
	}
}

type Entry struct {
	Val    string
	Exists bool
}

// setAsJSON converts the given value as JSON into the cache
// with the given key.
func (c *redisCache) setAsJSON(ctx context.Context, key string, value interface{}, t time.Duration) error {
	encodedData, err := json.Marshal(value)

	if err != nil {
		return err
	}

	err = c.redis.Set(ctx, key, encodedData, t).Err()

	return err
}

// getFromJSON gets the entry from the given key
// as JSON format and umarshal it to the given destination
// interface structure type
func (c *redisCache) getFromJSON(ctx context.Context, key string, dest interface{}) error {
	entry := c.redis.Get(ctx, key)

	switch {
	case entry.Err() == redis.Nil:
		return entry.Err()
	case entry.Err() != nil:
		return entry.Err()
	default:
		err := json.Unmarshal([]byte(entry.Val()), dest)
		return err
	}
}

// get the entry from the given key
// as Entry. If the entry does not exist, the Entry
// has the flag Exists set false and no error will be returned.
// Or if an error happens, an error will be returned.
// Otherwise, the Value and the Exists flag will be returned, without
// any error.
func (c *redisCache) get(ctx context.Context, key string) (Entry, error) {
	entry := c.redis.Get(ctx, key)

	sReturn := Entry{Exists: false}

	switch {
	case entry.Err() == redis.Nil:
		return sReturn, nil
	case entry.Err() != nil:
		return sReturn, entry.Err()
	default:
		sReturn.Exists = true
		sReturn.Val = entry.Val()
		return sReturn, nil
	}
}

func (c *redisCache) set(ctx context.Context, key string, value interface{}, t time.Duration) error {
	err := c.redis.Set(ctx, key, value, t).Err()
	return err
}

// FlushAll the cache and flush the db
func (c *redisCache) FlushAll() error {
	ctx := context.Background()
	return c.redis.FlushDB(ctx).Err()
}
