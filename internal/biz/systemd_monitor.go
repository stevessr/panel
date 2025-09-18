package biz

import (
	"time"
)

// SystemdMonitorItem systemd 服务监控项
type SystemdMonitorItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AppSlug   string    `gorm:"not null;default:'';index" json:"app_slug"`           // 应用标识
	Service   string    `gorm:"not null;default:''" json:"service"`                 // 服务名称
	Active    bool      `gorm:"not null;default:false" json:"active"`               // 是否活跃
	Enabled   bool      `gorm:"not null;default:false" json:"enabled"`              // 是否启用
	Status    string    `gorm:"not null;default:''" json:"status"`                  // 状态描述
	CheckedAt time.Time `gorm:"not null" json:"checked_at"`                         // 最后检查时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SystemdMonitorConfig systemd 监控配置
type SystemdMonitorConfig struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	AppSlug  string `gorm:"not null;default:'';index" json:"app_slug"`  // 应用标识
	Service  string `gorm:"not null;default:''" json:"service"`         // 服务名称
	Enabled  bool   `gorm:"not null;default:true" json:"enabled"`       // 是否启用监控
	Priority int    `gorm:"not null;default:0" json:"priority"`         // 优先级
}

// SystemdMonitorRepo systemd 监控仓库接口
type SystemdMonitorRepo interface {
	// 配置管理
	GetConfigs(appSlug string) ([]*SystemdMonitorConfig, error)
	GetAllConfigs() ([]*SystemdMonitorConfig, error)
	AddConfig(appSlug, service string, enabled bool, priority int) error
	UpdateConfig(id uint, enabled bool, priority int) error
	RemoveConfig(id uint) error

	// 监控数据管理
	GetMonitorItems(appSlug string) ([]*SystemdMonitorItem, error)
	GetAllMonitorItems() ([]*SystemdMonitorItem, error)
	UpdateMonitorStatus(appSlug, service string, active, enabled bool, status string) error
	CleanOldMonitorData(days int) error

	// 服务状态检查
	CheckServiceStatus(service string) (active, enabled bool, status string, err error)
	CheckAllConfiguredServices() error
}