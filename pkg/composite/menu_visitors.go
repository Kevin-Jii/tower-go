package composite

import (
	"strings"
)

// PermissionCollector 权限收集访问者
type PermissionCollector struct {
	Permissions []string
}

func NewPermissionCollector() *PermissionCollector {
	return &PermissionCollector{
		Permissions: make([]string, 0),
	}
}

func (v *PermissionCollector) VisitDirectory(node *MenuNode) {
	if node.menu.Permission != "" {
		v.Permissions = append(v.Permissions, node.menu.Permission)
	}
}

func (v *PermissionCollector) VisitMenu(node *MenuNode) {
	if node.menu.Permission != "" {
		v.Permissions = append(v.Permissions, node.menu.Permission)
	}
}

func (v *PermissionCollector) VisitButton(node *MenuNode) {
	if node.menu.Permission != "" {
		v.Permissions = append(v.Permissions, node.menu.Permission)
	}
}

// PathCollector 路径收集访问者
type PathCollector struct {
	Paths []string
}

func NewPathCollector() *PathCollector {
	return &PathCollector{
		Paths: make([]string, 0),
	}
}

func (v *PathCollector) VisitDirectory(node *MenuNode) {
	if node.menu.Path != "" {
		v.Paths = append(v.Paths, node.menu.Path)
	}
}

func (v *PathCollector) VisitMenu(node *MenuNode) {
	if node.menu.Path != "" {
		v.Paths = append(v.Paths, node.menu.Path)
	}
}

func (v *PathCollector) VisitButton(node *MenuNode) {
	// 按钮通常没有路径
}

// MenuCounter 菜单统计访问者
type MenuCounter struct {
	Directories int
	Menus       int
	Buttons     int
}

func NewMenuCounter() *MenuCounter {
	return &MenuCounter{}
}

func (v *MenuCounter) VisitDirectory(node *MenuNode) {
	v.Directories++
}

func (v *MenuCounter) VisitMenu(node *MenuNode) {
	v.Menus++
}

func (v *MenuCounter) VisitButton(node *MenuNode) {
	v.Buttons++
}

func (v *MenuCounter) Total() int {
	return v.Directories + v.Menus + v.Buttons
}

// MenuPrinter 菜单打印访问者（用于调试）
type MenuPrinter struct {
	Output strings.Builder
	indent int
}

func NewMenuPrinter() *MenuPrinter {
	return &MenuPrinter{}
}

func (v *MenuPrinter) VisitDirectory(node *MenuNode) {
	v.printNode(node, "📁")
}

func (v *MenuPrinter) VisitMenu(node *MenuNode) {
	v.printNode(node, "📄")
}

func (v *MenuPrinter) VisitButton(node *MenuNode) {
	v.printNode(node, "🔘")
}

func (v *MenuPrinter) printNode(node *MenuNode, icon string) {
	indent := strings.Repeat("  ", v.indent)
	v.Output.WriteString(indent)
	v.Output.WriteString(icon)
	v.Output.WriteString(" ")
	v.Output.WriteString(node.menu.Title)
	if node.menu.Path != "" {
		v.Output.WriteString(" (")
		v.Output.WriteString(node.menu.Path)
		v.Output.WriteString(")")
	}
	v.Output.WriteString("\n")
}

func (v *MenuPrinter) String() string {
	return v.Output.String()
}

// DepthTrackingVisitor 深度追踪访问者包装器
type DepthTrackingVisitor struct {
	inner   MenuVisitor
	depth   int
	OnEnter func(depth int)
	OnLeave func(depth int)
}

func NewDepthTrackingVisitor(inner MenuVisitor) *DepthTrackingVisitor {
	return &DepthTrackingVisitor{
		inner: inner,
	}
}

func (v *DepthTrackingVisitor) VisitDirectory(node *MenuNode) {
	if v.OnEnter != nil {
		v.OnEnter(v.depth)
	}
	v.inner.VisitDirectory(node)
	v.depth++
}

func (v *DepthTrackingVisitor) VisitMenu(node *MenuNode) {
	if v.OnEnter != nil {
		v.OnEnter(v.depth)
	}
	v.inner.VisitMenu(node)
	v.depth++
}

func (v *DepthTrackingVisitor) VisitButton(node *MenuNode) {
	if v.OnEnter != nil {
		v.OnEnter(v.depth)
	}
	v.inner.VisitButton(node)
}

// FilterVisitor 过滤访问者
type FilterVisitor struct {
	inner     MenuVisitor
	predicate func(*MenuNode) bool
}

func NewFilterVisitor(inner MenuVisitor, predicate func(*MenuNode) bool) *FilterVisitor {
	return &FilterVisitor{
		inner:     inner,
		predicate: predicate,
	}
}

func (v *FilterVisitor) VisitDirectory(node *MenuNode) {
	if v.predicate(node) {
		v.inner.VisitDirectory(node)
	}
}

func (v *FilterVisitor) VisitMenu(node *MenuNode) {
	if v.predicate(node) {
		v.inner.VisitMenu(node)
	}
}

func (v *FilterVisitor) VisitButton(node *MenuNode) {
	if v.predicate(node) {
		v.inner.VisitButton(node)
	}
}
