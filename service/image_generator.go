package service

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/fogleman/gg"
	"go.uber.org/zap"
)

// ImageGeneratorService 图片生成服务
type ImageGeneratorService struct {
	rustfsService *RustFSService
	fontPath      string
}

// NewImageGeneratorService 创建图片生成服务
func NewImageGeneratorService(rustfsService *RustFSService) *ImageGeneratorService {
	svc := &ImageGeneratorService{
		rustfsService: rustfsService,
	}
	svc.fontPath = svc.findFontPath()
	return svc
}

// findFontPath 查找可用的中文字体路径
func (s *ImageGeneratorService) findFontPath() string {
	// 优先使用项目内置字体
	fontPaths := []string{
		"pkg/assets/font/NotoSerifCJKsc-VF.ttf",
		"./pkg/assets/font/NotoSerifCJKsc-VF.ttf",
	}

	// 添加系统字体作为备选
	if runtime.GOOS == "windows" {
		fontPaths = append(fontPaths,
			"C:/Windows/Fonts/msyh.ttc",
			"C:/Windows/Fonts/msyhbd.ttc",
			"C:/Windows/Fonts/simhei.ttf",
		)
	} else {
		fontPaths = append(fontPaths,
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
			"/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
		)
	}
	if runtime.GOOS == "darwin" {
		fontPaths = append(fontPaths,
			"/System/Library/Fonts/PingFang.ttc",
			"/System/Library/Fonts/Supplemental/PingFang.ttc",
			"/System/Library/Fonts/Hiragino Sans GB.ttc",
			"/System/Library/Fonts/STHeiti Light.ttc",
			"/System/Library/Fonts/Supplemental/Songti.ttc",
		)
	}

	for _, path := range fontPaths {
		if _, err := os.Stat(path); err == nil {
			logging.LogInfo("找到可用字体", zap.String("path", path))
			return path
		}
	}

	logging.LogWarn("未找到中文字体，图片生成可能显示异常")
	return ""
}

// loadFont 加载字体
func (s *ImageGeneratorService) loadFont(dc *gg.Context, size float64) error {
	if s.fontPath == "" {
		return fmt.Errorf("未找到可用字体")
	}
	return dc.LoadFontFace(s.fontPath, size)
}

// AccountNotifyData 记账通知数据
type AccountNotifyData struct {
	StoreName    string
	AccountNo    string
	ChannelName  string
	AccountDate  string
	OperatorName string
	Items        []AccountItemData
	TotalAmount  float64
	OtherExpense float64
	NetIncome    float64
	ItemCount    int
	Remark       string
	CreateTime   string
}

// AccountItemData 记账明细数据
type AccountItemData struct {
	Name     string
	Quantity float64
	Unit     string
	Amount   float64
}

// 现代风格颜色定义
var (
	colorWhite      = color.White
	colorBgLight    = color.RGBA{244, 247, 252, 255}
	colorPrimary    = color.RGBA{37, 99, 235, 255}
	colorSuccess    = color.RGBA{22, 163, 74, 255}
	colorTextDark   = color.RGBA{15, 23, 42, 255}
	colorTextMedium = color.RGBA{71, 85, 105, 255}
	colorTextLight  = color.RGBA{100, 116, 139, 255}
	colorBorderLine = color.RGBA{226, 232, 240, 255}
	colorTableHead  = color.RGBA{241, 245, 249, 255}
)

func (s *ImageGeneratorService) useFont(dc *gg.Context, size float64) {
	if err := s.loadFont(dc, size); err != nil && logging.SugaredLogger != nil {
		logging.SugaredLogger.Warnw("load font failed, fallback may look bad", "size", size, "error", err)
	}
}

