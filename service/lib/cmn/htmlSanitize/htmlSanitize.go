package htmlSanitize

import (
	"sync"

	"github.com/microcosm-cc/bluemonday"
)

var (
	policyOnce sync.Once
	policy     *bluemonday.Policy
)

// GetPolicy 内容消毒策略（单例）
//
// 自定义代码模块是 UGC 内容（开发者投稿），直接 v-html 渲染存在 XSS 风险，
// 因此保存时必须消毒。白名单需要保留 <pre><code class="language-*">，
// 否则前台的代码高亮与复制功能会失效。
func GetPolicy() *bluemonday.Policy {
	policyOnce.Do(func() {
		p := bluemonday.UGCPolicy()

		// 代码块与排版所需元素
		p.AllowElements("pre", "code", "figure", "figcaption", "hr", "section", "article")

		// class 属性：代码块的语言标识 language-xxx 依赖它
		p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).Globally()

		// 保留富文本编辑器生成的常用内联样式
		p.AllowStyles(
			"text-align", "text-decoration", "vertical-align", "white-space",
			"color", "background", "background-color",
			"font-size", "font-weight", "font-style", "font-family", "line-height",
			"margin", "margin-top", "margin-right", "margin-bottom", "margin-left",
			"padding", "padding-top", "padding-right", "padding-bottom", "padding-left",
			"width", "height", "max-width", "display", "float",
			"border", "border-collapse", "border-color", "border-width", "border-style",
		).Globally()

		// 图片
		p.AllowAttrs("loading", "decoding").OnElements("img")
		p.AllowAttrs("width", "height").Matching(bluemonday.Number).OnElements("img", "td", "th")

		// 链接
		p.AllowAttrs("target").Matching(bluemonday.Paragraph).OnElements("a")
		p.RequireNoFollowOnFullyQualifiedLinks(true)
		p.AddTargetBlankToFullyQualifiedLinks(true)

		// 只放行安全协议（自动排除 javascript: 等）
		p.AllowURLSchemes("http", "https", "mailto")

		policy = p
	})

	return policy
}

// Sanitize 对富文本内容做消毒
func Sanitize(html string) string {
	return GetPolicy().Sanitize(html)
}
