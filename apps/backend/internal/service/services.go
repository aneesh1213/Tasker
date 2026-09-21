package service

import (
	"github.com/aneesh1213/tasker/internal/lib/job"
	"github.com/aneesh1213/tasker/internal/repository"
	"github.com/aneesh1213/tasker/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}
