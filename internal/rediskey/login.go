package rediskey

import (
	"fmt"
	"strconv"
)

func MakeBanUserKey(userId uint64) string {
	return "ban:" + strconv.FormatUint(userId, 10)
}

func MakeAccountKey(userid int64) string {
	return "account:" + strconv.FormatInt(userid, 10)
}

func MakePlayerCacheKey(playerId uint64) string {
	return fmt.Sprintf("playerCache:%v", playerId)
}

func MakeTokenKey(userid uint64) string {
	return "token:" + strconv.FormatInt(int64(userid), 10)
}
