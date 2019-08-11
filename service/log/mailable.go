package main

import (
	"time"
)

// TODO: 每五分钟检查记录中是否超时，如果超时，进行内存回收。
type mailAble struct {
	beginTime time.Time
	times     int
}

var mailAbleTimes = make(map[string]*mailAble)

// 允计发送短信的次数, 以防止日志爆量时的短信无法工作
func isMailAble(key string, minTimes int, now time.Time, timeout time.Duration) bool {
	sAble, ok := mailAbleTimes[key]

	// 首次始化记录时间
	if !ok {
		sAble = &mailAble{
			beginTime: now,
			times:     1,
		}
		mailAbleTimes[key] = sAble
	} else {
		// 检查是否在指定的时间内
		d := now.Sub(sAble.beginTime)
		if d < timeout {
			sAble.times++
		} else {
			// 已非超出时间段, 进行重计数
			sAble.beginTime = now
			sAble.times = 1
		}
	}

	// 如果未满足最小发送次数，不执行发送
	if sAble.times < minTimes {
		return false
	}

	// 首次触发时发送
	if sAble.times == minTimes {
		return true
	}

	// 其他时间使用固定次数的算法。
	// 在一定的时间内固定发送次数, 以减少暴量时的发送次数
	return timesToAble(sAble.times)
}
