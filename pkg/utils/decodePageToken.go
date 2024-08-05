package utils

import (
	"encoding/base64"
	"strconv"
)

func DecodePageToken(token string) (int, error) {
	bytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(bytes))
}
