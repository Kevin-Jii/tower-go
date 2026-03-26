package service

import (
	"github.com/Kevin-Jii/tower-go/pkg/xpyun/model"
)

const BASE_URL = "https://open.xpyun.net/api/openapi"

func xpyunPostJson(url string, request interface{}) *model.XPYunResp {
	return HttpPostJson(url, request)
}

// XpYunAddPrinters 添加打印机
func XpYunAddPrinters(request *model.AddPrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/addPrinters"
	return xpyunPostJson(url, request)
}

// XpYunSetVoiceType 设置语音类型
func XpYunSetVoiceType(request model.SetVoiceTypeRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/setVoiceType"
	return xpyunPostJson(url, request)
}

// XpYunPrint 打印小票
func XpYunPrint(request *model.PrintRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/print"
	return xpyunPostJson(url, request)
}

// XpYunPrintLabel 打印标签
func XpYunPrintLabel(request *model.PrintLabelRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/printLabel"
	return xpyunPostJson(url, request)
}

// XpYunDelPrinters 删除打印机
func XpYunDelPrinters(request *model.DelPrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/delPrinters"
	return xpyunPostJson(url, request)
}

// XpYunUpdatePrinter 更新打印机
func XpYunUpdatePrinter(request *model.UpdPrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/updPrinter"
	return xpyunPostJson(url, request)
}

// XpYunDelPrinterQueue 清空打印队列
func XpYunDelPrinterQueue(request *model.ClearPrintOrderRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/delPrinterQueue"
	return xpyunPostJson(url, request)
}

// XpYunQueryOrderState 查询订单状态
func XpYunQueryOrderState(request *model.QueryOrderStateRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/queryOrderState"
	return xpyunPostJson(url, request)
}

// XpYunQueryOrderStatis 查询订单统计
func XpYunQueryOrderStatis(request *model.QueryOrderStatisRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/queryOrderStatis"
	return xpyunPostJson(url, request)
}

// XpYunQueryPrinterStatus 查询打印机状态
// 返回状态: 0-离线, 1-在线正常, 2-在线异常
func XpYunQueryPrinterStatus(request *model.PrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/queryPrinterStatus"
	return xpyunPostJson(url, request)
}

// XpYunQueryPrintersStatus 批量查询打印机状态
func XpYunQueryPrintersStatus(request *model.PrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/queryPrintersStatus"
	return xpyunPostJson(url, request)
}

// XpYunPlayVoice 金额播报
func XpYunPlayVoice(request *model.VoiceRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/playVoice"
	return xpyunPostJson(url, request)
}

// XpYunPos POS指令打印
func XpYunPos(request *model.PosRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/pos"
	return xpyunPostJson(url, request)
}

// XpYunControlBox 钱箱控制
func XpYunControlBox(request *model.PrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/controlBox"
	return xpyunPostJson(url, request)
}

// XpYunPlayVoiceExt 扩展语音播报
func XpYunPlayVoiceExt(request *model.VoicePlayMsgRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/playVoiceExt"
	return xpyunPostJson(url, request)
}

// XpYunPlayCustomVoice 自定义语音播报
func XpYunPlayCustomVoice(request *model.VoicePlayMsgRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/playCustomVoice"
	return xpyunPostJson(url, request)
}

// XpYunUploadLogo 上传LOGO
func XpYunUploadLogo(request *model.UploadLogoRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/uploadLogo"
	return xpyunPostJson(url, request)
}

// XpYunDelUploadLogo 删除LOGO
func XpYunDelUploadLogo(request *model.PrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/delUploadLogo"
	return xpyunPostJson(url, request)
}

// XpYunPrinterInfo 获取打印机信息
func XpYunPrinterInfo(request *model.PrinterRequest) *model.XPYunResp {
	url := BASE_URL + "/xprinter/printerInfo"
	return xpyunPostJson(url, request)
}