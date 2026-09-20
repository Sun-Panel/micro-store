package admin

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"sun-panel/models"

	"gorm.io/gorm"
)

// 作者信任窗口：最近一次审核通过距今 ≤ 此值，其提交可走自动审核
const AuthorTrustWindow = 365 * 24 * time.Hour

// 提交后的审核模式
type ReviewMode int

const (
	ReviewModeManual ReviewMode = iota // 人工审核
	ReviewModeAuto                     // 自动审核
)

// DecideReviewMode 审核模式决策（策略集中点）。
// 作者受信任(trusted) 且 安全扫描通过(scanPass) -> 自动审核，否则人工。
// 以后调整策略（如白名单、特定类型免扫等）只改这里。
func DecideReviewMode(trusted, scanPass bool) ReviewMode {
	if trusted && scanPass {
		return ReviewModeAuto
	}
	return ReviewModeManual
}

// AuthorLatestApprovedTime 取作者名下最近一次「审核通过」的时间（跨其所有帖子）。
// 无通过记录返回零值。
func AuthorLatestApprovedTime(db *gorm.DB, authorId uint) (time.Time, error) {
	var rt time.Time
	err := db.Model(&models.CustomCodeReview{}).
		Joins("INNER JOIN custom_code ON custom_code.id = custom_code_review.custom_code_id").
		Where("custom_code.author_id = ? AND custom_code_review.status = ?", authorId, models.CustomCodeReviewStatusApproved).
		Select("MAX(custom_code_review.review_time)").
		Scan(&rt).Error
	return rt, err
}

// ShouldAutoApprove 作者近一年内是否审核通过过（受信任）。
// DB 出错时保守返回 false（走人工）。
func ShouldAutoApprove(db *gorm.DB, authorId uint) (bool, error) {
	latest, err := AuthorLatestApprovedTime(db, authorId)
	if err != nil {
		return false, err
	}
	if latest.IsZero() {
		return false, nil
	}
	return time.Since(latest) <= AuthorTrustWindow, nil
}

// PublishApprovedReview 将审核快照直接发布上线（人工 Approve 与自动审核共用同一逻辑）。
// reviewerId 传 0 表示系统自动审核；note 为审核说明。
func PublishApprovedReview(db *gorm.DB, mainID, reviewID uint, reviewerId uint, note string) error {
	mReview := &models.CustomCodeReview{}
	review, err := mReview.GetById(db, reviewID)
	if err != nil {
		return err
	}

	mBlock := &models.CustomCodeBlock{}
	snapshotBlocks, err := mBlock.GetListByCustomCodeId(db, mainID, reviewID)
	if err != nil {
		return err
	}

	now := time.Now()
	err = db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"title":             review.Title,
			"description":       review.Description,
			"keywords":          review.Keywords,
			"is_original":       review.IsOriginal,
			"source_note":       review.SourceNote,
			"versions":          review.Versions,
			"code_types":        review.CodeTypes,
			"status":            models.CustomCodeStatusOnline,
			"published_at":      now,
			"current_review_id": review.ID,
			"offline_reason":    "",
		}
		if err := tx.Model(&models.CustomCode{}).Where("id = ?", mainID).Updates(updates).Error; err != nil {
			return err
		}

		onlineBlocks := make([]models.CustomCodeBlock, 0, len(snapshotBlocks))
		for _, block := range snapshotBlocks {
			onlineBlocks = append(onlineBlocks, models.CustomCodeBlock{
				Sort:     block.Sort,
				CodeType: block.CodeType,
				Title:    block.Title,
				Note:     block.Note,
				Code:     block.Code,
				Images:   block.Images,
				OnlyId:   block.OnlyId,
			})
		}

		return (&models.CustomCodeBlock{}).ReplaceAll(tx, mainID, 0, onlineBlocks)
	})
	if err != nil {
		return err
	}

	return mReview.UpdateStatus(db, reviewID, models.CustomCodeReviewStatusApproved, reviewerId, note)
}

// ---------------- 自动化安全扫描门禁 ----------------

// AutoReviewSecurityScan 提交自动审核前的静态安全扫描（不执行代码）。
// 命中高危模式返回 pass=false 及原因；由调用方降级为人工审核。
func AutoReviewSecurityScan(blocks []models.CustomCodeBlock) (pass bool, reasons []string) {
	for _, b := range blocks {
		switch b.CodeType {
		case models.CodeTypeJS:
			if hits := scanJS(b.Code); len(hits) > 0 {
				reasons = append(reasons, fmt.Sprintf("JS块[%s]: %s", b.OnlyId, strings.Join(hits, "; ")))
			}
		case models.CodeTypeCSS:
			if hits := scanCSS(b.Code); len(hits) > 0 {
				reasons = append(reasons, fmt.Sprintf("CSS块[%s]: %s", b.OnlyId, strings.Join(hits, "; ")))
			}
		case models.CodeTypeFooter:
			if hits := scanHTML(b.Code); len(hits) > 0 {
				reasons = append(reasons, fmt.Sprintf("页脚块[%s]: %s", b.OnlyId, strings.Join(hits, "; ")))
			}
		}
	}
	return len(reasons) == 0, reasons
}

