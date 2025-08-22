package rs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type mGetRs[T, K any] interface {
	Sugar(context.Context) (map[string]*T, error)
	get(context.Context) (hits map[string]*T, miss map[string]K, err error)
	set(context.Context, map[string]*T) error
}

type MGetter[T, K any] func(context.Context, []K) (map[string]*T, error)

type MGetJson[T, K any] struct {
	Conn    ICache
	KeysMap map[string]K  // 缓存key map, value是每个key回源需要的参数
	Expire  time.Duration // 过期时间
	Getter  MGetter[T, K] // 回源方法，如果回源没有找到数据，缓存默认存空字符串

	keys []string
}

var _ mGetRs[any, any] = (*MGetJson[any, any])(nil)

// Sugar 批量获取缓存消息
func (m MGetJson[T, K]) Sugar(ctx context.Context) (map[string]*T, error) {
	if len(m.KeysMap) == 0 {
		return nil, nil
	}
	if m.Expire <= 0 {
		return nil, errors.New("cache key mush has expire time")
	}
	for k := range m.KeysMap {
		m.keys = append(m.keys, k)
	}

	hits, miss, err := m.get(ctx)
	if err != nil {
		return nil, err
	}
	if len(miss) == 0 {
		return hits, nil
	}
	// 回源
	var missSource []K
	var missKey []string
	for key := range miss {
		missKey = append(missKey, key)
		missSource = append(missSource, miss[key])
	}
	sgKey := strings.Join(missKey, "_")
	sgData, err, _ := sg.Do(sgKey, func() (any, error) {
		got, err := m.Getter(ctx, missSource)
		if err != nil {
			return nil, err
		}
		// 填充默认值
		var defaultValue *T
		fillData := make(map[string]*T, len(missSource))
		for key := range miss {
			vv, ok := got[key]
			if ok {
				fillData[key] = vv
			} else {
				fillData[key] = defaultValue
			}
		}
		return fillData, m.set(ctx, fillData)
	})
	if err != nil {
		return nil, fmt.Errorf("Sugar.sg.Do: %w", err)
	}
	dt, ok := sgData.(map[string]*T)
	if !ok {
		return nil, errors.New("sgData is not map[string]*T")
	}
	for k, v := range dt {
		hits[k] = v
	}
	return hits, nil
}

func (m MGetJson[T, K]) get(ctx context.Context) (hits map[string]*T, miss map[string]K, err error) {
	cacheResult, err := m.Conn.MGet(ctx, m.keys...).Result()
	if err != nil {
		err = fmt.Errorf("MGetJson m.Conn.MGet: %w", err)
		return nil, nil, err
	}
	miss = make(map[string]K)
	hits = make(map[string]*T)
	for i, v := range cacheResult {
		k := m.keys[i]
		if v == nil {
			miss[k] = m.KeysMap[k]
			continue
		}
		var data *T
		val, ok := v.(string)
		if !ok {
			err = fmt.Errorf("MGetJson val not string: %w", err)
			return
		}
		if val != "" {
			if err = json.Unmarshal([]byte(val), &data); err != nil {
				err = fmt.Errorf("MGetJson JsonUnmarshal: %w", err)
				return
			}
		}
		hits[k] = data
	}
	return
}

func (m MGetJson[T, K]) set(ctx context.Context, data map[string]*T) error {
	if len(data) == 0 {
		return nil
	}
	var pairs []string
	for k, v := range data {
		var str string
		if v != nil {
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("MGetJson.set val can not format to json: %w", err)
			}
			str = string(b)
		}
		pairs = append(pairs, k, str)
	}
	// 设置缓存
	err := m.Conn.MSet(ctx, pairs).Err()
	if err != nil {
		return fmt.Errorf("MGetJson.set mset: %w", err)
	}
	// 设置过期时间

	pipeline := m.Conn.Pipeline()
	for k := range data {
		err = pipeline.Expire(ctx, k, m.Expire).Err()
		if err != nil {
			return fmt.Errorf("MGetJson.set pipeline.Expire: %w", err)
		}
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("MGetJson.set pipeline.Exec: %w", err)
	}
	return nil
}
