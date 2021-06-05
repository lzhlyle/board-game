package core

import "fmt"

// Dashboard 仪表盘
type Dashboard interface {
	Title() string
	Refresh() string
}

type DefaultDashboard struct {
	title   string
	content string
}

func NewDefaultDashboard(title string) *DefaultDashboard {
	return &DefaultDashboard{title: title}
}

func (d *DefaultDashboard) Title() string {
	return d.title
}

func (d *DefaultDashboard) Refresh() string {
	return d.content
}

func (d *DefaultDashboard) Clear() {
	d.content = ""
}

func (d *DefaultDashboard) Appendf(format string, arg ...interface{}) *DefaultDashboard {
	d.content += fmt.Sprintf(format, arg...)
	return d
}

func (d *DefaultDashboard) LineEnd() {
	d.content += "\n"
}
