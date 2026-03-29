package service

import (
	"github.com/Kevin-Jii/tower-go/pkg/xpyun/model"
)

const DEFAULT_BASE_URL = "https://open.xpyun.net/api/openapi"

func xpyunPostJson(url string, request interface{}) *model.XPYunResp {
	return HttpPostJson(url, request)
}

// XpYunAddPrinters 添加打印机
func XpYunAddPrinters(request *model.AddPrinterRequest) *model.XPYunResp {
	return XpYunAddPrintersWithURL(request, DEFAULT_BASE_URL)
}

func XpYunAddPrintersWithURL(request *model.AddPrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/addPrinters"
	return xpyunPostJson(url, request)
}

// XpYunSetVoiceType 设置语音类型
func XpYunSetVoiceType(request model.SetVoiceTypeRequest) *model.XPYunResp {
	return XpYunSetVoiceTypeWithURL(request, DEFAULT_BASE_URL)
}

func XpYunSetVoiceTypeWithURL(request model.SetVoiceTypeRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/setVoiceType"
	return xpyunPostJson(url, request)
}

// XpYunPrint 打印小票
func XpYunPrint(request *model.PrintRequest) *model.XPYunResp {
	return XpYunPrintWithURL(request, DEFAULT_BASE_URL)
}

func XpYunPrintWithURL(request *model.PrintRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/print"
	return xpyunPostJson(url, request)
}

// XpYunPrintLabel 打印标签
func XpYunPrintLabel(request *model.PrintLabelRequest) *model.XPYunResp {
	return XpYunPrintLabelWithURL(request, DEFAULT_BASE_URL)
}

func XpYunPrintLabelWithURL(request *model.PrintLabelRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/printLabel"
	return xpyunPostJson(url, request)
}

// XpYunDelPrinters 删除打印机
func XpYunDelPrinters(request *model.DelPrinterRequest) *model.XPYunResp {
	return XpYunDelPrintersWithURL(request, DEFAULT_BASE_URL)
}

func XpYunDelPrintersWithURL(request *model.DelPrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/delPrinters"
	return xpyunPostJson(url, request)
}

// XpYunUpdatePrinter 更新打印机
func XpYunUpdatePrinter(request *model.UpdPrinterRequest) *model.XPYunResp {
	return XpYunUpdatePrinterWithURL(request, DEFAULT_BASE_URL)
}

func XpYunUpdatePrinterWithURL(request *model.UpdPrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/updPrinter"
	return xpyunPostJson(url, request)
}

// XpYunDelPrinterQueue 清空打印队列
func XpYunDelPrinterQueue(request *model.ClearPrintOrderRequest) *model.XPYunResp {
	return XpYunDelPrinterQueueWithURL(request, DEFAULT_BASE_URL)
}

func XpYunDelPrinterQueueWithURL(request *model.ClearPrintOrderRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/delPrinterQueue"
	return xpyunPostJson(url, request)
}

// XpYunQueryOrderState 查询订单状态
func XpYunQueryOrderState(request *model.QueryOrderStateRequest) *model.XPYunResp {
	return XpYunQueryOrderStateWithURL(request, DEFAULT_BASE_URL)
}

func XpYunQueryOrderStateWithURL(request *model.QueryOrderStateRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/queryOrderState"
	return xpyunPostJson(url, request)
}

// XpYunQueryOrderStatis 查询订单统计
func XpYunQueryOrderStatis(request *model.QueryOrderStatisRequest) *model.XPYunResp {
	return XpYunQueryOrderStatisWithURL(request, DEFAULT_BASE_URL)
}

func XpYunQueryOrderStatisWithURL(request *model.QueryOrderStatisRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/queryOrderStatis"
	return xpyunPostJson(url, request)
}

// XpYunQueryPrinterStatus 查询打印机状态
// 返回状态: 0-离线, 1-在线正常, 2-在线异常
func XpYunQueryPrinterStatus(request *model.PrinterRequest) *model.XPYunResp {
	return XpYunQueryPrinterStatusWithURL(request, DEFAULT_BASE_URL)
}

func XpYunQueryPrinterStatusWithURL(request *model.PrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/queryPrinterStatus"
	return xpyunPostJson(url, request)
}

// XpYunQueryPrintersStatus 批量查询打印机状态
func XpYunQueryPrintersStatus(request *model.PrinterRequest) *model.XPYunResp {
	return XpYunQueryPrintersStatusWithURL(request, DEFAULT_BASE_URL)
}

func XpYunQueryPrintersStatusWithURL(request *model.PrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/queryPrintersStatus"
	return xpyunPostJson(url, request)
}

// XpYunPlayVoice 金额播报
func XpYunPlayVoice(request *model.VoiceRequest) *model.XPYunResp {
	return XpYunPlayVoiceWithURL(request, DEFAULT_BASE_URL)
}

func XpYunPlayVoiceWithURL(request *model.VoiceRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/playVoice"
	return xpyunPostJson(url, request)
}

// XpYunPos POS指令打印
func XpYunPos(request *model.PosRequest) *model.XPYunResp {
	return XpYunPosWithURL(request, DEFAULT_BASE_URL)
}

func XpYunPosWithURL(request *model.PosRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/pos"
	return xpyunPostJson(url, request)
}

// XpYunControlBox 钱箱控制
func XpYunControlBox(request *model.PrinterRequest) *model.XPYunResp {
	return XpYunControlBoxWithURL(request, DEFAULT_BASE_URL)
}

func XpYunControlBoxWithURL(request *model.PrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/controlBox"
	return xpyunPostJson(url, request)
}

// XpYunPlayVoiceExt 扩展语音播报
func XpYunPlayVoiceExt(request *model.VoicePlayMsgRequest) *model.XPYunResp {
	return XpYunPlayVoiceExtWithURL(request, DEFAULT_BASE_URL)
}

func XpYunPlayVoiceExtWithURL(request *model.VoicePlayMsgRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/playVoiceExt"
	return xpyunPostJson(url, request)
}

// XpYunPlayCustomVoice 自定义语音播报
func XpYunPlayCustomVoice(request *model.VoicePlayMsgRequest) *model.XPYunResp {
	return XpYunPlayCustomVoiceWithURL(request, DEFAULT_BASE_URL)
}

func XpYunPlayCustomVoiceWithURL(request *model.VoicePlayMsgRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/playCustomVoice"
	return xpyunPostJson(url, request)
}

// XpYunUploadLogo 上传LOGO
func XpYunUploadLogo(request *model.UploadLogoRequest) *model.XPYunResp {
	return XpYunUploadLogoWithURL(request, DEFAULT_BASE_URL)
}

func XpYunUploadLogoWithURL(request *model.UploadLogoRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/uploadLogo"
	return xpyunPostJson(url, request)
}

// XpYunDelUploadLogo 删除LOGO
func XpYunDelUploadLogo(request *model.PrinterRequest) *model.XPYunResp {
	return XpYunDelUploadLogoWithURL(request, DEFAULT_BASE_URL)
}

func XpYunDelUploadLogoWithURL(request *model.PrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/delUploadLogo"
	return xpyunPostJson(url, request)
}

// XpYunPrinterInfo 获取打印机信息
func XpYunPrinterInfo(request *model.PrinterRequest) *model.XPYunResp {
	return XpYunPrinterInfoWithURL(request, DEFAULT_BASE_URL)
}

func XpYunPrinterInfoWithURL(request *model.PrinterRequest, baseURL string) *model.XPYunResp {
	url := baseURL + "/xprinter/printerInfo"
	return xpyunPostJson(url, request)
}
