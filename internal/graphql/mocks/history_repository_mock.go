package mocks

import (
	"sync"
	"time"
)

// RecentMemberCall records one AddRecentMember call.
type RecentMemberCall struct {
	Key      string
	Member   string
	Score    int64
	Capacity int64
}

// HistoryRepository is a repository.Repository fake that records
// AddRecentMember calls and returns Err from them. All other methods
// return zero values.
type HistoryRepository struct {
	Err   error
	mu    sync.Mutex
	Calls []RecentMemberCall
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
