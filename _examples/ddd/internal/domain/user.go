// Package domain 定义业务领域模型和仓储接口。
//
// 领域层是 DDD 分层的核心，不依赖任何外部框架或基础设施。
// data 层实现此包中定义的接口，service 层通过接口调用。
package domain

import "time"

// User 是用户聚合根。
type User struct {
	ID        uint64    // 主键
	Name      string    // 用户名
	Email     string    // 邮箱
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间
}

// CreateUserRequest 创建用户请求。
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest 更新用户请求。
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// UserRepository 定义用户持久化接口。
// data 层负责实现此接口。
type UserRepository interface {
	// Create 创建用户并返回 ID。
	Create(req *CreateUserRequest) (*User, error)
	// GetByID 按 ID 查询用户。
	GetByID(id uint64) (*User, error)
	// List 分页列出用户。
	List(offset, limit int) ([]*User, error)
	// Update 更新用户。
	Update(id uint64, req *UpdateUserRequest) (*User, error)
	// Delete 删除用户。
	Delete(id uint64) error
}
