package repository

import (
	"context"
	"sync"

	"example.com/student-service/student"
)

type memoryRepo struct {
	mu       sync.RWMutex
	students map[string]student.Student
}

// NewMemoryRepository создаёт in-memory хранилище студентов
func NewMemoryRepository() student.Repository {
	return &memoryRepo{
		students: make(map[string]student.Student),
	}
}

func (r *memoryRepo) CreateStudent(ctx context.Context, st student.Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.students[st.ID] = st
	return nil
}

func (r *memoryRepo) GetStudentByID(ctx context.Context, id string) (student.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	st, ok := r.students[id]
	if !ok {
		return student.Student{}, student.ErrStudentNotFound
	}
	return st, nil
}

func (r *memoryRepo) ChangeStudentGroup(ctx context.Context, id string, group string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	st, ok := r.students[id]
	if !ok {
		return student.ErrStudentNotFound
	}
	st.Group = group
	r.students[id] = st
	return nil
}
