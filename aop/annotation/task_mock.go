package annotation

type TaskMock struct {
	Category    int
	TargetParam []string //monsterId,level,count
}

func (TaskMock) OnChange() {

}

type event struct {
	Category                int
	monsterId, level, count int
}

func OnEvent(e event) {
	for _, mock := range taskInstances {
		if mock.Category == e.Category {
			//tmp1, tmp2, tmp3 := mock[0], mock[1], mock[2]
		}
	}
}

var (
	taskInstances map[int]TaskMock
)
