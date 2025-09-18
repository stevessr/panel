package data

import (
	"log/slog"
	"time"

	"github.com/leonelquinteros/gotext"
	"gorm.io/gorm"

	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/pkg/systemctl"
)

type systemdMonitorRepo struct {
	db  *gorm.DB
	log *slog.Logger
	t   *gotext.Locale
}

func NewSystemdMonitorRepo(db *gorm.DB, log *slog.Logger, t *gotext.Locale) biz.SystemdMonitorRepo {
	return &systemdMonitorRepo{
		db:  db,
		log: log,
		t:   t,
	}
}

// GetConfigs 获取应用的监控配置
func (r *systemdMonitorRepo) GetConfigs(appSlug string) ([]*biz.SystemdMonitorConfig, error) {
	var configs []*biz.SystemdMonitorConfig
	err := r.db.Where("app_slug = ?", appSlug).Order("priority DESC, service ASC").Find(&configs).Error
	return configs, err
}

// GetAllConfigs 获取所有监控配置
func (r *systemdMonitorRepo) GetAllConfigs() ([]*biz.SystemdMonitorConfig, error) {
	var configs []*biz.SystemdMonitorConfig
	err := r.db.Where("enabled = ?", true).Order("app_slug ASC, priority DESC, service ASC").Find(&configs).Error
	return configs, err
}

// AddConfig 添加监控配置
func (r *systemdMonitorRepo) AddConfig(appSlug, service string, enabled bool, priority int) error {
	config := &biz.SystemdMonitorConfig{
		AppSlug:  appSlug,
		Service:  service,
		Enabled:  enabled,
		Priority: priority,
	}
	return r.db.Create(config).Error
}

// UpdateConfig 更新监控配置
func (r *systemdMonitorRepo) UpdateConfig(id uint, enabled bool, priority int) error {
	return r.db.Model(&biz.SystemdMonitorConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
		"enabled":  enabled,
		"priority": priority,
	}).Error
}

// RemoveConfig 删除监控配置
func (r *systemdMonitorRepo) RemoveConfig(id uint) error {
	return r.db.Delete(&biz.SystemdMonitorConfig{}, id).Error
}

// GetMonitorItems 获取应用的监控数据
func (r *systemdMonitorRepo) GetMonitorItems(appSlug string) ([]*biz.SystemdMonitorItem, error) {
	var items []*biz.SystemdMonitorItem
	err := r.db.Where("app_slug = ?", appSlug).Order("checked_at DESC").Find(&items).Error
	return items, err
}

// GetAllMonitorItems 获取所有监控数据
func (r *systemdMonitorRepo) GetAllMonitorItems() ([]*biz.SystemdMonitorItem, error) {
	var items []*biz.SystemdMonitorItem
	err := r.db.Order("app_slug ASC, checked_at DESC").Find(&items).Error
	return items, err
}

// UpdateMonitorStatus 更新监控状态
func (r *systemdMonitorRepo) UpdateMonitorStatus(appSlug, service string, active, enabled bool, status string) error {
	now := time.Now()
	
	// 首先尝试更新现有记录
	result := r.db.Model(&biz.SystemdMonitorItem{}).
		Where("app_slug = ? AND service = ?", appSlug, service).
		Updates(map[string]interface{}{
			"active":     active,
			"enabled":    enabled,
			"status":     status,
			"checked_at": now,
		})
	
	if result.Error != nil {
		return result.Error
	}
	
	// 如果没有记录被更新，则创建新记录
	if result.RowsAffected == 0 {
		item := &biz.SystemdMonitorItem{
			AppSlug:   appSlug,
			Service:   service,
			Active:    active,
			Enabled:   enabled,
			Status:    status,
			CheckedAt: now,
		}
		return r.db.Create(item).Error
	}
	
	return nil
}

// CleanOldMonitorData 清理过期的监控数据
func (r *systemdMonitorRepo) CleanOldMonitorData(days int) error {
	if days <= 0 {
		return nil
	}
	
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("checked_at < ?", cutoff).Delete(&biz.SystemdMonitorItem{}).Error
}

// CheckServiceStatus 检查服务状态
func (r *systemdMonitorRepo) CheckServiceStatus(service string) (active, enabled bool, status string, err error) {
	active, err = systemctl.Status(service)
	if err != nil {
		r.log.Warn("检查服务状态失败", slog.String("service", service), slog.Any("err", err))
		return false, false, "检查失败", err
	}
	
	enabled, err = systemctl.IsEnabled(service)
	if err != nil {
		r.log.Warn("检查服务启用状态失败", slog.String("service", service), slog.Any("err", err))
		return active, false, "检查失败", err
	}
	
	if active {
		status = "运行中"
	} else {
		status = "已停止"
	}
	
	if !enabled {
		status += "(未启用)"
	}
	
	return active, enabled, status, nil
}

// CheckAllConfiguredServices 检查所有配置的服务状态
func (r *systemdMonitorRepo) CheckAllConfiguredServices() error {
	configs, err := r.GetAllConfigs()
	if err != nil {
		return err
	}
	
	for _, config := range configs {
		active, enabled, status, err := r.CheckServiceStatus(config.Service)
		if err != nil {
			r.log.Warn("检查服务状态失败", 
				slog.String("app", config.AppSlug),
				slog.String("service", config.Service),
				slog.Any("err", err))
			continue
		}
		
		err = r.UpdateMonitorStatus(config.AppSlug, config.Service, active, enabled, status)
		if err != nil {
			r.log.Warn("更新监控状态失败",
				slog.String("app", config.AppSlug),
				slog.String("service", config.Service),
				slog.Any("err", err))
		}
	}
	
	return nil
}