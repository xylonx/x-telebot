package util

import "strconv"

func GenerateID(chatID int64, username string) string {
	id := strconv.FormatInt(chatID, 10)
	return id + ":" + username
}
