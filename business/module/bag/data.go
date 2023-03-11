package bag

import (
	"greatestworks/business/module/hub"
)

// Data 对应DB -> mongo
type Data struct {
	hub.DataAsPublisher
}
