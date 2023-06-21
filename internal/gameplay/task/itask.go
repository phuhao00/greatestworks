package task

import (
	"greatestworks/internal/note/event"
)

type ITask interface {
	SetStatus(Status)
	OnEvent(event event.IEvent)
	GetTaskData() ITaskData
}

// ITaskData task_category_group data  save task_category_group data eg:target data or child task_category_group data
type ITaskData interface {
	GetProgress()
	GetTotalProgress()
	CheckComplete()
}
