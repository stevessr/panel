package job

import (
	"log/slog"

	"github.com/tnborg/panel/internal/app"
	"github.com/tnborg/panel/internal/biz"
)

// SystemdMonitoring systemd 服务监控
type SystemdMonitoring struct {
	log                 *slog.Logger
	systemdMonitorRepo biz.SystemdMonitorRepo
	settingRepo        biz.SettingRepo
}

func NewSystemdMonitoring(log *slog.Logger, systemdMonitorRepo biz.SystemdMonitorRepo, settingRepo biz.SettingRepo) *SystemdMonitoring {
	return &SystemdMonitoring{
		log:                 log,
		systemdMonitorRepo: systemdMonitorRepo,
		settingRepo:        settingRepo,
	}
}

func (r *SystemdMonitoring) Run() {
	if app.Status != app.StatusNormal {
		return
	}

	// 检查是否启用了 systemd 监控
	systemdMonitorEnabled, err := r.settingRepo.Get("systemd_monitor")
	if err != nil || systemdMonitorEnabled != "1" {
		// 如果设置不存在，默认启用
		if err != nil {
			r.log.Debug("[SystemdMonitor] systemd监控设置不存在，使用默认启用")
		} else {
			r.log.Debug("[SystemdMonitor] systemd监控已禁用")
			return
		}
	}

	// 检查所有配置的服务状态
	err = r.systemdMonitorRepo.CheckAllConfiguredServices()
	if err != nil {
		r.log.Warn("[SystemdMonitor] 检查服务状态失败", slog.Any("err", err))
		return
	}

	// 清理过期数据（保留7天的监控数据）
	err = r.systemdMonitorRepo.CleanOldMonitorData(7)
	if err != nil {
		r.log.Warn("[SystemdMonitor] 清理过期监控数据失败", slog.Any("err", err))
	}

	r.log.Debug("[SystemdMonitor] systemd服务监控检查完成")
}