package redis

import (
	"fmt"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
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

var (
	zScoreLens = L.MakeLens(
		func(z r.Z) float64 { return z.Score },
		func(z r.Z, v float64) r.Z {
			z.Score = v
			return z
		},
	)
	// toScoredMember converts one go-redis member to the repository type.
	toScoredMember = func(z r.Z) repository.ScoredMember {
		return F.Pipe2(
			repository.ScoredMember{},
			scoredMemberNameLens.Set(fmt.Sprint(z.Member)),
			scoredMemberScoreLens.Set(int64(z.Score)),
		)
	}
	scoredMemberNameLens = L.MakeLens(
		func(m repository.ScoredMember) string { return m.Member },
		func(m repository.ScoredMember, v string) repository.ScoredMember {
			m.Member = v
			return m
		},
	)
	scoredMemberScoreLens = L.MakeLens(
		func(m repository.ScoredMember) int64 { return m.Score },
		func(m repository.ScoredMember, v int64) repository.ScoredMember {
			m.Score = v
			return m
		},
	)
	// nonNilScoredMembers replaces the nil slice that A.Map returns for an
	// empty input with an empty non-nil slice.
	nonNilScoredMembers = F.Flow2(
		O.FromPredicate(predarrays.IsNonEmpty[repository.ScoredMember]()),
		O.GetOrElse(func() []repository.ScoredMember {
			return []repository.ScoredMember{}
		}),
	)
	emptyRecentMembers = IOE.Of[error]([]repository.ScoredMember{})
	zMemberLens        = L.MakeLens(
		func(z r.Z) any { return z.Member },
		func(z r.Z, v any) r.Z {
			z.Member = v
			return z
		},
	)
)

// AddRecentMember saves member in the sorted set at key with the given
// score, then trims the set to the capacity highest scores. A repeat of
// a member overwrites its score. ZADD and ZREMRANGEBYRANK run in one
// MULTI/EXEC transaction, so concurrent writes never leave the set over
// capacity.
func (rs *RedisStorage) AddRecentMember(
	key, member string,
	score, capacity int64,
) error {
	return F.Pipe7(
		r.Z{},
		zScoreLens.Set(float64(score)),
		zMemberLens.Set(member),
		F.Ref[r.Z],
		rs.addCapped(key, capacity),
		IOE.MapLeft[[]r.Cmder](
			ER.OnError("error in adding recent member to %s", key),
		),
		ioeutils.ToEither[error, []r.Cmder],
		E.Fold(F.Identity[error], F.Constant1[[]r.Cmder, error](nil)),
	)
}

// addCapped returns the effect that runs ZADD for one member and then
// ZREMRANGEBYRANK in one MULTI/EXEC transaction.
func (rs *RedisStorage) addCapped(
	key string,
	capacity int64,
) func(*r.Z) IOE.IOEither[error, []r.Cmder] {
	return func(z *r.Z) IOE.IOEither[error, []r.Cmder] {
		return IOE.TryCatchError(func() ([]r.Cmder, error) {
			return rs.client.TxPipelined(func(pipe r.Pipeliner) error {
				pipe.ZAdd(key, z)
				pipe.ZRemRangeByRank(key, 0, -(capacity + 1))
				return nil
			})
		})
	}
}

// RecentMembers returns up to count members of the sorted set at key with
// the highest scores, newest first. A count below 1 sends no command and
// gives an empty non-nil slice. A missing key gives the same.
func (rs *RedisStorage) RecentMembers(
	key string,
	count int64,
) ([]repository.ScoredMember, error) {
	result := F.Pipe4(
		count,
		O.FromPredicate(N.MoreThan(int64(0))),
		O.Fold(
			F.Constant(emptyRecentMembers),
			rs.readRecent(key),
		),
		IOE.MapLeft[[]repository.ScoredMember](
			ER.OnError("error in reading recent members of %s", key),
		),
		ioeutils.ToEither[error, []repository.ScoredMember],
	)

	return E.UnwrapError(result)
}

// readRecent returns the effect that runs ZREVRANGE WITHSCORES for the
// count newest members.
func (rs *RedisStorage) readRecent(
	key string,
) func(int64) IOE.IOEither[error, []repository.ScoredMember] {
	return func(count int64) IOE.IOEither[error, []repository.ScoredMember] {
		return F.Pipe2(
			IOE.TryCatchError(func() ([]r.Z, error) {
				return rs.client.ZRevRangeWithScores(key, 0, count-1).Result()
			}),
			IOE.Map[error](A.Map(toScoredMember)),
			IOE.Map[error](nonNilScoredMembers),
		)
	}
}
