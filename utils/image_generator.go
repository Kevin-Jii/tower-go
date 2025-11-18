package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/xuri/excelize/v2"
	"golang.org/x/image/font"
	_ "golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

const (
	imageWidth    = 1400 // 进一步增加宽度避免拥挤
	headerHeight  = 100  // 相应调整高度
	rowHeight     = 65   // 增加行高
	padding       = 35   // 增加边距
	fontSize      = 15   // 字体大小
	titleFontSize = 22   // 标题字体
	dpi           = 144  // 提高DPI：72 -> 144（2倍清晰度）
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
	// 获取项目根目录
	rootDir, err := os.Getwd()
	if err != nil {
		log.Printf("无法获取当前工作目录: %v", err)
		rootDir = "."
	}

	// 尝试多个字体路径（按优先级）
	fontPaths := []string{
		// 项目自带字体
		filepath.Join(rootDir, "pkg/assets/font/NotoSerifCJKsc-VF.ttf"),
		filepath.Join(rootDir, "assets/font/NotoSerifCJKsc-VF.ttf"),
		filepath.Join(rootDir, "fonts/NotoSerifCJKsc-VF.ttf"),
		
		// CentOS/RHEL 系统字体路径
		"/usr/share/fonts/wqy-microhei/wqy-microhei.ttc",
		"/usr/share/fonts/wqy-zenhei/wqy-zenhei.ttc",
		"/usr/share/fonts/chinese/TrueType/uming.ttc",
		"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
		"/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
		
		// Ubuntu/Debian 系统字体路径
		"/usr/share/fonts/truetype/droid/DroidSansFallbackFull.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
		"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
		
		// Windows 字体路径
		"C:\\Windows\\Fonts\\msyh.ttc",     // 微软雅黑
		"C:\\Windows\\Fonts\\simhei.ttf",   // 黑体
		"C:\\Windows\\Fonts\\simsun.ttc",   // 宋体
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
			log.Printf("✅ 成功加载字体: %s", path)
			return
		}
	}

	log.Println("⚠️  警告: 无法加载中文字体，图片中的中文可能无法正常显示")
	log.Println("💡 提示: 请在 CentOS 上安装中文字体:")
	log.Println("   sudo yum install -y wqy-microhei-fonts")
	log.Println("   或")
	log.Println("   sudo yum install -y wqy-zenhei-fonts")
}

