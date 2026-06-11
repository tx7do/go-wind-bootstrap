// Package service 实现业务逻辑层。
//
// service 层编排 domain 层的业务规则，调用 data 层的仓储接口完成用例。
// 它不直接依赖任何基础设施（数据库、缓存、HTTP），只依赖 domain 接口。
package service

import (
	"fmt"

	"github.com/tx7do/go-wind-bootstrap/_examples/ddd/internal/domain"
)

// UserService 用户业务逻辑。
type UserService struct {
	repo domain.UserRepository
}

// NewUserService 创建用户服务。
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser 创建用户。
func (s *UserService) CreateUser(req *domain.CreateUserRequest) (*domain.User, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	return s.repo.Create(req)
}

// GetUser 获取用户。
func (s *UserService) GetUser(id uint64) (*domain.User, error) {
	return s.repo.GetByID(id)
}

// ListUsers 列出用户。
func (s *UserService) ListUsers(offset, limit int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return s.repo.List(offset, limit)
}

// UpdateUser 更新用户。
func (s *UserService) UpdateUser(id uint64, req *domain.UpdateUserRequest) (*domain.User, error) {
	return s.repo.Update(id, req)
}

// DeleteUser 删除用户。
func (s *UserService) DeleteUser(id uint64) error {
	return s.repo.Delete(id)
}
