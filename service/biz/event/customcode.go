package event

// 自定义代码领域事件
const (
	// CustomCodeSubmitted 作者提交完成（新增/编辑后触发）
	//   Auto=true  : 受信任作者且安全扫描通过，已在请求内同步上线
	//   Auto=false : 进入待审核，需由异步 handler 处理通知审核员/审计等副作用
	CustomCodeSubmitted = "custom_code.submitted"
)

// CustomCodeSubmittedPayload 自定义代码提交事件载荷
type CustomCodeSubmittedPayload struct {
	CustomCodeID uint
	ReviewID     uint
	AuthorID     uint
	Auto         bool
}
