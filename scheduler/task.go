package scheduler

type Task struct {
	ID       string
	Name     string
	Schedule string
	JobFunc  func()
}
