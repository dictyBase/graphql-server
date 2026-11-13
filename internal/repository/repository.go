package repository

import "time"

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
}
