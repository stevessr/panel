package request

// SystemdMonitorConfig systemd 监控配置请求
type SystemdMonitorConfig struct {
	ID       uint   `json:"id" form:"id"`
	AppSlug  string `json:"app_slug" form:"app_slug" validate:"required"`
	Service  string `json:"service" form:"service" validate:"required"`
	Enabled  bool   `json:"enabled" form:"enabled"`
	Priority int    `json:"priority" form:"priority"`
}

// SystemdMonitorService 服务名称请求
type SystemdMonitorService struct {
	Service string `json:"service" form:"service" validate:"required"`
}

// SystemdMonitorApp 应用标识请求
type SystemdMonitorApp struct {
	AppSlug string `json:"app_slug" form:"app_slug" validate:"required"`
}