package service

import (
	"net/http"

	"github.com/leonelquinteros/gotext"

	"github.com/tnborg/panel/internal/biz"
	"github.com/tnborg/panel/internal/http/request"
)

type SystemdMonitorService struct {
	t                   *gotext.Locale
	systemdMonitorRepo biz.SystemdMonitorRepo
}

func NewSystemdMonitorService(t *gotext.Locale, systemdMonitorRepo biz.SystemdMonitorRepo) *SystemdMonitorService {
	return &SystemdMonitorService{
		t:                   t,
		systemdMonitorRepo: systemdMonitorRepo,
	}
}

// GetConfigs 获取应用的监控配置
func (s *SystemdMonitorService) GetConfigs(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemdMonitorApp](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	configs, err := s.systemdMonitorRepo.GetConfigs(req.AppSlug)
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("获取监控配置失败: %v", err))
		return
	}

	Success(w, configs)
}

// GetAllConfigs 获取所有监控配置
func (s *SystemdMonitorService) GetAllConfigs(w http.ResponseWriter, r *http.Request) {
	configs, err := s.systemdMonitorRepo.GetAllConfigs()
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("获取监控配置失败: %v", err))
		return
	}

	Success(w, configs)
}

// AddConfig 添加监控配置
func (s *SystemdMonitorService) AddConfig(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemdMonitorConfig](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	err = s.systemdMonitorRepo.AddConfig(req.AppSlug, req.Service, req.Enabled, req.Priority)
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("添加监控配置失败: %v", err))
		return
	}

	Success(w, nil)
}

// UpdateConfig 更新监控配置
func (s *SystemdMonitorService) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemdMonitorConfig](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	err = s.systemdMonitorRepo.UpdateConfig(req.ID, req.Enabled, req.Priority)
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("更新监控配置失败: %v", err))
		return
	}

	Success(w, nil)
}

// RemoveConfig 删除监控配置
func (s *SystemdMonitorService) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemdMonitorConfig](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	err = s.systemdMonitorRepo.RemoveConfig(req.ID)
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("删除监控配置失败: %v", err))
		return
	}

	Success(w, nil)
}

// GetMonitorItems 获取应用的监控数据
func (s *SystemdMonitorService) GetMonitorItems(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemdMonitorApp](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	items, err := s.systemdMonitorRepo.GetMonitorItems(req.AppSlug)
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("获取监控数据失败: %v", err))
		return
	}

	Success(w, items)
}

// GetAllMonitorItems 获取所有监控数据
func (s *SystemdMonitorService) GetAllMonitorItems(w http.ResponseWriter, r *http.Request) {
	items, err := s.systemdMonitorRepo.GetAllMonitorItems()
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("获取监控数据失败: %v", err))
		return
	}

	Success(w, items)
}

// CheckServiceStatus 检查服务状态
func (s *SystemdMonitorService) CheckServiceStatus(w http.ResponseWriter, r *http.Request) {
	req, err := Bind[request.SystemdMonitorService](r)
	if err != nil {
		Error(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}

	active, enabled, status, err := s.systemdMonitorRepo.CheckServiceStatus(req.Service)
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("检查服务状态失败: %v", err))
		return
	}

	result := map[string]interface{}{
		"service": req.Service,
		"active":  active,
		"enabled": enabled,
		"status":  status,
	}

	Success(w, result)
}

// CheckAllServices 检查所有配置的服务状态
func (s *SystemdMonitorService) CheckAllServices(w http.ResponseWriter, r *http.Request) {
	err := s.systemdMonitorRepo.CheckAllConfiguredServices()
	if err != nil {
		Error(w, http.StatusInternalServerError, s.t.Get("检查所有服务状态失败: %v", err))
		return
	}

	Success(w, nil)
}