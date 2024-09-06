package timewheel

import "time"

/*
   @Author: orbit-w
   @File: const
   @2024 8月 周一 00:14
*/

const (
	taskTimeout = time.Second * 5
)

const (
	StateNormal = iota
	StateClosed
)

const (
	LvSecond = iota
	LvMinute
	LvHour
)

const (
	HundredMsName     = "HUNDRED_MS"
	HundredMsInterval = 100
	HundredMsScales   = 10

	//HourName 小时
	HourName = "HOUR"
	//HourInterval 小时间隔ms为精度
	HourInterval = 60 * 60 * 1e3
	//HourScales  12小时制
	HourScales = 12

	//MinuteName 分钟
	MinuteName = "MINUTE"
	//MinuteInterval 每分钟时间间隔
	MinuteInterval = 60 * 1e3
	//MinuteScales 60分钟
	MinuteScales = 60

	//SecondName  秒
	SecondName = "SECOND"
	//SecondInterval 秒的间隔
	SecondInterval = 1e3
	//SecondScales  60秒
	SecondScales = 60
)
