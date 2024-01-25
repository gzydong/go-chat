package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const tableCacheExpiration = 10 * time.Minute // 缓存过期时间常量

type TablePrimaryType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type IModelCache[S TablePrimaryType] interface {
	TableName() string
	TablePrimaryIdValue() S
}

type TableCache[T IModelCache[V], V TablePrimaryType] struct {
	redis     *redis.Client
	prefix    string
	tableName string
}

func NewTableCache[T IModelCache[V], V TablePrimaryType](redis *redis.Client) TableCache[T, V] {
	return TableCache[T, V]{redis: redis, prefix: "tablecache", tableName: any(new(T)).(IModelCache[V]).TableName()}
}

func (t *TableCache[T, V]) buildKey(primaryId V) string {
	return fmt.Sprintf("%s:%s:%d", t.prefix, t.tableName, primaryId)
}

func (t *TableCache[T, V]) Get(ctx context.Context, primaryId V) (*T, error) {
	value, err := t.redis.Get(ctx, t.buildKey(primaryId)).Result()
	if err != nil {
		return nil, err
	}

	var model T
	if err := json.Unmarshal([]byte(value), &model); err != nil {
		return nil, fmt.Errorf("unmarshal cache value: %w", err)
	}

	return &model, nil
}

func (t *TableCache[T, V]) Set(ctx context.Context, data T) error {
	bt, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	err = t.redis.Set(ctx, t.buildKey(data.TablePrimaryIdValue()), bt, tableCacheExpiration).Err()
	if err != nil {
		return fmt.Errorf("set cache: %w", err)
	}

	return nil
}

func (t *TableCache[T, V]) GetOrSet(ctx context.Context, primaryId V, loader func(ctx context.Context) (*T, error)) (*T, error) {
	model, err := t.Get(ctx, primaryId)
	if err == nil {
		return model, nil
	}

	model, err = loader(ctx)
	if err != nil {
		return nil, err
	}

	if err := t.Set(ctx, *model); err != nil {
		return nil, fmt.Errorf("get or set cache: %w", err)
	}

	return model, nil
}

func (t *TableCache[T, V]) Del(ctx context.Context, primaryId V) error {
	return t.redis.Del(ctx, t.buildKey(primaryId)).Err()
}

// CacheGetOrSet 缓存获取或设置
func CacheGetOrSet[T any](
	ctx context.Context,
	rds *redis.Client,
	key string,
	loader func(ctx context.Context) (*T, error),
	ex time.Duration,
) (*T, error) {
	value, err := rds.Get(ctx, key).Result()
	if err == nil {
		model := new(T)
		if err := json.Unmarshal([]byte(value), model); err == nil {
			return model, nil
		}
	}

	model, err := loader(ctx)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(model)
	if err == nil {
		if err := rds.Set(ctx, key, string(body), ex).Err(); err != nil {
			fmt.Println("CacheGetOrSet GetOrSet Err:", err)
		}
	}

	return model, nil
}
