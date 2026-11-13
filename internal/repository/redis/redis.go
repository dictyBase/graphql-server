package redis

import (
	"fmt"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	"github.com/dictyBase/graphql-server/internal/repository"
	r "github.com/go-redis/redis/v7"
)

type RedisStorage struct {
	client *r.Client
}

func NewCache(addr string) (repository.Repository, error) {
	client := r.NewClient(&r.Options{
		Addr: addr,
	})
	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("error pinging redis %s", err)
	}
	return &RedisStorage{client: client}, nil
}

func (rs *RedisStorage) Get(key string) (string, error) {
	val, err := rs.client.Get(key).Result()
	if err != nil {
		return "", fmt.Errorf("did not find %s in cache %s", key, err)
	}
	return val, nil
}

func (rs *RedisStorage) Set(key string, value string) error {
	return rs.client.Set(key, value, 24*time.Hour).Err()
}

func (rs *RedisStorage) SetWithTTL(
	key, value string,
	expiration time.Duration,
) error {
	err := rs.client.Set(key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("error in setting key with expiration %s", err)
	}

	return nil
}

func (rs *RedisStorage) Exists(key string) (bool, error) {
	h, err := rs.client.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if h == 0 {
		return false, nil
	}
	return true, nil
}

func (rs *RedisStorage) HGet(hash, key string) (string, error) {
	val, err := rs.client.HGet(hash, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rs *RedisStorage) HSet(hash, key, value string) error {
	return rs.client.HSet(hash, key, value).Err()
}

func (rs *RedisStorage) HExists(key, field string) (bool, error) {
	h, err := rs.client.HExists(key, field).Result()
	if err != nil {
		return false, err
	}
	return h, err
}

// AddRecentMember saves member in the sorted set at key with the given
// score, then trims the set to the capacity highest scores. A repeat of
// a member overwrites its score. ZADD and ZREMRANGEBYRANK run in one
// MULTI/EXEC transaction, so concurrent writes never leave the set over
// capacity.
func (rs *RedisStorage) AddRecentMember(
	key, member string,
	score, capacity int64,
) error {
	return F.Pipe3(
		IOE.TryCatchError(func() ([]r.Cmder, error) {
			return rs.client.TxPipelined(func(pipe r.Pipeliner) error {
				pipe.ZAdd(key, &r.Z{Score: float64(score), Member: member})
				pipe.ZRemRangeByRank(key, 0, -(capacity + 1))
				return nil
			})
		}),
		IOE.MapLeft[[]r.Cmder](
			ER.OnError("error in adding recent member to %s", key),
		),
		ioeutils.ToEither[error, []r.Cmder],
		E.Fold(F.Identity[error], F.Constant1[[]r.Cmder, error](nil)),
	)
}
