package userCache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pedramktb/schwarzit-probearbeit/internal/datasource"
	"github.com/pedramktb/schwarzit-probearbeit/internal/logging"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
	"go.uber.org/zap"
)

const (
	// CacheTTL is the time-to-live for cached user data
	cacheTTL           = 24 * time.Hour
	cacheUpdateTimeout = 5 * time.Second
)

type cache struct {
	*redis.Client
	datasource.Getter[types.User]
	datasource.Saver[types.User]
	datasource.Deleter[types.User]
	datasource.UserByEmailGetter
}

func create(
	r *redis.Client,
	getter datasource.Getter[types.User],
	saver datasource.Saver[types.User],
	deleter datasource.Deleter[types.User],
	userByEmailGetter datasource.UserByEmailGetter,
) *cache {
	return &cache{
		r,
		getter,
		saver,
		deleter,
		userByEmailGetter,
	}
}

func keyFromID(id uuid.UUID) string {
	return "user:" + id.String()
}

func (c *cache) Get(ctx context.Context, id uuid.UUID) (types.User, error) {
	cached, err := c.Client.Get(ctx, keyFromID(id)).Result()
	if err == nil {
		// Cache hit
		var user types.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return user, nil
		} else {
			logging.FromContext(ctx).Warn("failed to unmarshal cached user data", zap.String("cached_user", cached), zap.Error(err))
		}
	} else if !errors.Is(err, redis.Nil) {
		logging.FromContext(ctx).Warn("cache error", zap.Error(err))
	}

	// Cache miss, fetch from wrapped getter
	// TODO: use redisJSON for better performance
	user, err := c.Getter.Get(ctx, id)
	if err != nil {
		return user, err
	}

	c.setUserCache(user)

	return user, nil
}

func (c *cache) Save(ctx context.Context, user types.User) (types.User, error) {
	savedUser, err := c.Saver.Save(ctx, user)
	if err != nil {
		return savedUser, err
	}

	c.delUserCache(savedUser.ID, &user.Email)

	return savedUser, nil
}

func (c *cache) Delete(ctx context.Context, id uuid.UUID) error {
	err := c.Deleter.Delete(ctx, id)
	if err != nil {
		return err
	}

	c.delUserCache(id, nil)

	return nil
}

func (c *cache) GetByEmail(ctx context.Context, email string) (types.User, error) {
	cached, err := c.Client.Get(ctx, "user_id:"+email).Result()
	if err == nil {
		// Cache hit
		var user_id uuid.UUID
		if err := json.Unmarshal([]byte(cached), &user_id); err == nil {
			return c.Get(ctx, user_id)
		} else {
			logging.FromContext(ctx).Warn("failed to unmarshal cached user id", zap.String("cached_user_id", cached), zap.Error(err))
		}
	} else if !errors.Is(err, redis.Nil) {
		logging.FromContext(ctx).Warn("cache error", zap.Error(err))
	}

	// Cache miss, fetch from wrapped getter
	user, err := c.UserByEmailGetter.GetByEmail(ctx, email)
	if err != nil {
		return user, err
	}

	c.setUserCache(user)

	return user, nil
}

func (c *cache) setUserCache(user types.User) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheUpdateTimeout)
		defer cancel()
		// Cache the result
		userData, err := json.Marshal(user)
		if err != nil {
			logging.FromContext(ctx).Warn("failed to marshal user for caching", zap.Any("caching_user", user), zap.Error(err))
		}

		err = c.Client.Set(ctx, keyFromID(user.ID), userData, cacheTTL).Err()
		if err != nil {
			logging.FromContext(ctx).Warn("failed to cache user", zap.Any("caching_user", user), zap.Error(err))
		}

		// Cache the result id by email (less space than saving the whole user and add the ability to invalidate by id alone)
		// TODO: use redisJSON for better performance
		idData, err := json.Marshal(user.ID)
		if err != nil {
			logging.FromContext(ctx).Warn("failed to marshal user id for caching", zap.String("caching_user_id", user.ID.String()), zap.Error(err))
		}

		err = c.Client.Set(ctx, "user:email:"+user.Email, idData, cacheTTL).Err()
		if err != nil {
			logging.FromContext(ctx).Warn("failed to cache user id", zap.String("caching_user_id", user.ID.String()), zap.Error(err))
		}
	}()
}

func (c *cache) delUserCache(id uuid.UUID, email *string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cacheUpdateTimeout)
		defer cancel()

		// Invalidate cache
		if err := c.Client.Del(ctx, keyFromID(id)).Err(); !errors.Is(err, redis.Nil) && err != nil {
			logging.FromContext(ctx).Error("failed to invalidate cache", zap.Error(err))
		}

		// Invalidate cache by email to reduce 1 lookup
		if email != nil {
			if err := c.Client.Del(ctx, "user:email:"+*email).Err(); !errors.Is(err, redis.Nil) && err != nil {
				logging.FromContext(ctx).Error("failed to invalidate cache", zap.Error(err))
			}
		}
	}()
}
