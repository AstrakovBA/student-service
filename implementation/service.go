package implementation

import (
	"context"
	"time"

	"github.com/go-kit/log"       // <-- ИСПРАВЛЕНО
	"github.com/go-kit/log/level" // <-- ИСПРАВЛЕНО
	"github.com/google/uuid"

	"example.com/student-service/student"
)

type service struct {
	repository student.Repository
	logger     log.Logger
}

// NewService создаёт новый экземпляр сервиса студентов
func NewService(repo student.Repository, logger log.Logger) student.Service {
	return &service{
		repository: repo,
		logger:     logger,
	}
}

func (s *service) Create(ctx context.Context, st student.Student) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	id := uuid.New().String()
	st.ID = id
	st.EnrolledOn = time.Now().Unix()

	if err := s.repository.CreateStudent(ctx, st); err != nil {
		_ = level.Error(logger).Log("err", err) // <-- ИСПРАВЛЕНО: добавлено _ =
		return "", student.ErrCmdRepository
	}

	return id, nil
}

func (s *service) GetByID(ctx context.Context, id string) (student.Student, error) {
	logger := log.With(s.logger, "method", "GetByID")

	st, err := s.repository.GetStudentByID(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err) // <-- ИСПРАВЛЕНО: добавлено _ =
		return st, student.ErrStudentNotFound
	}

	return st, nil
}

func (s *service) ChangeGroup(ctx context.Context, id string, group string) error {
	logger := log.With(s.logger, "method", "ChangeGroup")

	if err := s.repository.ChangeStudentGroup(ctx, id, group); err != nil {
		_ = level.Error(logger).Log("err", err) // <-- ИСПРАВЛЕНО: добавлено _ =
		return student.ErrCmdRepository
	}

	return nil
}
