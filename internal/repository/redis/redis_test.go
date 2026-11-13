package redis

import (
	"fmt"
	"os"
	"sync"
	"testing"

	r "github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

var redisAddr = fmt.Sprintf(
	"%s:%s",
	os.Getenv("REDIS_SERVICE_HOST"),
	os.Getenv("REDIS_SERVICE_PORT"),
)

const (
	testHash = "testHash"
)

// CheckRedisEnv checks for the presence of the following
// environment variables
//
//	REDIS_SERVICE_HOST
//	REDIS_SERVICE_PORT
func CheckRedisEnv() error {
	envs := []string{
		"REDIS_SERVICE_HOST",
		"REDIS_SERVICE_PORT",
	}
	for _, e := range envs {
		if len(os.Getenv(e)) == 0 {
			return fmt.Errorf("env %s is not set", e)
		}
	}
	return nil
}

type TestRedis struct {
	client *r.Client
}

// NewTestRedisFromEnv is a constructor for a TestRedis instance.
func NewTestRedisFromEnv() (*TestRedis, error) {
	tr := new(TestRedis)
	if err := CheckRedisEnv(); err != nil {
		return tr, err
	}
	tr.client = r.NewClient(&r.Options{
		Addr: redisAddr,
	})
	err := tr.client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("error pinging redis %s", err)
	}
	return tr, nil
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func newTestCache(t *testing.T) *RedisStorage {
	t.Helper()
	repo, err := NewCache(redisAddr)
	if err != nil {
		t.Skipf("redis unavailable: %v", err)
	}
	redisStorage, ok := repo.(*RedisStorage)
	if !ok {
		t.Fatalf("expected *RedisStorage, got %T", repo)
	}
	return redisStorage
}

func TestSet(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	err := repo.Set("art", "vandelay")
	assert.NoError(err, "error in setting key")
}

func TestGet(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	err := repo.Set("art", "vandelay")
	assert.NoError(err, "error in setting key")
	key, err := repo.Get("art")
	assert.NoError(err, "error getting key")
	assert.Equal(key, "vandelay", "should retrieve correct value")
	_, err = repo.Get("newman")
	assert.Error(err, "should return error when key doesn't exist")
}

func TestExists(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	err := repo.Set("art", "vandelay")
	assert.NoError(err, "error in setting key")
	lookup, err := repo.Exists("art")
	assert.NoError(err, "error finding key")
	assert.True(lookup, "should find previously set key")
	badLookup, err := repo.Exists("obrien-murphy")
	assert.NoError(err, "error finding key")
	assert.False(badLookup, "should not find random key")
}

func TestHSet(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	err := repo.HSet(testHash, "art", "vandelay")
	assert.NoError(err, "error in setting key")
}

func TestHGet(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	err := repo.HSet(testHash, "art", "vandelay")
	assert.NoError(err, "error in setting key")
	key, err := repo.HGet(testHash, "art")
	assert.NoError(err, "error getting key")
	assert.Equal(key, "vandelay", "should retrieve correct value")
	_, err = repo.HGet("fakeHash", "banana")
	assert.Error(err, "should return error when hash doesn't exist")
}

func TestHExists(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	err := repo.HSet(testHash, "art", "vandelay")
	assert.NoError(err, "error in setting key")
	lookup, err := repo.HExists(testHash, "art")
	assert.NoError(err, "error finding key")
	assert.True(lookup, "should find previously set key")
	badLookup, err := repo.HExists(testHash, "obrien-murphy")
	assert.NoError(err, "error finding key")
	assert.False(badLookup, "should not find random key")
}

const testRecentCapacity = int64(250)

// newRecentKey returns an isolated sorted-set key that is removed
// when the test ends.
func newRecentKey(t *testing.T, repo *RedisStorage) string {
	t.Helper()
	key := "test:recent:" + t.Name()
	t.Cleanup(func() { repo.client.Del(key) })
	repo.client.Del(key)
	return key
}

func recentMembers(
	t *testing.T,
	repo *RedisStorage,
	key string,
) []string {
	t.Helper()
	members, err := repo.client.ZRange(key, 0, -1).Result()
	assert.NoError(t, err, "error reading members")
	return members
}

func TestAddRecentMemberInsert(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	assert.NoError(repo.AddRecentMember(key, "csa", 100, testRecentCapacity))
	got, err := repo.client.ZRangeWithScores(key, 0, -1).Result()
	assert.NoError(err)
	assert.Equal([]r.Z{{Score: 100, Member: "csa"}}, got)
}

func TestAddRecentMemberRepeatUpdatesScore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	assert.NoError(repo.AddRecentMember(key, "csa", 100, testRecentCapacity))
	assert.NoError(repo.AddRecentMember(key, "csa", 200, testRecentCapacity))
	card, err := repo.client.ZCard(key).Result()
	assert.NoError(err)
	assert.Equal(int64(1), card)
	score, err := repo.client.ZScore(key, "csa").Result()
	assert.NoError(err)
	assert.Equal(float64(200), score)
}

