package utils

import (
	"encoding/base64"
	"strconv"
)

func EncodePageToken(index int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(index)))
}
