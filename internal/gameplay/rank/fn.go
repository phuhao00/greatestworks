package rank

import "time"

var (
	startTime int64
	finalTime int64
)

func calScore(score int64, timeUnit int64, timeBitLen uint32, sortType uint32) int64 {
	var (
		scoreWithTime int64
		timeFactor    int64
	)

	nowTime := time.Now().Unix()

	if sortType == 1 {
		timeFactor = (nowTime - startTime) / timeUnit
	} else {
		timeFactor = (finalTime - nowTime) / timeUnit
	}

	scoreWithTime = (score << timeBitLen) | timeFactor

	return scoreWithTime
}

func getRealScore(tmScore int64, timeBitLen uint32) int64 {
	var realScore int64
	realScore = tmScore >> timeBitLen
	return realScore
}

func getRealScoreTime(tmScore int64, timeUnit int64, timeBitLen uint32, sortType uint32) int64 {
	var realTM int64
	timeFactor := tmScore & ((1 << timeBitLen) - 1)
	if sortType == 1 {
		realTM = (timeFactor * timeUnit) + startTime
	} else {
		realTM = finalTime - (timeFactor * timeUnit)
	}
	return realTM
}
