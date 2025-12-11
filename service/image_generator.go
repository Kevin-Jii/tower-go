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

// GenerateAccountNotifyImage 生成记账通知图片（现代电子回单风格）
func (s *ImageGeneratorService) GenerateAccountNotifyImage(data *AccountNotifyData) (string, error) {
	// 2倍分辨率
	scale := 2.0

	// 布局参数
	cardMargin := 16.0 * scale
	cardPadding := 24.0 * scale
	cardRadius := 12.0 * scale
	width := int(420.0 * scale)

	// 计算高度
	headerHeight := 70.0 * scale
	infoHeight := 90.0 * scale
	tableHeaderHeight := 40.0 * scale
	rowHeight := 44.0 * scale
	footerHeight := 100.0 * scale

	cardWidth := float64(width) - cardMargin*2
	cardHeight := headerHeight + infoHeight + tableHeaderHeight + float64(len(data.Items))*rowHeight + footerHeight + cardPadding
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

	// ========== 头部区域 ==========
	// 蓝色顶部条
	dc.SetColor(colorPrimary)
	drawRoundedRectTop(dc, cardMargin, cardMargin, cardWidth, headerHeight, cardRadius)
	dc.Fill()

	// 标题
	dc.SetColor(colorWhite)
	s.loadFont(dc, 22*scale)
	dc.DrawStringAnchored("记账通知", cardMargin+cardWidth/2, cardMargin+headerHeight/2-8*scale, 0.5, 0.5)

	// 门店名称
	dc.SetColor(color.RGBA{255, 255, 255, 200})
	s.loadFont(dc, 12*scale)
	dc.DrawStringAnchored(data.StoreName, cardMargin+cardWidth/2, cardMargin+headerHeight/2+12*scale, 0.5, 0.5)

	y := cardMargin + headerHeight + cardPadding

	// ========== 信息区域 ==========
	leftX := cardMargin + cardPadding
	rightX := cardMargin + cardWidth - cardPadding

	// 第一行：编号 + 日期
	dc.SetColor(colorTextLight)
	s.loadFont(dc, 11*scale)
	dc.DrawString("编号", leftX, y)
	dc.DrawStringAnchored("日期", rightX, y, 1, 0.5)
	y += 18 * scale

	dc.SetColor(colorTextDark)
	s.loadFont(dc, 13*scale)
	dc.DrawString(data.AccountNo, leftX, y)
	dc.DrawStringAnchored(data.AccountDate, rightX, y, 1, 0.5)
	y += 28 * scale

	// 第二行：渠道 + 操作人
	dc.SetColor(colorTextLight)
	s.loadFont(dc, 11*scale)
	dc.DrawString("渠道", leftX, y)
	dc.DrawStringAnchored("操作人", rightX, y, 1, 0.5)
	y += 18 * scale

	dc.SetColor(colorTextDark)
	s.loadFont(dc, 13*scale)
	dc.DrawString(data.ChannelName, leftX, y)
	dc.DrawStringAnchored(data.OperatorName, rightX, y, 1, 0.5)
	y += 24 * scale

	// ========== 分隔线 ==========
	dc.SetColor(colorBorderLine)
	dc.SetLineWidth(1 * scale)
	dc.DrawLine(leftX, y, rightX, y)
	dc.Stroke()
	y += 16 * scale

	// ========== 表格区域 ==========
	tableX := leftX
	tableWidth := rightX - leftX
	col1Width := tableWidth * 0.5 // 商品名称
	col2Width := tableWidth * 0.2 // 数量
	_ = tableWidth * 0.3          // col3Width 金额

	// 表头背景
	dc.SetColor(colorTableHead)
	dc.DrawRectangle(tableX, y-8*scale, tableWidth, tableHeaderHeight)
	dc.Fill()

	// 表头文字
	dc.SetColor(colorTextLight)
	s.loadFont(dc, 11*scale)
	dc.DrawString("商品名称", tableX+8*scale, y+10*scale)
	dc.DrawStringAnchored("数量", tableX+col1Width+col2Width/2, y+10*scale, 0.5, 0.5)
	dc.DrawStringAnchored("金额", tableX+tableWidth-8*scale, y+10*scale, 1, 0.5)
	y += tableHeaderHeight

	// 商品列表
	for i, item := range data.Items {
		// 斑马纹
		if i%2 == 1 {
			dc.SetColor(color.RGBA{250, 251, 252, 255})
			dc.DrawRectangle(tableX, y-8*scale, tableWidth, rowHeight)
			dc.Fill()
		}

		rowY := y + 14*scale

		// 商品名称（左对齐）
		dc.SetColor(colorTextDark)
		s.loadFont(dc, 13*scale)
		name := item.Name
		if len([]rune(name)) > 10 {
			name = string([]rune(name)[:10]) + "..."
		}
		dc.DrawString(name, tableX+8*scale, rowY)

		// 数量（居中）
		dc.SetColor(colorTextMedium)
		s.loadFont(dc, 12*scale)
		qtyStr := fmt.Sprintf("×%.0f%s", item.Quantity, item.Unit)
		dc.DrawStringAnchored(qtyStr, tableX+col1Width+col2Width/2, rowY, 0.5, 0.5)

		// 金额（右对齐）
		dc.SetColor(colorTextDark)
		s.loadFont(dc, 13*scale)
		amountStr := fmt.Sprintf("¥%.2f", item.Amount)
		dc.DrawStringAnchored(amountStr, tableX+tableWidth-8*scale, rowY, 1, 0.5)

		y += rowHeight
	}

	y += 16 * scale

	// ========== 底部分隔线 ==========
	dc.SetColor(colorBorderLine)
	dc.SetLineWidth(1 * scale)
	dc.DrawLine(leftX, y, rightX, y)
	dc.Stroke()
	y += 20 * scale

	// ========== 合计区域 ==========
	// 笔数
	dc.SetColor(colorTextLight)
	s.loadFont(dc, 12*scale)
	dc.DrawString(fmt.Sprintf("共 %d 项商品", len(data.Items)), leftX, y)

	// 合计金额
	dc.SetColor(colorTextMedium)
	s.loadFont(dc, 14*scale)
	dc.DrawStringAnchored("合计：", rightX-120*scale, y, 1, 0.5)

	dc.SetColor(colorPrimary)
	s.loadFont(dc, 24*scale)
	totalStr := fmt.Sprintf("¥%.2f", data.TotalAmount)
	dc.DrawStringAnchored(totalStr, rightX, y, 1, 0.5)

	y += 36 * scale

	// ========== 已入账印章 ==========
	stampX := rightX - 60*scale
	stampY := y - 20*scale
	drawStamp(dc, stampX, stampY, 36*scale, "已入账", colorSuccess)

	// ========== 底部时间 ==========
	y += 10 * scale
	dc.SetColor(colorTextLight)
	s.loadFont(dc, 10*scale)
	dc.DrawStringAnchored(data.CreateTime, cardMargin+cardWidth/2, y, 0.5, 0.5)

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
