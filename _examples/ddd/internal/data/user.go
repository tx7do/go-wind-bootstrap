// Package data 实现领域层定义的仓储接口。
//
// data 层负责所有与基础设施（数据库、缓存、消息队列）的交互。
// 它依赖 domain 层的接口定义，但不被 domain 层感知。
package data

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/tx7do/go-wind-bootstrap/_examples/ddd/internal/domain"
)

// 确保 userRepo 实现了 domain.UserRepository 接口。
var _ domain.UserRepository = (*userRepo)(nil)

// userPO 是 GORM 持久化对象（数据库表结构）。
type userPO struct {
	ID        uint64         `gorm:"primaryKey"`
	Name      string         `gorm:"size:128;not null"`
	Email     string         `gorm:"size:256;uniqueIndex;not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定 GORM 表名。
func (userPO) TableName() string { return "users" }

// NewUserRepository 创建基于 GORM 的用户仓储。
//
// db 参数来自 bootstrap 的 Database 适配层：
//
//	rawDB := ctx.Database(bootstrap.DatabaseTypeGorm).(*gormCrud.Client)
//	db := rawDB.DB
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	// 自动迁移表结构
	_ = db.AutoMigrate(&userPO{})

	return &userRepo{db: db}
}

type userRepo struct {
	db *gorm.DB
}

func (r *userRepo) Create(req *domain.CreateUserRequest) (*domain.User, error) {
	po := &userPO{
		Name:  req.Name,
		Email: req.Email,
	}
	if err := r.db.Create(po).Error; err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return poToEntity(po), nil
}

func (r *userRepo) GetByID(id uint64) (*domain.User, error) {
	var po userPO
	if err := r.db.First(&po, id).Error; err != nil {
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}
	return poToEntity(&po), nil
}

func (r *userRepo) List(offset, limit int) ([]*domain.User, error) {
	var pos []userPO
	if err := r.db.Offset(offset).Limit(limit).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	users := make([]*domain.User, len(pos))
	for i, po := range pos {
		users[i] = poToEntity(&po)
	}
	return users, nil
}

func (r *userRepo) Update(id uint64, req *domain.UpdateUserRequest) (*domain.User, error) {
	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if len(updates) == 0 {
		return r.GetByID(id)
	}
	if err := r.db.Model(&userPO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update user %d: %w", id, err)
	}
	return r.GetByID(id)
}

func (r *userRepo) Delete(id uint64) error {
	if err := r.db.Delete(&userPO{}, id).Error; err != nil {
		return fmt.Errorf("delete user %d: %w", id, err)
	}
	return nil
}

// poToEntity 将持久化对象转换为领域实体。
func poToEntity(po *userPO) *domain.User {
	return &domain.User{
		ID:        po.ID,
		Name:      po.Name,
		Email:     po.Email,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}
