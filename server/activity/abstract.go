package main

type IActivity interface {
	OnDayReset()
	OnInit()
}
