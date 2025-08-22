package rs

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// returningToTheSource 是一个接口，用于缓存和构建某种类型的数据。
type returningToTheSource[T any] interface {
	// Sugar 是函数的公共入口点，它将尝试从缓存获取数据，
	// 如果不存在，将构建数据并设置缓存。
	Sugar(context.Context) (*T, error)

	// Get 是一个公有函数，用于从缓存获取数据。
	Get(context.Context) (*T, bool, error)

	// private methods
	// get 是一个私有函数，用于从缓存获取数据。
	get(context.Context) (*T, bool, error)

	// set 是一个私有函数，用于将数据设置到缓存中。
	set(context.Context, *T) error
}

// 确保 JSON 结构实现了 ReturningToTheSource 接口
var _ returningToTheSource[any] = (*JSON[any])(nil)

type Getter[T any] func() (*T, error)

// JSON 结构包含与 Redis 缓存的连接、操作的 key、
// 数据失效时间和构建缓存数据的 Builder 函数。
type JSON[T any] struct {
	Conn   ICache
	Key    string
	Expire time.Duration // 默认缓存时间 1 小时
	Getter Getter[T]
}

// Sugar 函数首先尝试从缓存中获得数据；如果数据不存在，
// 它会调用提供的 Builder 函数来构建数据，并将其存储在缓存中。
func (m *JSON[T]) Sugar(ctx context.Context) (*T, error) {
	data, ok, err := m.get(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}
	getKey := fmt.Sprintf("rs-sg-json-%s", m.Key)
	sgData, err, _ := sg.Do(getKey, func() (any, error) {
		got, err := m.Getter()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, m.set(ctx, nil)
			}
			return nil, err
		}
		return got, m.set(ctx, got)
	})
	if err != nil {
		return nil, fmt.Errorf("Sugar.sg.Do: %w", err)
	}
	if sgData == nil {
		return (*T)(nil), nil
	}
	dt, ok := sgData.(*T)
	if !ok {
		return nil, errors.New("sgData is not *T")
	}
	return dt, nil
}

func (m *JSON[T]) Get(ctx context.Context) (data *T, ok bool, err error) {
	return m.get(ctx)
}

// get 函数从 Redis 缓存中获取给定键的值（如果存在），
// 然后将该值的 JSON 数据解析为给定的 data 结构。
func (m *JSON[T]) get(ctx context.Context) (data *T, ok bool, err error) {
	val, err := m.Conn.Get(ctx, m.Key).Result()
	if err != nil {
		if errors.Is(err, RedisNil) {
			err = nil
			return
		}
		err = fmt.Errorf("Model.Get: %w", err)
		return
	}
	ok = true
	if val == "" {
		return
	}
	if err = json.Unmarshal([]byte(val), &data); err != nil {
		err = fmt.Errorf("Model.Get.JsonUnmarshal: %w", err)
		return
	}
	return
}

// set 函数将给定的值（以JSON 格式）存储到 Redis 缓存中。
func (m *JSON[T]) set(ctx context.Context, val *T) error {
	if m.Expire == 0 {
		return errors.New("cache key mush has expire time")
	}
	var str string
	if val != nil {
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Errorf("val can not format to json: %w", err)
		}
		str = string(b)
	}
	if err := m.Conn.Set(ctx, m.Key, str, m.Expire).Err(); err != nil {
		return fmt.Errorf("Model.Set: %w", err)
	}
	return nil
}
