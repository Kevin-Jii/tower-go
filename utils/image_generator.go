package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	imageWidth    = 800
	headerHeight  = 80
	rowHeight     = 50
	padding       = 20
	fontSize      = 14
	titleFontSize = 20
)

var (
	bgColor     = color.RGBA{255, 255, 255, 255} // 白色背景
	headerColor = color.RGBA{67, 160, 71, 255}   // 绿色标题栏
	textColor   = color.RGBA{33, 33, 33, 255}    // 深灰色文字
	lineColor   = color.RGBA{224, 224, 224, 255} // 浅灰色分割线
	whiteColor  = color.RGBA{255, 255, 255, 255} // 白色

	// 字体缓存
	titleFont  *truetype.Font
	normalFont *truetype.Font
)

// 初始化字体
func init() {
	// 尝试加载Windows系统字体（微软雅黑）
	fontPaths := []string{
		"C:/Windows/Fonts/msyh.ttc",                       // 微软雅黑
		"C:/Windows/Fonts/simhei.ttf",                     // 黑体
		"C:/Windows/Fonts/simsun.ttc",                     // 宋体
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", // Linux
		"/System/Library/Fonts/PingFang.ttc",              // macOS
	}

	for _, path := range fontPaths {
		fontBytes, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		font, err := freetype.ParseFont(fontBytes)
		if err == nil {
			titleFont = font
			normalFont = font
			log.Printf("成功加载字体: %s", path)
			return
		}
	}

	log.Println("警告: 无法加载中文字体，图片中的中文可能无法正常显示")
}

// GenerateMenuReportImage 生成报菜记录单PNG图片
func GenerateMenuReportImage(order *model.MenuReportOrder, storeName, userName string) ([]byte, error) {
	// 计算图片总高度
	itemCount := len(order.Items)
	totalHeight := headerHeight + padding*4 + rowHeight*itemCount + 150

	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, totalHeight))

	// 填充背景色
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 绘制标题栏
	headerRect := image.Rect(0, 0, imageWidth, headerHeight)
	draw.Draw(img, headerRect, &image.Uniform{headerColor}, image.Point{}, draw.Src)

	// 绘制标题文字
	drawTextWithFont(img, "报菜记录单", padding, 50, whiteColor, titleFontSize)

	// 绘制基本信息
	currentY := headerHeight + padding*2

	info := []string{
		fmt.Sprintf("门店名称: %s", storeName),
		fmt.Sprintf("操作人员: %s", userName),
		fmt.Sprintf("报菜时间: %s", order.CreatedAt.Format("2006-01-02 15:04:05")),
		fmt.Sprintf("记录单ID: %d", order.ID),
	}

	for _, line := range info {
		drawTextWithFont(img, line, padding, currentY, textColor, fontSize)
		currentY += 25
	}

	// 绘制分割线
	currentY += 10
	drawLine(img, padding, currentY, imageWidth-padding, currentY, lineColor)
	currentY += 20

	// 绘制表头
	drawTextWithFont(img, "菜品明细:", padding, currentY, textColor, fontSize)
	currentY += 30

	// 绘制菜品列表
	for i, item := range order.Items {
		dishName := "未知菜品"
		if item.Dish != nil {
			dishName = item.Dish.Name
		}

		text := fmt.Sprintf("%d. %s x %d", i+1, dishName, item.Quantity)
		if item.Remark != "" {
			text += fmt.Sprintf(" (%s)", item.Remark)
		}

		drawTextWithFont(img, text, padding+20, currentY, textColor, fontSize)
		currentY += rowHeight
	}

	// 绘制总数
	currentY += 10
	totalQty := 0
	for _, item := range order.Items {
		totalQty += item.Quantity
	}
	drawTextWithFont(img, fmt.Sprintf("总数量: %d", totalQty), padding, currentY, textColor, fontSize)
	currentY += 30

	// 绘制备注
	if order.Remark != "" {
		drawLine(img, padding, currentY, imageWidth-padding, currentY, lineColor)
		currentY += 20
		drawTextWithFont(img, "备注:", padding, currentY, textColor, fontSize)
		currentY += 25
		drawTextWithFont(img, order.Remark, padding+20, currentY, color.RGBA{100, 100, 100, 255}, fontSize)
	}

	// 编码为PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// drawTextWithFont 使用TrueType字体绘制文字（支持中文）
func drawTextWithFont(img *image.RGBA, text string, x, y int, col color.Color, size float64) {
	if normalFont == nil {
		// 如果字体加载失败，使用基本方法绘制（仅支持ASCII）
		point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
		d := &font.Drawer{
			Dst: img,
			Src: image.NewUniform(col),
			Dot: point,
		}
		d.DrawString(text)
		return
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(normalFont)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(col))

	pt := freetype.Pt(x, y)
	_, err := c.DrawString(text, pt)
	if err != nil {
		log.Printf("绘制文字失败: %v", err)
	}
}

// drawLine 绘制直线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	for x := x1; x <= x2; x++ {
		img.Set(x, y1, col)
	}
}
