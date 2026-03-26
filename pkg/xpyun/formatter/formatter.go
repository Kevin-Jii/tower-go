package formatter

import (
	"fmt"
	"strconv"
	"github.com/Kevin-Jii/tower-go/pkg/xpyun/util"
)

const ROW_MAX_CHAR_LEN = 32
const MAX_NAME_CHAR_LEN = 20
const LAST_ROW_MAX_NAME_CHAR_LEN = 16
const MAX_QUANTITY_CHAR_LEN = 6
const MAX_PRICE_CHAR_LEN = 6

const ROW_MAX_CHAR_LEN80 = 48
const MAX_NAME_CHAR_LEN80 = 27
const LAST_ROW_MAX_NAME_CHAR_LEN80 = 24
const MAX_QUANTITY_CHAR_LEN80 = 7

var orderNameEmpty = util.StrRepeat(" ", MAX_NAME_CHAR_LEN)

// FormatPrintOrderItem 格式化菜品列表（用于58mm打印机）
// 58mm打印机一行可打印32个字符，汉字按2个字符算
// 分3列：名称(20字符) 数量(6字符) 单价(6字符)
func FormatPrintOrderItem(foodName string, quantity int, price float64) string {
	foodNameLen := util.CalcGbkLenForPrint(foodName)
	quantityStr := strconv.Itoa(quantity)
	quantityLen := util.CalcAsciiLenForPrint(quantityStr)
	priceStr := fmt.Sprintf("%.2f", price)
	priceLen := util.CalcAsciiLenForPrint(priceStr)

	result := foodName
	mod := foodNameLen % ROW_MAX_CHAR_LEN
	if mod <= LAST_ROW_MAX_NAME_CHAR_LEN {
		result = result + util.StrRepeat(" ", MAX_NAME_CHAR_LEN-mod)
	} else {
		result = result + "<BR>" + orderNameEmpty
	}

	result = result + quantityStr + util.StrRepeat(" ", MAX_QUANTITY_CHAR_LEN-quantityLen)
	result = result + priceStr + util.StrRepeat(" ", MAX_PRICE_CHAR_LEN-priceLen)
	result = result + "<BR>"

	return result
}

// FormatPrintOrderItem80 格式化菜品列表（用于80mm打印机）
// 80mm打印机一行可打印48个字符，汉字按2个字符算
// 分4列：名称(27字符) 数量(7字符) 单价(7字符) 总价(7字符)
func FormatPrintOrderItem80(foodName string, quantity int, price float64, total float64) string {
	foodNameLen := util.CalcGbkLenForPrint(foodName)
	quantityStr := strconv.Itoa(quantity)
	quantityLen := util.CalcAsciiLenForPrint(quantityStr)
	priceStr := fmt.Sprintf("%.2f", price)
	priceLen := util.CalcAsciiLenForPrint(priceStr)
	totalStr := fmt.Sprintf("%.2f", total)
	totalLen := util.CalcAsciiLenForPrint(totalStr)

	result := foodName
	mod := foodNameLen % ROW_MAX_CHAR_LEN80
	if mod <= LAST_ROW_MAX_NAME_CHAR_LEN80 {
		result = result + util.StrRepeat(" ", MAX_NAME_CHAR_LEN80-mod)
		result = result + quantityStr + util.StrRepeat(" ", MAX_QUANTITY_CHAR_LEN80-quantityLen)
		result = result + priceStr + util.StrRepeat(" ", MAX_QUANTITY_CHAR_LEN80-priceLen)
		result = result + totalStr + util.StrRepeat(" ", MAX_QUANTITY_CHAR_LEN80-totalLen)
	} else {
		foods := splitStrArray(foodName, 12)
		tempStr := ""
		for i := 0; i < len(foods); i++ {
			if i == 0 {
				mod := util.CalcGbkLenForPrint(foods[i])
				tempStr = tempStr + foods[i] + util.StrRepeat(" ", MAX_NAME_CHAR_LEN80-mod)
				tempStr = tempStr + quantityStr + util.StrRepeat(" ", MAX_QUANTITY_CHAR_LEN80-quantityLen)
				tempStr = tempStr + priceStr + util.StrRepeat(" ", MAX_QUANTITY_CHAR_LEN80-priceLen)
				tempStr = tempStr + totalStr + util.StrRepeat(" ", MAX_QUANTITY_CHAR_LEN80-totalLen)
				tempStr = tempStr + "<BR>"
			} else if i == len(foods)-1 {
				tempStr = tempStr + foods[i]
			} else {
				tempStr = tempStr + foods[i] + "<BR>"
			}
		}
		result = tempStr
	}

	result = result + "<BR>"
	return result
}

func splitStrArray(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}