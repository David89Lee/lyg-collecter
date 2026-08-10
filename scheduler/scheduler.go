package scheduler

import "sync"
import "github.com/robfig/cron/v3"

type Scheduler struct {
	c     *cron.Cron
	mu    sync.Mutex
	tasks map[string]cron.EntryID
	funcs map[string]*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		c:     cron.New(),
		tasks: make(map[string]cron.EntryID),
		funcs: make(map[string]*Task),
	}
}

func (s *Scheduler) Start() {
	s.c.Start()
}
func (s *Scheduler) Stop() {
	s.c.Stop()
}
func (s *Scheduler) AddTask(t *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, err := s.c.AddFunc(t.Schedule, t.JobFunc)
	if err != nil {
		return err
	}
	s.tasks[t.ID] = id
	s.funcs[t.ID] = t
	return nil
}

func (s *Scheduler) UpdateTask(newTask *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.funcs[newTask.ID]; ok {
		s.c.Remove(s.tasks[newTask.ID])
		delete(s.tasks, newTask.ID)
		delete(s.funcs, newTask.ID)
		id, err := s.c.AddFunc(newTask.Schedule, newTask.JobFunc)
		if err != nil {
			return err
		}
		s.tasks[newTask.ID] = id
		s.funcs[newTask.ID] = newTask
		return nil
	}
	return s.AddTask(newTask)
}

func (s *Scheduler) RemoveTask(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id, ok := s.tasks[taskID]; ok {
		s.c.Remove(id)
		delete(s.tasks, taskID)
		delete(s.funcs, taskID)
	}
}
