package student

import "context"

// Student представляет студента
type Student struct {
	ID         string  `json:"id,omitempty"`
	FullName   string  `json:"full_name"`
	Group      string  `json:"group"`
	GPA        float32 `json:"gpa"`
	EnrolledOn int64   `json:"enrolled_on,omitempty"`
}

// Service описывает бизнес-логику сервиса студентов
type Service interface {
	Create(ctx context.Context, student Student) (string, error)
	GetByID(ctx context.Context, id string) (Student, error)
	ChangeGroup(ctx context.Context, id string, group string) error
}

// Repository описывает хранилище студентов
type Repository interface {
	CreateStudent(ctx context.Context, student Student) error
	GetStudentByID(ctx context.Context, id string) (Student, error)
	ChangeStudentGroup(ctx context.Context, id string, group string) error
}