// GenerateMenuReportImage 生成报菜记录单PNG图片（与Excel样式一致）
func GenerateMenuReportImage(order *model.MenuReportOrder, storeName, userName, storePhone, storeAddress string) ([]byte, error) {
	// 计算图片总高度
	itemCount := len(order.Items)
	headerInfoHeight := 100 // 顶部信息区域（增加高度）
	tableHeaderHeight := 50 // 表头高度
	totalRowHeight := 60    // 合计行高度
	footerHeight := 60      // 底部信息区域（电话和地址）
	totalHeight := headerInfoHeight + tableHeaderHeight + rowHeight*itemCount + totalRowHeight + footerHeight + padding*2

	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, totalHeight))

	// 填充背景色
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 定义颜色
	redColor := color.RGBA{255, 0, 0, 255}          // 红色标题
	grayTextColor := color.RGBA{102, 102, 102, 255} // 灰色文字
	tableHeaderBg := color.RGBA{217, 217, 217, 255} // 表头灰色背景
	borderColor := color.RGBA{204, 204, 204, 255}   // 边框颜色

	currentY := padding + 15

	// 绘制标题（左侧红色，更大字体）
	drawTextWithFont(img, fmt.Sprintf("%s报菜明细", storeName), padding+10, currentY+40, redColor, 28)

	// 绘制右侧信息（只显示申报人）
	rightX := imageWidth - padding - 200
	drawTextWithFont(img, fmt.Sprintf("申报人：%s", userName), rightX, currentY+40, grayTextColor, 14)

	currentY += 80

	// 绘制表头背景
	headerRect := image.Rect(padding, currentY, imageWidth-padding, currentY+tableHeaderHeight)
	draw.Draw(img, headerRect, &image.Uniform{tableHeaderBg}, image.Point{}, draw.Src)

	// 绘制表头边框
	drawLine(img, padding, currentY, imageWidth-padding, currentY, borderColor)
	drawLine(img, padding, currentY+tableHeaderHeight, imageWidth-padding, currentY+tableHeaderHeight, borderColor)
	drawLine(img, padding, currentY, padding, currentY+tableHeaderHeight, borderColor)
	drawLine(img, imageWidth-padding, currentY, imageWidth-padding, currentY+tableHeaderHeight, borderColor)

	// 绘制表头文字
	// 调整列宽，确保内容不拥挤
	colWidths := []int{280, 140, 110, 110, 110, 280} // 各列宽度（更合理的分配）
	headers := []string{"商品名称", "商品规格", "数量", "单价", "金额", "备注"}
	currentX := padding + 20
	for i, header := range headers {
		// 居中显示表头文字
		headerX := currentX + (colWidths[i]-len(header)*int(fontSize))/2
		if headerX < currentX {
			headerX = currentX + 10
		}
		drawTextWithFont(img, header, headerX, currentY+32, textColor, fontSize)
		// 绘制列分隔线
		if i < len(headers)-1 {
			lineX := currentX + colWidths[i]
			drawVerticalLine(img, lineX, currentY, lineX, currentY+tableHeaderHeight, borderColor)
		}
		currentX += colWidths[i]
	}

	currentY += tableHeaderHeight

	// 绘制菜品行并计算总金额
	var totalAmount float64 = 0
	for _, item := range order.Items {
		dishName := "未知菜品"
		price := 0.0
		if item.Dish != nil {
			dishName = item.Dish.Name
			price = item.Dish.Price
		}

		// 计算金额
		amount := price * float64(item.Quantity)
		totalAmount += amount

		// 绘制行下划线（单元格底部边框）
		drawLine(img, padding, currentY+rowHeight, imageWidth-padding, currentY+rowHeight, borderColor)
		drawLine(img, padding, currentY, padding, currentY+rowHeight, borderColor)
		drawLine(img, imageWidth-padding, currentY, imageWidth-padding, currentY+rowHeight, borderColor)

		// 绘制单元格内容
		currentX := padding + 20

		// 商品名称（左对齐）
		drawTextWithFont(img, dishName, currentX+10, currentY+40, textColor, 14)
		currentX += colWidths[0]

		// 商品规格（居中）
		specText := "斤"
		specX := currentX + (colWidths[1]-len(specText)*7)/2
		drawTextWithFont(img, specText, specX, currentY+40, textColor, 14)
		currentX += colWidths[1]

		// 数量（居中）
		qtyText := fmt.Sprintf("%d", item.Quantity)
		qtyX := currentX + (colWidths[2]-len(qtyText)*10)/2
		drawTextWithFont(img, qtyText, qtyX, currentY+40, textColor, 14)
		currentX += colWidths[2]

		// 单价（居中）
		priceText := fmt.Sprintf("%.2f", price)
		priceX := currentX + (colWidths[3]-len(priceText)*8)/2
		drawTextWithFont(img, priceText, priceX, currentY+40, textColor, 14)
		currentX += colWidths[3]

		// 金额（居中）
		amountText := fmt.Sprintf("%.2f", amount)
		amountX := currentX + (colWidths[4]-len(amountText)*8)/2
		drawTextWithFont(img, amountText, amountX, currentY+40, textColor, 14)
		currentX += colWidths[4]

		// 备注（左对齐）
		if item.Remark != "" {
			drawTextWithFont(img, item.Remark, currentX+10, currentY+40, textColor, 14)
		}

		// 绘制列分隔线
		currentX = padding + 20
		for i := 0; i < len(headers)-1; i++ {
			lineX := currentX + colWidths[i]
			drawVerticalLine(img, lineX, currentY, lineX, currentY+rowHeight, borderColor)
			currentX += colWidths[i]
		}

		currentY += rowHeight
	}

	// 绘制总金额行
	// 绘制行背景（使用淡绿色）
	lightGreenBg := color.RGBA{200, 230, 201, 255} // 淡绿色背景
	totalRowRect := image.Rect(padding, currentY, imageWidth-padding, currentY+rowHeight)
	draw.Draw(img, totalRowRect, &image.Uniform{lightGreenBg}, image.Point{}, draw.Src)

	// 绘制边框
	drawLine(img, padding, currentY, imageWidth-padding, currentY, borderColor)
	drawLine(img, padding, currentY+rowHeight, imageWidth-padding, currentY+rowHeight, borderColor)
	drawLine(img, padding, currentY, padding, currentY+rowHeight, borderColor)
	drawLine(img, imageWidth-padding, currentY, imageWidth-padding, currentY+rowHeight, borderColor)

	// 绘制"合计"文字（左侧）
	currentX = padding + 20
	drawTextWithFont(img, "合计", currentX+10, currentY+40, textColor, 15)
	currentX += colWidths[0] + colWidths[1] + colWidths[2] + colWidths[3]

	// 绘制总金额（在金额列）
	totalAmountText := fmt.Sprintf("%.2f", totalAmount)
	totalAmountX := currentX + (colWidths[4]-len(totalAmountText)*9)/2
	drawTextWithFont(img, totalAmountText, totalAmountX, currentY+40, textColor, 15)

	// 绘制列分隔线
	currentX = padding + 20
	for i := 0; i < len(headers)-1; i++ {
		lineX := currentX + colWidths[i]
		drawVerticalLine(img, lineX, currentY, lineX, currentY+rowHeight, borderColor)
		currentX += colWidths[i]
	}

	currentY += rowHeight

	// 绘制底部信息
	currentY += 20

	// 左下角：负责人电话
	phoneText := "负责人电话："
	if storePhone != "" {
		phoneText += storePhone
	} else {
		phoneText += "未设置"
	}
	drawTextWithFont(img, phoneText, padding+20, currentY+20, grayTextColor, 13)

	// 右下角：门店地址
	addressText := "门店地址："
	if storeAddress != "" {
		addressText += storeAddress
	} else {
		addressText += "未设置"
	}
	// 计算地址文本位置（右对齐）
	addressX := imageWidth - padding - len(addressText)*7 - 20
	if addressX < imageWidth/2 {
		addressX = imageWidth / 2
	}
	drawTextWithFont(img, addressText, addressX, currentY+20, grayTextColor, 13)

	// 编码为PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// drawVerticalLine 绘制垂直线（加粗，更清晰）