// GenerateAccountNotifyImage 生成记账通知图片（清晰账单卡片风格）
func (s *ImageGeneratorService) GenerateAccountNotifyImage(data *AccountNotifyData) (string, error) {
	scale := 2.0
	width := int(440 * scale)
	cardMargin := 14.0 * scale
	cardPadding := 16.0 * scale
	cardRadius := 12.0 * scale
	rowHeight := 22.0 * scale
	metaHeight := 84.0 * scale
	footerHeight := 92.0 * scale
	maxRows := 10

	showRows := len(data.Items)
	if showRows > maxRows {
		showRows = maxRows
	}
	extraRows := len(data.Items) - showRows
	extraLineHeight := 0.0
	if extraRows > 0 {
		extraLineHeight = rowHeight
	}

	tableHeaderHeight := 26.0 * scale
	tableHeight := tableHeaderHeight + float64(showRows)*rowHeight + extraLineHeight
	headerHeight := 88.0 * scale
	cardWidth := float64(width) - cardMargin*2
	cardHeight := cardPadding*2 + headerHeight + metaHeight + tableHeight + footerHeight
	totalHeight := int(cardHeight + cardMargin*2 + 2)

	dc := gg.NewContext(width, totalHeight)
	dc.SetColor(colorBgLight)
	dc.Clear()

	dc.SetColor(color.RGBA{15, 23, 42, 18})
	dc.DrawRoundedRectangle(cardMargin+3*scale, cardMargin+4*scale, cardWidth, cardHeight, cardRadius)
	dc.Fill()

	dc.SetColor(colorWhite)
	dc.DrawRoundedRectangle(cardMargin, cardMargin, cardWidth, cardHeight, cardRadius)
	dc.Fill()

	y := cardMargin + cardPadding + 12*scale
	dc.SetColor(colorTextDark)
	s.useFont(dc, 18*scale)
	dc.DrawString(cardEllipsis(data.StoreName, 18), cardMargin+cardPadding, y)
	y += 20 * scale

	dc.SetColor(colorTextLight)
	s.useFont(dc, 12*scale)
	dc.DrawString("记账通知", cardMargin+cardPadding, y)

	dc.SetColor(colorPrimary)
	s.useFont(dc, 20*scale)
	dc.DrawStringAnchored(fmt.Sprintf("¥ %.2f", data.TotalAmount), cardMargin+cardWidth-cardPadding, cardMargin+cardPadding+24*scale, 1, 0.5)

	y += 24 * scale
	dc.SetColor(colorBorderLine)
	dc.DrawLine(cardMargin+cardPadding, y, cardMargin+cardWidth-cardPadding, y)
	dc.Stroke()
	y += 16 * scale

	kv := [][2]string{
		{"单号", data.AccountNo},
		{"渠道", data.ChannelName},
		{"操作人", data.OperatorName},
		{"记账日期", data.AccountDate},
		{"创建时间", data.CreateTime},
	}
	kvLeft := cardMargin + cardPadding
	kvRight := cardMargin + cardWidth/2 + 6*scale
	for i, item := range kv {
		if item[1] == "" {
			item[1] = "-"
		}
		lineY := y + float64(i%3)*18*scale
		colX := kvLeft
		if i >= 3 {
			colX = kvRight
		}
		dc.SetColor(colorTextLight)
		s.useFont(dc, 11*scale)
		dc.DrawString(item[0], colX, lineY)
		dc.SetColor(colorTextMedium)
		s.useFont(dc, 11*scale)
		dc.DrawString(cardEllipsis(item[1], 24), colX+42*scale, lineY)
	}
	y += metaHeight

	tableX := cardMargin + cardPadding
	tableW := cardWidth - cardPadding*2
	dc.SetColor(colorTableHead)
	dc.DrawRoundedRectangle(tableX, y, tableW, tableHeaderHeight, 6*scale)
	dc.Fill()
	dc.SetColor(colorTextMedium)
	s.useFont(dc, 11*scale)
	dc.DrawString("# 商品", tableX+8*scale, y+16*scale)
	dc.DrawStringAnchored("数量", tableX+tableW-120*scale, y+16*scale, 1, 0.5)
	dc.DrawStringAnchored("金额", tableX+tableW-10*scale, y+16*scale, 1, 0.5)
	y += tableHeaderHeight

	for i := 0; i < showRows; i++ {
		item := data.Items[i]
		if i%2 == 0 {
			dc.SetColor(color.RGBA{248, 250, 252, 255})
			dc.DrawRectangle(tableX, y, tableW, rowHeight)
			dc.Fill()
		}
		name := strings.TrimSpace(item.Name)
		if name == "" {
			name = fmt.Sprintf("商品%d", i+1)
		}
		dc.SetColor(colorTextDark)
		s.useFont(dc, 11*scale)
		dc.DrawString(cardEllipsis(fmt.Sprintf("%d. %s", i+1, name), 24), tableX+8*scale, y+14*scale)
		qty := fmt.Sprintf("%.2f%s", item.Quantity, strings.TrimSpace(item.Unit))
		dc.SetColor(colorTextMedium)
		dc.DrawStringAnchored(qty, tableX+tableW-120*scale, y+14*scale, 1, 0.5)
		dc.SetColor(colorTextDark)
		dc.DrawStringAnchored(fmt.Sprintf("¥%.2f", item.Amount), tableX+tableW-10*scale, y+14*scale, 1, 0.5)
		y += rowHeight
	}
	if extraRows > 0 {
		dc.SetColor(colorTextLight)
		s.useFont(dc, 10*scale)
		dc.DrawString(fmt.Sprintf("...其余 %d 项未展开", extraRows), tableX+8*scale, y+14*scale)
		y += rowHeight
	}

	y += 8 * scale
	dc.SetColor(colorBorderLine)
	dc.DrawLine(tableX, y, tableX+tableW, y)
	dc.Stroke()
	y += 18 * scale

	dc.SetColor(colorTextMedium)
	s.useFont(dc, 11*scale)
	dc.DrawString(fmt.Sprintf("商品数: %d", maxInt(data.ItemCount, len(data.Items))), tableX, y)
	dc.DrawString(fmt.Sprintf("其它支出: ¥%.2f", data.OtherExpense), tableX+120*scale, y)
	dc.SetColor(colorSuccess)
	dc.DrawStringAnchored(fmt.Sprintf("净收入: ¥%.2f", data.NetIncome), tableX+tableW, y, 1, 0.5)
	y += 16 * scale
	if strings.TrimSpace(data.Remark) != "" {
		dc.SetColor(colorTextLight)
		s.useFont(dc, 10*scale)
		dc.DrawString(cardEllipsis("备注: "+data.Remark, 50), tableX, y)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return "", fmt.Errorf("编码图片失败: %v", err)
	}

	// 上传到RustFS
	if s.rustfsService == nil {
		return "", fmt.Errorf("RustFS服务未启用")
	}

	now := time.Now()
	folder := fmt.Sprintf("notify/account/%s", now.Format("2006/01/02"))
	filename := fmt.Sprintf("%s_%s.png", now.Format("150405"), strings.ReplaceAll(data.AccountNo, "-", ""))

	imageData := buf.Bytes()
	result, err := s.rustfsService.UploadToNotify(folder, filename, bytes.NewReader(imageData), int64(len(imageData)), "image/png")
	if err != nil {
		return "", fmt.Errorf("上传图片失败: %v", err)
	}

	notifyBucket := s.rustfsService.GetNotifyBucket()
	presignedURL, err := s.rustfsService.GetPresignedURLForBucket(notifyBucket, result.Path, 7*24*time.Hour)
	if err != nil {
		logging.LogWarn("生成预签名URL失败，使用公开URL", zap.Error(err))
		presignedURL = result.URL
	}

	logging.LogInfo("记账通知图片生成成功", zap.String("url", presignedURL))
	return presignedURL, nil
}

func cardEllipsis(s string, limit int) string {
	r := []rune(strings.TrimSpace(s))
	if len(r) <= limit {
		return string(r)
	}
	if limit <= 1 {
		return "…"
	}
	return string(r[:limit-1]) + "…"
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
