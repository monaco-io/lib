package rs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type pipelineGetRs[T, K any] interface {
	Sugar(context.Context) (map[string]*T, error)
	get(context.Context) (hits map[string]*T, miss map[string]K, err error)
	set(context.Context, map[string]*T) error
}

type PipelineGetGetter[T, K any] func(context.Context, []K) (map[string]*T, error)

// PipelineGetJson pipeline get 命令, 注意控制key的数量
type PipelineGetJson[T, K any] struct {
	Conn    ICache
	KeysMap map[string]K  // 缓存key map, value是每个key回源需要的参数
	Expire  time.Duration // 过期时间
	Getter  MGetter[T, K] // 回源方法，如果回源没有找到数据，缓存默认存空字符串
}

var _ pipelineGetRs[any, any] = (*PipelineGetJson[any, any])(nil)

// Sugar 批量获取缓存消息
func (p PipelineGetJson[T, K]) Sugar(ctx context.Context) (map[string]*T, error) {
	if len(p.KeysMap) == 0 {
		return nil, nil
	}
	if p.Expire <= 0 {
		return nil, errors.New("cache key mush has expire time")
	}
	hits, miss, err := p.get(ctx)
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
		got, err := p.Getter(ctx, missSource)
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
		return fillData, p.set(ctx, fillData)
	})
	if err != nil {
		return nil, fmt.Errorf("PipelineGetJson Sugar.sg.Do: %w", err)
	}
	dt, ok := sgData.(map[string]*T)
	if !ok {
		return nil, errors.New("PipelineGetJson sgData is not map[string]*T")
	}
	for k, v := range dt {
		hits[k] = v
	}
	return hits, nil
}

// get 获取缓存
func (p PipelineGetJson[T, K]) get(ctx context.Context) (hits map[string]*T, miss map[string]K, err error) {
	var (
		pipeline  = p.Conn.Pipeline()
		resultMap = map[string]*redis.StringCmd{}
	)
	miss = make(map[string]K)
	hits = make(map[string]*T)

	for k := range p.KeysMap {
		resultMap[k] = pipeline.Get(ctx, k)
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		if !errors.Is(err, RedisNil) {
			err = fmt.Errorf("PipelineGetJson pipeline.Exec: %w", err)
			return
		}
		err = nil
	}
	for k, result := range resultMap {
		resStr, rErr := result.Result()
		if rErr != nil {
			if errors.Is(rErr, RedisNil) {
				miss[k] = p.KeysMap[k]
				continue
			}
			err = fmt.Errorf("PipelineGetJson result.Result(): %w", rErr)
			return
		}
		var data *T
		if resStr != "" {
			if err = json.Unmarshal([]byte(resStr), &data); err != nil {
				err = fmt.Errorf("PipelineGetJson json.Unmarshal: %w", err)
				return
			}
		}
		hits[k] = data
	}
	return
}

// set 设置缓存
func (p PipelineGetJson[T, K]) set(ctx context.Context, data map[string]*T) error {
	if len(data) == 0 {
		return errors.New("data should not be empty")
	}
	pipeline := p.Conn.Pipeline()
	for k, v := range data {
		var str string
		if v != nil {
			b, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("PipelineGetJson.set val can not format to json: %w", err)
			}
			str = string(b)
		}
		err := pipeline.Set(ctx, k, str, p.Expire).Err()
		if err != nil {
			return fmt.Errorf("PipelineGetJson.set pipeline.Set(k, str, p.Expire).Err(): %w", err)
		}
	}
	_, err := pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("PipelineGetJson.set pipeline.Exec: %w", err)
	}
	return nil
}
