package trace

import (
	"sync"

	"github.com/google/uuid"
)

const Header = "TRACE-ID"

var _ T = (*Trace)(nil)

// 為了讓外面可以相依於這個 interface，而不是 Trace struct
type T interface {
	i()
	ID() string
	WithRequest(req *Request) *Trace
	WithResponse(resp *Response) *Trace
	AppendDialog(dialog *Dialog) *Trace
	AppendDebug(debug *Debug) *Trace
	AppendSQL(sql *SQL) *Trace
	AppendRedis(redis *Redis) *Trace
}

// Trace 紀錄的參數
// 這裡定義了 Trace 的結構，包含了 Trace 的所有資訊，有隱藏了實作細節
// 在 core.go裡面呼叫的時候必須得使用 trace底下的實作細節，所以那邊需要轉型
// 這裡可以考慮直接相依於 interface 包含所有 function
type Trace struct {
	mux                sync.Mutex
	Identifier         string    `json:"trace_id"`             // 鏈路 ID
	Request            *Request  `json:"request"`              // 請求資訊
	Response           *Response `json:"response"`             // 回應資訊
	ThirdPartyRequests []*Dialog `json:"third_party_requests"` // 第三方服務呼叫資訊
	Debugs             []*Debug  `json:"debugs"`               // 除錯資訊
	SQLs               []*SQL    `json:"sqls"`                 // 執行的 SQL 資訊
	Redis              []*Redis  `json:"redis"`                // 執行的 Redis 資訊
	Success            bool      `json:"success"`              // 請求結果 true 或 false
	CostSeconds        float64   `json:"cost_seconds"`         // 執行時間 (單位：秒)
}

// Request 請求資訊
type Request struct {
	TTL        string      `json:"ttl"`         // 請求逾時時間
	Method     string      `json:"method"`      // 請求方式
	DecodedURL string      `json:"decoded_url"` // 請求位址
	Header     interface{} `json:"header"`      // 請求 Header 資訊
	Body       interface{} `json:"body"`        // 請求 Body 資訊
}

// Response 回應資訊
type Response struct {
	Header          interface{} `json:"header"`                      // Header 資訊
	Body            interface{} `json:"body"`                        // Body 資訊
	BusinessCode    int         `json:"business_code,omitempty"`     // 業務代碼
	BusinessCodeMsg string      `json:"business_code_msg,omitempty"` // 提示訊息
	HttpCode        int         `json:"http_code"`                   // HTTP 狀態碼
	HttpCodeMsg     string      `json:"http_code_msg"`               // HTTP 狀態碼說明
	CostSeconds     float64     `json:"cost_seconds"`                // 執行時間 (單位：秒)
}

// New 建立一個新的 Trace，若 ID 為空則自動生成 UUID
func New(id string) *Trace {
	if id == "" {
		id = uuid.New().String()
	}

	return &Trace{
		Identifier: id,
	}
}

func (t *Trace) i() {}

// ID 唯一識別碼
func (t *Trace) ID() string {
	return t.Identifier
}

// WithRequest 設定 request
func (t *Trace) WithRequest(req *Request) *Trace {
	t.Request = req
	return t
}

// WithResponse 設定 response
func (t *Trace) WithResponse(resp *Response) *Trace {
	t.Response = resp
	return t
}

// AppendDialog 安全地追加內部呼叫過程 (Dialog)
func (t *Trace) AppendDialog(dialog *Dialog) *Trace {
	if dialog == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.ThirdPartyRequests = append(t.ThirdPartyRequests, dialog)
	return t
}

// AppendDebug 追加除錯訊息
func (t *Trace) AppendDebug(debug *Debug) *Trace {
	if debug == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.Debugs = append(t.Debugs, debug)
	return t
}

// AppendSQL 追加 SQL 資訊
func (t *Trace) AppendSQL(sql *SQL) *Trace {
	if sql == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.SQLs = append(t.SQLs, sql)
	return t
}

// AppendRedis 追加 Redis 資訊
func (t *Trace) AppendRedis(redis *Redis) *Trace {
	if redis == nil {
		return t
	}

	t.mux.Lock()
	defer t.mux.Unlock()

	t.Redis = append(t.Redis, redis)
	return t
}
