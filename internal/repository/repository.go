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
}
