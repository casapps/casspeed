package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Task represents a scheduled task
type Task struct {
	ID          string
	Name        string
	Schedule    string
	Handler     func(context.Context) error
	Enabled     bool
	LastRun     time.Time
	LastStatus  string
	LastError   string
	NextRun     time.Time
	RunCount    int
	FailCount   int
	RetryOnFail bool
	RetryDelay  time.Duration
	MaxRetries  int
}

// Scheduler manages scheduled tasks
type Scheduler struct {
	cron     *cron.Cron
	tasks    map[string]*Task
	mu       sync.RWMutex
	timezone *time.Location
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a new scheduler
func New(timezone string) (*Scheduler, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cron:     cron.New(cron.WithLocation(loc)),
		tasks:    make(map[string]*Task),
		timezone: loc,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// AddTask adds a task to the scheduler
func (s *Scheduler) AddTask(task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.ID]; exists {
		return fmt.Errorf("task %s already exists", task.ID)
	}

	// Add to cron if enabled
	if task.Enabled {
		_, err := s.cron.AddFunc(task.Schedule, func() {
			s.runTask(task)
		})
		if err != nil {
			return fmt.Errorf("invalid schedule %s: %w", task.Schedule, err)
		}
	}

	s.tasks[task.ID] = task
	return nil
}

// RemoveTask removes a task from the scheduler
func (s *Scheduler) RemoveTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskID]; !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	delete(s.tasks, taskID)
	return nil
}

// GetTask returns a task by ID
func (s *Scheduler) GetTask(taskID string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	return task, nil
}

// GetAllTasks returns all tasks
func (s *Scheduler) GetAllTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// RunTaskNow triggers immediate task execution
func (s *Scheduler) RunTaskNow(taskID string) error {
	s.mu.RLock()
	task, exists := s.tasks[taskID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	go s.runTask(task)
	return nil
}

// runTask executes a task
func (s *Scheduler) runTask(task *Task) {
	if !task.Enabled {
		return
	}

	task.LastRun = time.Now()

	// Execute task handler
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Minute)
	defer cancel()

	err := task.Handler(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	if err != nil {
		task.LastStatus = "failed"
		task.LastError = err.Error()
		task.FailCount++

		// Retry logic
		if task.RetryOnFail && task.FailCount < task.MaxRetries {
			go s.retryTask(task)
		}
	} else {
		task.LastStatus = "success"
		task.LastError = ""
		task.RunCount++
		task.FailCount = 0
	}
}

// retryTask retries a failed task after delay
func (s *Scheduler) retryTask(task *Task) {
	time.Sleep(task.RetryDelay)

	if task.FailCount < task.MaxRetries {
		s.runTask(task)
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop stops the scheduler gracefully
func (s *Scheduler) Stop() {
	s.cancel()
	ctx := s.cron.Stop()
	<-ctx.Done()
}

// EnableTask enables a task
func (s *Scheduler) EnableTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	task.Enabled = true
	return nil
}

// DisableTask disables a task
func (s *Scheduler) DisableTask(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}

	task.Enabled = false
	return nil
}

// GetStatus returns scheduler status
func (s *Scheduler) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	taskStatuses := make([]map[string]interface{}, 0, len(s.tasks))
	for _, task := range s.tasks {
		taskStatuses = append(taskStatuses, map[string]interface{}{
			"id":          task.ID,
			"name":        task.Name,
			"enabled":     task.Enabled,
			"last_run":    task.LastRun,
			"last_status": task.LastStatus,
			"last_error":  task.LastError,
			"next_run":    task.NextRun,
			"run_count":   task.RunCount,
			"fail_count":  task.FailCount,
		})
	}

	return map[string]interface{}{
		"timezone": s.timezone.String(),
		"tasks":    taskStatuses,
	}
}
