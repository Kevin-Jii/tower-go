package service

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"math"
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
	colorBgLight    = color.RGBA{245, 247, 250, 255} // #F5F7FA
	colorPrimary    = color.RGBA{64, 158, 255, 255}  // #409EFF 主色蓝
	colorSuccess    = color.RGBA{103, 194, 58, 255}  // #67C23A 成功绿
	colorTextDark   = color.RGBA{48, 49, 51, 255}    // #303133 主文字
	colorTextMedium = color.RGBA{96, 98, 102, 255}   // #606266 常规文字
	colorTextLight  = color.RGBA{144, 147, 153, 255} // #909399 次要文字
	colorBorderLine = color.RGBA{235, 238, 245, 255} // #EBEEF5 边框
	colorTableHead  = color.RGBA{245, 247, 250, 255} // 表头背景
)

// GenerateAccountNotifyImage 生成记账通知图片（优化版电子回单风格）
func (s *ImageGeneratorService) GenerateAccountNotifyImage(data *AccountNotifyData) (string, error) {
	// 2倍分辨率
	scale := 2.0

	// 布局参数
	cardMargin := 16.0 * scale
	cardPadding := 24.0 * scale
	cardRadius := 12.0 * scale
	width := int(420.0 * scale)

	// 计算高度
	headerHeight := 120.0 * scale // 增加头部高度以容纳门店名称和金额
	rowHeight := 50.0 * scale
	serratedHeight := 12.0 * scale // 锯齿装饰高度
	footerHeight := 80.0 * scale   // 底部辅助信息区域

	cardWidth := float64(width) - cardMargin*2
	cardHeight := headerHeight + float64(len(data.Items))*rowHeight + serratedHeight + footerHeight + cardPadding*2
	totalHeight := int(cardHeight + cardMargin*2)

	// 创建画布
	dc := gg.NewContext(width, totalHeight)

	// 背景色
	dc.SetColor(colorBgLight)
	dc.Clear()

	// 绘制卡片阴影
	dc.SetColor(color.RGBA{0, 0, 0, 20})
	drawRoundedRect(dc, cardMargin+4*scale, cardMargin+4*scale, cardWidth, cardHeight, cardRadius)
	dc.Fill()

	// 绘制白色卡片
	dc.SetColor(colorWhite)
	drawRoundedRect(dc, cardMargin, cardMargin, cardWidth, cardHeight, cardRadius)
	dc.Fill()

	// ========== 绘制防伪水印背景 ==========
	drawWatermark(dc, cardMargin, cardMargin, cardWidth, cardHeight, scale)

	// ========== 头部区域（门店名称 + 金额）==========
	y := cardMargin + cardPadding + 20*scale

	// 门店名称（大字标）
	dc.SetColor(colorTextDark)
	s.loadFont(dc, 28*scale)
	dc.DrawStringAnchored(data.StoreName, cardMargin+cardWidth/2, y, 0.5, 0.5)
	y += 50 * scale

	// 金额（大字显示）
	dc.SetColor(colorPrimary)
	s.loadFont(dc, 36*scale)
	totalStr := fmt.Sprintf("¥%.2f", data.TotalAmount)
	dc.DrawStringAnchored(totalStr, cardMargin+cardWidth/2, y, 0.5, 0.5)
	y += 50 * scale

	// ========== 详情区域（左右结构）==========
	leftX := cardMargin + cardPadding
	rightX := cardMargin + cardWidth - cardPadding

	// 商品列表
	for _, item := range data.Items {
		// 左边：商品名称 / 重量
		dc.SetColor(colorTextDark)
		s.loadFont(dc, 15*scale)
		itemText := fmt.Sprintf("%s / %.0f%s", item.Name, item.Quantity, item.Unit)
		dc.DrawString(itemText, leftX, y)

		// 右边：金额
		dc.SetColor(colorTextDark)
		s.loadFont(dc, 15*scale)
		amountStr := fmt.Sprintf("¥%.2f", item.Amount)
		dc.DrawStringAnchored(amountStr, rightX, y, 1, 0.5)

		y += rowHeight
	}

	y += 10 * scale

	// ========== 锯齿装饰线 ==========
	drawSerratedLine(dc, leftX, y, rightX-leftX, serratedHeight, scale)
	y += serratedHeight + 20*scale

	// ========== 底部辅助信息区域 ==========
	// 第一行：编号 + 渠道
	dc.SetColor(color.RGBA{160, 160, 160, 255}) // 更浅的灰色
	s.loadFont(dc, 10*scale)
	infoText1 := fmt.Sprintf("编号: %s  |  渠道: %s", data.AccountNo, data.ChannelName)
	dc.DrawStringAnchored(infoText1, cardMargin+cardWidth/2, y, 0.5, 0.5)
	y += 18 * scale

	// 第二行：操作人 + 日期
	infoText2 := fmt.Sprintf("操作人: %s  |  日期: %s", data.OperatorName, data.AccountDate)
	dc.DrawStringAnchored(infoText2, cardMargin+cardWidth/2, y, 0.5, 0.5)
	y += 18 * scale

	// 第三行：创建时间
	dc.DrawStringAnchored(data.CreateTime, cardMargin+cardWidth/2, y, 0.5, 0.5)

	// ========== 已入账印章 ==========
	stampX := rightX - 50*scale
	stampY := cardMargin + headerHeight + 30*scale
	drawStamp(dc, stampX, stampY, 36*scale, "已入账", colorSuccess)

	// 导出为PNG
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

// drawRoundedRect 绘制圆角矩形
func drawRoundedRect(dc *gg.Context, x, y, w, h, r float64) {
	dc.NewSubPath()
	dc.MoveTo(x+r, y)
	dc.LineTo(x+w-r, y)
	dc.DrawArc(x+w-r, y+r, r, -math.Pi/2, 0)
	dc.LineTo(x+w, y+h-r)
	dc.DrawArc(x+w-r, y+h-r, r, 0, math.Pi/2)
	dc.LineTo(x+r, y+h)
	dc.DrawArc(x+r, y+h-r, r, math.Pi/2, math.Pi)
	dc.LineTo(x, y+r)
	dc.DrawArc(x+r, y+r, r, math.Pi, 3*math.Pi/2)
	dc.ClosePath()
}

// drawRoundedRectTop 绘制顶部圆角矩形
func drawRoundedRectTop(dc *gg.Context, x, y, w, h, r float64) {
	dc.NewSubPath()
	dc.MoveTo(x+r, y)
	dc.LineTo(x+w-r, y)
	dc.DrawArc(x+w-r, y+r, r, -math.Pi/2, 0)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.LineTo(x, y+r)
	dc.DrawArc(x+r, y+r, r, math.Pi, 3*math.Pi/2)
	dc.ClosePath()
}

// drawStamp 绘制印章
func drawStamp(dc *gg.Context, x, y, size float64, text string, col color.RGBA) {
	// 旋转 -15 度
	dc.Push()
	dc.RotateAbout(-0.26, x, y) // 约 -15 度

	// 外圈
	dc.SetColor(col)
	dc.SetLineWidth(size * 0.08)
	dc.DrawCircle(x, y, size)
	dc.Stroke()

	// 内圈
	dc.SetLineWidth(size * 0.04)
	dc.DrawCircle(x, y, size*0.75)
	dc.Stroke()

	// 文字
	dc.SetColor(col)
	// 这里无法直接设置字体，使用已加载的字体
	dc.DrawStringAnchored(text, x, y, 0.5, 0.5)

	dc.Pop()
}

// drawSerratedLine 绘制锯齿装饰线（撕票效果）
func drawSerratedLine(dc *gg.Context, x, y, width, height float64, scale float64) {
	dc.SetColor(colorBorderLine)
	dc.SetLineWidth(1 * scale)

	// 绘制上边线
	dc.DrawLine(x, y, x+width, y)
	dc.Stroke()

	// 绘制锯齿
	toothWidth := 8.0 * scale
	toothCount := int(width / toothWidth)
	actualToothWidth := width / float64(toothCount)

	dc.NewSubPath()
	dc.MoveTo(x, y)
	for i := 0; i < toothCount; i++ {
		xPos := x + float64(i)*actualToothWidth
		// 向下的三角形锯齿
		dc.LineTo(xPos+actualToothWidth/2, y+height)
		dc.LineTo(xPos+actualToothWidth, y)
	}
	dc.SetColor(colorBgLight)
	dc.Fill()

	// 绘制锯齿边框
	dc.SetColor(colorBorderLine)
	dc.SetLineWidth(1 * scale)
	dc.MoveTo(x, y)
	for i := 0; i < toothCount; i++ {
		xPos := x + float64(i)*actualToothWidth
		dc.LineTo(xPos+actualToothWidth/2, y+height)
		dc.LineTo(xPos+actualToothWidth, y)
	}
	dc.Stroke()

	// 绘制下边线
	dc.DrawLine(x, y+height, x+width, y+height)
	dc.Stroke()
}

// drawWatermark 绘制防伪水印背景
func drawWatermark(dc *gg.Context, x, y, width, height float64, scale float64) {
	dc.Push()

	// 设置半透明水印颜色
	dc.SetColor(color.RGBA{64, 158, 255, 15}) // 非常淡的蓝色

	// 旋转角度
	centerX := x + width/2
	centerY := y + height/2
	dc.RotateAbout(-0.35, centerX, centerY) // 约 -20 度

	// 绘制多个水印文字
	watermarkText := "已入账"
	spacing := 80.0 * scale

	// 计算需要绘制的行列数
	rows := int(height/spacing) + 3
	cols := int(width/spacing) + 3

	for row := -1; row < rows; row++ {
		for col := -1; col < cols; col++ {
			wmX := centerX + float64(col-cols/2)*spacing
			wmY := centerY + float64(row-rows/2)*spacing
			dc.DrawStringAnchored(watermarkText, wmX, wmY, 0.5, 0.5)
		}
	}

	dc.Pop()
}
