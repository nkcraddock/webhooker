package db

import (
	"encoding/json"
	"fmt"

	"github.com/nkcraddock/webhooker/domain"
	"gopkg.in/redis.v3"
)

type redisStore struct {
	getRedisClient func() *redis.Client
}

func RedisHookerStore(clientProvider func() *redis.Client) domain.Store {
	return &redisStore{
		getRedisClient: clientProvider,
	}
}

func (s *redisStore) SaveHook(h *domain.Hook) error {
	return s.save("hooks", h.Id, h)
}

func (s *redisStore) GetHook(id string) (*domain.Hook, error) {
	hook := new(domain.Hook)
	if err := s.get("hooks", id, hook); err != nil {
		return nil, err
	}
	return hook, nil
}

func (s *redisStore) GetHooks(query string) ([]*domain.Hook, error) {
	results, err := s.list("hooks")
	if err != nil {
		return nil, err
	}

	hooks := make([]*domain.Hook, len(results))
	i := 0

	for _, h := range results {
		hook := new(domain.Hook)
		if err = json.Unmarshal([]byte(h), hook); err != nil {
			return nil, err
		}

		hooks[i] = hook
		i += 1
	}

	return hooks, nil
}

func (s *redisStore) GetFilters(hook string) ([]*domain.Filter, error) {
	col := filterColKey(hook)
	results, err := s.list(col)
	if err != nil {
		return nil, err
	}

	filters := make([]*domain.Filter, len(results))
	i := 0

	for _, f := range results {
		filter := new(domain.Filter)
		if err = json.Unmarshal([]byte(f), filter); err != nil {
			return nil, err
		}

		filters[i] = filter
		i += 1
	}

	return filters, nil
}

func (s *redisStore) SaveFilter(f *domain.Filter) error {
	col := filterColKey(f.Hook)
	return s.save(col, f.Id, f)
}

func (s *redisStore) DeleteFilter(hook, id string) error {
	col := filterColKey(hook)
	return s.getRedisClient().HDel(col, id).Err()
}

func (s *redisStore) save(col, id string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.getRedisClient().HSet(col, id, string(jsonData)).Err()
}

func (s *redisStore) get(col, id string, data interface{}) error {
	jsonData, err := s.getRedisClient().HGet(col, id).Result()

	if err == redis.Nil {
		return domain.ErrorNotFound
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonData), data)
}

func (s *redisStore) list(col string) (map[string]string, error) {
	return s.getRedisClient().HGetAllMap(col).Result()
}

func (s *redisStore) push(col, val string) error {
	return s.getRedisClient().LPush(col, val).Err()
}

func filterColKey(hook string) string {
	return fmt.Sprintf("hook:%s:filters", hook)
}