func TestAddRecentMemberLastWriteWins(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	assert.NoError(repo.AddRecentMember(key, "csa", 200, testRecentCapacity))
	assert.NoError(repo.AddRecentMember(key, "csa", 100, testRecentCapacity))
	card, err := repo.client.ZCard(key).Result()
	assert.NoError(err)
	assert.Equal(int64(1), card)
	score, err := repo.client.ZScore(key, "csa").Result()
	assert.NoError(err)
	assert.Equal(float64(100), score)
}

func TestAddRecentMemberTrimsToCapacity(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	for i := int64(1); i <= testRecentCapacity+1; i++ {
		assert.NoError(repo.AddRecentMember(
			key, fmt.Sprintf("q%03d", i), i, testRecentCapacity,
		))
	}
	card, err := repo.client.ZCard(key).Result()
	assert.NoError(err)
	assert.Equal(testRecentCapacity, card)
	_, err = repo.client.ZScore(key, "q001").Result()
	assert.ErrorIs(err, r.Nil, "oldest member must be evicted")
	newest, err := repo.client.ZRevRange(key, 0, 0).Result()
	assert.NoError(err)
	assert.Equal([]string{"q251"}, newest)
	oldest, err := repo.client.ZRange(key, 0, 0).Result()
	assert.NoError(err)
	assert.Equal([]string{"q002"}, oldest)
}

func TestAddRecentMemberRepeatRescuesFromEviction(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	for _, w := range []struct {
		member string
		score  int64
	}{{"a", 1}, {"b", 2}, {"c", 3}, {"a", 4}, {"d", 5}} {
		assert.NoError(repo.AddRecentMember(key, w.member, w.score, 3))
	}
	assert.ElementsMatch([]string{"a", "c", "d"}, recentMembers(t, repo, key))
}

func TestAddRecentMemberEqualScoreEviction(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	for _, m := range []string{"a", "b", "c"} {
		assert.NoError(repo.AddRecentMember(key, m, 100, 2))
	}
	assert.Equal([]string{"b", "c"}, recentMembers(t, repo, key))
}

func TestAddRecentMemberConcurrentSameMember(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			assert.NoError(
				repo.AddRecentMember(key, "csa", int64(i+1), testRecentCapacity),
			)
		}(i)
	}
	wg.Wait()
	card, err := repo.client.ZCard(key).Result()
	assert.NoError(err)
	assert.Equal(int64(1), card)
	score, err := repo.client.ZScore(key, "csa").Result()
	assert.NoError(err)
	assert.GreaterOrEqual(score, float64(1))
	assert.LessOrEqual(score, float64(50))
}

func TestAddRecentMemberConcurrentCapacity(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	repo := newTestCache(t)
	key := newRecentKey(t, repo)
	var wg sync.WaitGroup
	for i := range 300 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			assert.NoError(repo.AddRecentMember(
				key, fmt.Sprintf("q%03d", i), int64(i), testRecentCapacity,
			))
		}(i)
	}
	wg.Wait()
	card, err := repo.client.ZCard(key).Result()
	assert.NoError(err)
	assert.Equal(testRecentCapacity, card)
	oldest, err := repo.client.ZRange(key, 0, 0).Result()
	assert.NoError(err)
	assert.Equal([]string{"q050"}, oldest, "scores 0..49 must be evicted")
}
