package repository

import (
	"github.com/aneesh1213/tasker/internal/server"
)

type Repositories struct {
	Todo *TodoRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Todo: NewTodoRepository(s),
	}
}