var (
	reJSExec       = regexp.MustCompile(`(?i)\beval\s*\(`)
	reJSNewFunc    = regexp.MustCompile(`(?i)\bnew\s+Function\s*\(`)
	reJSTimeoutStr = regexp.MustCompile("(?i)\\bset(?:Timeout|Interval)\\s*\\(\\s*(['\"\x60])")
	reJSDOMInject  = regexp.MustCompile(`(?i)(?:document\.write|innerHTML\s*=|outerHTML\s*=|<\s*script)`)
	reJSExecAtob   = regexp.MustCompile(`(?i)(?:eval|new\s+Function)\s*\(\s*atob`)
	reJSExtNav     = regexp.MustCompile(`(?i)(?:location\.(?:href|replace|assign)|window\.open)\s*(?:\(?\s*=?)\s*['"]https?://`)
	reJSExtSend    = regexp.MustCompile(`(?i)\b(?:fetch|XMLHttpRequest|navigator\.sendBeacon)\s*\(\s*['"]https?://`)
	reJSWebSocket  = regexp.MustCompile(`(?i)\bnew\s+WebSocket\s*\(\s*['"]wss?://`)
	reJSMining     = regexp.MustCompile(`(?i)(coinhive|minero|jsecoin|cryptonight)`)

	reCSSExtURL   = regexp.MustCompile(`(?i)url\s*\(\s*['"]?https?://`)
	reCSSImport   = regexp.MustCompile(`(?i)@import`)
	reCSSPosition = regexp.MustCompile(`(?i)position\s*:\s*(?:fixed|absolute)`)
	reCSSOpacity0 = regexp.MustCompile(`(?i)opacity\s*:\s*0`)
	reCSSZIndex   = regexp.MustCompile(`(?i)z-index`)

	// reHTMLScript  = regexp.MustCompile(`(?i)<\s*script`)
	reHTMLIframe  = regexp.MustCompile(`(?i)<\s*(?:iframe|object|embed)`)
	reHTMLForm    = regexp.MustCompile(`(?i)<\s*form`)
	reHTMLEvent   = regexp.MustCompile(`(?i)\son[a-z]+\s*=`)
	reHTMLExtHref = regexp.MustCompile(`(?i)(?:href|src)\s*=\s*['"]https?://`)
)

func scanJS(code string) []string {
	var hits []string
	if reJSExec.MatchString(code) {
		hits = append(hits, "动态执行 eval()")
	}
	if reJSNewFunc.MatchString(code) {
		hits = append(hits, "动态执行 new Function()")
	}
	if reJSTimeoutStr.MatchString(code) {
		hits = append(hits, "setTimeout/setInterval 传入字符串")
	}
	if reJSDOMInject.MatchString(code) {
		hits = append(hits, "疑似注入 DOM(<script>/document.write/innerHTML)")
	}
	if reJSExecAtob.MatchString(code) {
		hits = append(hits, "解码后动态执行(atob+eval/Function)")
	}
	if reJSExtNav.MatchString(code) {
		hits = append(hits, "跳转/打开外部链接")
	}
	if reJSExtSend.MatchString(code) {
		hits = append(hits, "向外部域名发起请求(fetch/XHR/Beacon)")
	}
	if reJSWebSocket.MatchString(code) {
		hits = append(hits, "外部 WebSocket 连接")
	}
	if reJSMining.MatchString(code) {
		hits = append(hits, "疑似挖矿特征")
	}
	return hits
}

func scanCSS(code string) []string {
	var hits []string
	if reCSSExtURL.MatchString(code) {
		hits = append(hits, "外链资源 url(http(s)://)")
	}
	if reCSSImport.MatchString(code) {
		hits = append(hits, "@import 外部样式")
	}
	// 点击劫持：固定/绝对定位 + 透明 + 高 z-index 三者同时出现
	if reCSSPosition.MatchString(code) && reCSSOpacity0.MatchString(code) && reCSSZIndex.MatchString(code) {
		hits = append(hits, "疑似点击劫持(定位+透明+高z-index)")
	}
	return hits
}

func scanHTML(code string) []string {
	var hits []string
	// if reHTMLScript.MatchString(code) {
	// 	hits = append(hits, "页脚禁止 <script>")
	// }
	if reHTMLIframe.MatchString(code) {
		hits = append(hits, "含 <iframe>/<object>/<embed>")
	}
	if reHTMLForm.MatchString(code) {
		hits = append(hits, "含 <form>(疑似钓鱼)")
	}
	if reHTMLEvent.MatchString(code) {
		hits = append(hits, "内联事件处理器(on*)")
	}
	if reHTMLExtHref.MatchString(code) {
		hits = append(hits, "外链 href/src")
	}
	return hits
}
