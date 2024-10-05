package util

import (
	"github.com/google/uuid"
	"github.com/rs/xid"
	"strconv"
)

// NewXID The function "NewXID" generates a new unique identifier (XID) and returns it as a string.
func NewXID() string {
	return xid.New().String()
}

// MustNewUUID The function generates a new UUID and panics if there is an error.
func MustNewUUID() string {
	v, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return v.String()
}

func IDToString(id int64) string {
	return strconv.FormatInt(id, 10)
}

func IDToInt64(idStr string) int64 {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	return id
}