func drawVerticalLine(img *image.RGBA, x, y1, x2, y2 int, col color.Color) {
	// 绘制2像素宽的线条
	for y := y1; y <= y2; y++ {
		img.Set(x, y, col)
		if x+1 < img.Bounds().Max.X {
			img.Set(x+1, y, col)
		}
	}
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
	c.SetDPI(dpi) // 使用更高的DPI提升清晰度
	c.SetFont(normalFont)
	c.SetFontSize(size)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(col))
	// 启用抗锯齿
	c.SetHinting(font.HintingFull)

	pt := freetype.Pt(x, y)
	_, err := c.DrawString(text, pt)
	if err != nil {
		log.Printf("绘制文字失败: %v", err)
	}
}

// drawLine 绘制直线（加粗，更清晰）
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	// 绘制2像素宽的线条
	for x := x1; x <= x2; x++ {
		img.Set(x, y1, col)
		if y1+1 < img.Bounds().Max.Y {
			img.Set(x, y1+1, col)
		}
	}
}

// GenerateMenuReportExcel 生成报菜记录单Excel文件（流式）
func GenerateMenuReportExcel(order *model.MenuReportOrder, storeName, userName, storePhone, storeAddress string) (io.Reader, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "报菜明细"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("创建工作表失败: %w", err)
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // 删除默认工作表

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 20) // 商品名称
	f.SetColWidth(sheetName, "B", "B", 12) // 商品规格
	f.SetColWidth(sheetName, "C", "C", 10) // 数量
	f.SetColWidth(sheetName, "D", "D", 10) // 单价
	f.SetColWidth(sheetName, "E", "E", 10) // 金额
	f.SetColWidth(sheetName, "F", "F", 25) // 备注

	// 定义样式
	// 标题样式（红色粗体大字）
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   22,
			Color:  "FF0000",
			Family: "微软雅黑",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})

	// 右上角信息样式
	infoStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Color:  "666666",
			Family: "微软雅黑",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
			Vertical:   "center",
		},
	})

	// 表头样式（灰色背景）
	tableHeaderStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
			Family: "微软雅黑",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"D9D9D9"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// 内容样式（浅灰色背景）
	contentStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "微软雅黑",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F2F2F2"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
		},
	})

	// 居中内容样式
	centerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "微软雅黑",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F2F2F2"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
		},
	})

	// 第1行：标题和右侧信息
	f.SetRowHeight(sheetName, 1, 35)
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("%s报菜明细", storeName))
	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)

	// 右侧信息
	f.SetCellValue(sheetName, "D1", fmt.Sprintf("申报人：%s", userName))
	f.SetCellStyle(sheetName, "D1", "F1", infoStyle)

	// 第2行：日期和电话
	f.SetRowHeight(sheetName, 2, 20)
	f.SetCellValue(sheetName, "D2", fmt.Sprintf("申报日期：%s", order.CreatedAt.Format("2006.1.2 15:04:05")))
	f.SetCellStyle(sheetName, "D2", "F2", infoStyle)

	// 第3行：门店负责人电话
	f.SetRowHeight(sheetName, 3, 20)
	phoneText := "门店负责人电话："
	if storePhone != "" {
		phoneText += storePhone
	} else {
		phoneText += "未设置"
	}
	f.SetCellValue(sheetName, "D3", phoneText)
	f.SetCellStyle(sheetName, "D3", "F3", infoStyle)

	// 第4行：表头
	currentRow := 4
	f.SetRowHeight(sheetName, currentRow, 25)
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), "商品名称")
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), "商品规格")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), "数量")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), "单价")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), "金额")
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), "备注")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("F%d", currentRow), tableHeaderStyle)
	currentRow++

	// 菜品明细并计算总金额
	var totalAmount float64 = 0
	for _, item := range order.Items {
		dishName := "未知菜品"
		price := 0.0
		if item.Dish != nil {
			dishName = item.Dish.Name
			price = item.Dish.Price
		}

		// 计算金额
		amount := price * float64(item.Quantity)
		totalAmount += amount

		f.SetRowHeight(sheetName, currentRow, 22)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), dishName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", currentRow), "斤")
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", currentRow), item.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", currentRow), price)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), amount)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", currentRow), item.Remark)

		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("A%d", currentRow), contentStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("B%d", currentRow), fmt.Sprintf("B%d", currentRow), centerStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", currentRow), fmt.Sprintf("C%d", currentRow), centerStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("D%d", currentRow), fmt.Sprintf("D%d", currentRow), centerStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("E%d", currentRow), fmt.Sprintf("E%d", currentRow), centerStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("F%d", currentRow), fmt.Sprintf("F%d", currentRow), contentStyle)

		currentRow++
	}

	// 添加合计行
	f.SetRowHeight(sheetName, currentRow, 25)
	f.MergeCell(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("D%d", currentRow))
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", currentRow), "合计")
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", currentRow), totalAmount)
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", currentRow), fmt.Sprintf("E%d", currentRow), tableHeaderStyle)
	currentRow++

	// 写入到 buffer（流式）
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("写入Excel失败: %w", err)
	}

	return &buf, nil
}

// GenerateMenuReportExcelAndImage 同时生成Excel和图片，返回Excel流和图片字节
func GenerateMenuReportExcelAndImage(order *model.MenuReportOrder, storeName, userName, storePhone, storeAddress string) (excelReader io.Reader, imageBytes []byte, err error) {
	// 生成Excel
	excelReader, err = GenerateMenuReportExcel(order, storeName, userName, storePhone, storeAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("生成Excel失败: %w", err)
	}

	// 生成图片
	imageBytes, err = GenerateMenuReportImage(order, storeName, userName, storePhone, storeAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("生成图片失败: %w", err)
	}

	return excelReader, imageBytes, nil
}
