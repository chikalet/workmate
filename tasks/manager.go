package tasks

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type Manager struct {
	mu     sync.RWMutex
	tasks  map[int]*Task
	nextID int
}

func NewManager() *Manager {
	return &Manager{
		tasks: make(map[int]*Task),
	}
}

func (m *Manager) CreateTask() int {
	m.mu.Lock()
	id := m.nextID
	task := &Task{
		ID:        id,
		Status:    StatusCreated,
		CreatedAt: time.Now(),
	}
	m.tasks[id] = task
	m.nextID++
	m.mu.Unlock()

	go m.processTask(task)
	return id
}

func (m *Manager) processTask(task *Task) {
	ctx, cancel := context.WithCancel(context.Background())
	task.cancel = cancel

	select {
	case <-time.After(1 * time.Second):
		m.mu.Lock()
		if t, ok := m.tasks[task.ID]; ok {
			t.Status = StatusInProgress
			t.startedAt = time.Now()
		}
		m.mu.Unlock()
	case <-ctx.Done():
		return
	}

	duration := time.Duration(180+rand.Intn(121)) * time.Second

	select {
	case <-time.After(duration):
		m.mu.Lock()
		if t, ok := m.tasks[task.ID]; ok {
			t.Status = StatusCompleted
			t.DurationSeconds = int64(time.Since(t.startedAt) / time.Second)
		}
		m.mu.Unlock()
	case <-ctx.Done():
		return
	}
}

func (m *Manager) GetTask(id int) (*Task, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if task, ok := m.tasks[id]; ok {
		copyTask := *task
		return &copyTask, true
	}
	return nil, false
}

func (m *Manager) DeleteTask(id int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	task, ok := m.tasks[id]
	if !ok {
		return false
	}
	if task.cancel != nil {
		task.cancel()
	}
	delete(m.tasks, id)
	return true
}
