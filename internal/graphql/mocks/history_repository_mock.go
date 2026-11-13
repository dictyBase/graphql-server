package mocks

import (
	"sync"
	"time"

	"github.com/dictyBase/graphql-server/internal/repository"
)

// RecentMemberCall records one AddRecentMember call.
type RecentMemberCall struct {
	Key      string
	Member   string
	Score    int64
	Capacity int64
}

// HistoryRepository is a repository.Repository fake that records
// AddRecentMember calls and returns Err from them. RecentMembers records
// each count in ReadCounts, returns ReadErr when it is set, and otherwise
// returns the first count items of Members. All other methods return zero
// values.
type HistoryRepository struct {
	Err        error
	Members    []repository.ScoredMember
	ReadErr    error
	mu         sync.Mutex
	Calls      []RecentMemberCall
	ReadCounts []int64
}

func (h *HistoryRepository) Get(string) (string, error)  { return "", nil }
func (h *HistoryRepository) Set(string, string) error    { return nil }
func (h *HistoryRepository) Exists(string) (bool, error) { return false, nil }
func (h *HistoryRepository) HSet(string, string, string) error {
	return nil
}

func (h *HistoryRepository) HGet(string, string) (string, error) {
	return "", nil
}

func (h *HistoryRepository) HExists(string, string) (bool, error) {
	return false, nil
}

func (h *HistoryRepository) SetWithTTL(
	string, string, time.Duration,
) error {
	return nil
}

// AddRecentMember records the call under the mutex and returns Err.
func (h *HistoryRepository) AddRecentMember(
	key, member string,
	score, capacity int64,
) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Calls = append(h.Calls, RecentMemberCall{
		Key: key, Member: member, Score: score, Capacity: capacity,
	})
	return h.Err
}

// RecentMembers records the count and returns ReadErr or the first count
// items of Members.
func (h *HistoryRepository) RecentMembers(
	_ string,
	count int64,
) ([]repository.ScoredMember, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.ReadCounts = append(h.ReadCounts, count)
	if h.ReadErr != nil {
		return nil, h.ReadErr
	}
	return h.Members[:min(int(count), len(h.Members))], nil
}
