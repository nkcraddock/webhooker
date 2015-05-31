package db

import (
	"encoding/json"

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
