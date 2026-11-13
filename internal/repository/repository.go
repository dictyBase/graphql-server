package repository

import "time"

// ScoredMember is one member of a sorted set with its integer score.
type ScoredMember struct {
	Member string
	Score  int64
}

type Repository interface {
	Get(string) (string, error)
	Set(string, string) error
	SetWithTTL(string, string, time.Duration) error
	Exists(string) (bool, error)
	HGet(string, string) (string, error)
	HSet(string, string, string) error
	HExists(string, string) (bool, error)
	// AddRecentMember saves a member with a score in a sorted set and
	// trims the set to the capacity highest scores.
	AddRecentMember(key, member string, score, capacity int64) error
	// RecentMembers returns up to count members of a sorted set with the
	// highest scores, newest first. A count below 1 or a missing key gives
	// an empty non-nil slice.
	RecentMembers(key string, count int64) ([]ScoredMember, error)
}
