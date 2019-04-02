package models

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/vektah/gqlgen/graphql"
)

func MarshalTimestamp(t time.Time) graphql.Marshaler {
	ct := t.Format(time.RFC3339Nano)
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(ct))
	})
}

func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if s, ok := v.(int); ok {
		return time.Unix(int64(s), 0), nil
	}
	return time.Time{}, errors.New("time should be a unix timestamp")
}
