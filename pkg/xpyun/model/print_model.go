package model

// PrintRequest 打印小票请求
type PrintRequest struct {
	RestRequest `json:",inline"`
	Sn          string  `json:"sn"`
	Content     string  `json:"content"`
	Copies      int     `json:"copies,omitempty"`
	Mode        int     `json:"mode,omitempty"`
	PayType     int     `json:"payType,omitempty"`
	PayMode     int     `json:"payMode,omitempty"`
	Money       float64 `json:"money,omitempty"`
	Voice       int     `json:"voice,omitempty"`
	Idempotent  string  `json:"idempotent,omitempty"`
	ExpiresIn   int     `json:"expiresIn,omitempty"`
	Cutter      int     `json:"cutter,omitempty"`
	BackurlFlag int     `json:"backurlFlag,omitempty"`
	Attached    string  `json:"attached,omitempty"`
	TtsVoiceType    string `json:"ttsVoiceType,omitempty"`
	TtsVoiceTime    int    `json:"ttsVoiceTime,omitempty"`
	TtsVoiceInterval int   `json:"ttsVoiceInterval,omitempty"`
	Tts          string  `json:"tts,omitempty"`
}

// PrintLabelRequest 打印标签请求
type PrintLabelRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
	Content     string `json:"content"`
	Copies      int    `json:"copies,omitempty"`
	Mode        int    `json:"mode,omitempty"`
	Voice       int    `json:"voice,omitempty"`
	Idempotent  string `json:"idempotent,omitempty"`
	ExpiresIn   int    `json:"expiresIn,omitempty"`
	BackurlFlag int    `json:"backurlFlag,omitempty"`
	Attached    string `json:"attached,omitempty"`
}

// PosRequest POS指令请求
type PosRequest struct {
	RestRequest `json:",inline"`
	Sn          string `json:"sn"`
	Content     string `json:"content"`
	Copies      int    `json:"copies,omitempty"`
	Mode        int    `json:"mode,omitempty"`
	Voice       int    `json:"voice,omitempty"`
	Idempotent  string `json:"idempotent,omitempty"`
}