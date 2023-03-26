package main

import "errors"

var (
	NoEndpoint = errors.New("do not exit available endpoint ")
	NoZoneId   = errors.New("no zone list available")
)
