package task

import (
	"greatestworks/internal/note/event"
)

type ITask interface {
	SetStatus(Status)
	OnEvent(event event.IEvent)
	GetTaskData() ITaskData
}

// ITaskData task data  save task data eg:target data or child task data
type ITaskData interface {
	GetProgress()
	GetTotalProgress()
	CheckComplete()
}
