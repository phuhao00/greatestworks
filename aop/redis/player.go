package redis

import "strconv"

func MakeAccountKey(userid int64) string {
	return "account:" + strconv.FormatInt(userid, 10)
}
