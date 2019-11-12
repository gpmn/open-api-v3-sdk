package okex

import (
	"strings"
)

/*
币币账户信息
获取币币账户资产列表(仅展示拥有资金的币对)，查询各币种的余额、冻结和可用等信息。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/accounts

*/
/*available:0.96245415 holds:0 frozen:0 hold:0 id:6193104 currency:USDT balance:0.96245415*/
type SpotAccount struct {
	Currency  string  `json:"currency"`
	AccountID string  `json:"id"`
	Balance   float64 `json:"balance,string"`
	Available float64 `json:"available,string"`
	Hold      float64 `json:"hold,string"`
}

// GetSpotAccounts :
func (client *Client) GetSpotAccounts() (accounts []SpotAccount, err error) {
	if _, err = client.Request(GET, SPOT_ACCOUNTS, nil, &accounts); err != nil {
		return nil, err
	}
	return
}

/*
单一币种账户信息
获取币币账户单个币种的余额、冻结和可用等信息。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/accounts/<currency>
*/
func (client *Client) GetSpotAccountsCurrency(currency string) (*map[string]interface{}, error) {
	r := map[string]interface{}{}
	uri := GetCurrencyUri(SPOT_ACCOUNTS_CURRENCY, currency)

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
账单流水查询
列出账户资产流水。账户资产流水是指导致账户余额增加或减少的行为。流水会分页，并且按时间倒序排序和存储，最新的排在最前面。请参阅分页部分以获取第一页之后的其他记录。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/accounts/<currency>/ledger
*/
func (client *Client) GetSpotAccountsCurrencyLeger(currency string, optionalParams *map[string]string) (*[]map[string]interface{}, error) {
	r := []map[string]interface{}{}

	baseUri := GetCurrencyUri(SPOT_ACCOUNTS_CURRENCY_LEDGER, currency)
	uri := baseUri
	if optionalParams != nil && len(*optionalParams) > 0 {
		uri = BuildParams(baseUri, *optionalParams)
	}

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
获取订单列表
列出您当前所有的订单信息。这个请求支持分页，并且按时间倒序排序和存储，最新的排在最前面。请参阅分页部分以获取第一页之后的其他纪录。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/orders
*/
func (client *Client) GetSpotOrders(status, instrument_id string, options *map[string]string) (*[]map[string]interface{}, error) {
	r := []map[string]interface{}{}

	fullOptions := NewParams()
	fullOptions["instrument_id"] = instrument_id
	fullOptions["status"] = status
	if options != nil && len(*options) > 0 {
		fullOptions["before"] = (*options)["before"]
		fullOptions["after"] = (*options)["after"]
		fullOptions["limit"] = (*options)["limit"]
	}

	uri := BuildParams(SPOT_ORDERS, fullOptions)

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
获取所有未成交订单
列出您当前所有的订单信息。这个请求支持分页，并且按时间倒序排序和存储，最新的排在最前面。请参阅分页部分以获取第一页之后的其他纪录。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/orders_pending
*/
// created_at:2019-04-16T06:14:27.000Z price:0.107 type:limit client_oid: order_id:2664645705550848 timestamp:2019-04-16T06:14:27.000Z filled_size:0 order_type:0 status:open filled_notional:0 funds: instrument_id:IOTA-USDT notional: product_id:IOTA-USDT side:buy size:8.994]
type SpotOrder struct {
	OrderID        string  `json:"order_id"`               // 订单ID
	InstrumentID   string  `json:"instrument_id"`          // 币对名称
	ProductID      string  `json:"product_id"`             // 币对名称
	ClientOid      string  `json:"client_oid"`             // 用户设置的订单ID
	CreatedAt      string  `json:"created_at"`             // 下单时间
	Timestamp      string  `json:"timestamp"`              // 订单创建时间
	Status         string  `json:"status"`                 // all:所有状态 open:未成交 part_filled:部分成交 canceling:撤销中 filled:已成交 cancelled:已撤销 ordering:下单中，failure：下单失败
	Type           string  `json:"type"`                   // limit,market(默认是limit)
	Side           string  `json:"side"`                   // buy or sell
	Funds          string  `json:"funds"`                  // ??
	NotionalStr    string  `json:"notional"`               // 买入金额，市价买入时返回
	FilledNotional float64 `json:"filled_notional,string"` // 已成交金额
	FilledSize     float64 `json:"filled_size,string"`     // 已成交数量
	OrderType      float64 `json:"order_type,string"`      // 参数填数字，0：普通委托（order type不填或填0都是普通委托） 1：只做Maker（Post only） 2：全部成交或立即取消（FOK） 3：立即成交并取消剩余（IOC）
	Price          float64 `json:"price,string"`           // 价格
	Size           float64 `json:"size,string"`            // 交易货币数量
}

func (client *Client) GetSpotOrdersPending(options *map[string]string) (orders []SpotOrder, err error) {
	fullOptions := NewParams()
	uri := SPOT_ORDERS_PENDING
	if options != nil && len(*options) > 0 {
		fullOptions["instrument_id"] = (*options)["instrument_id"]
		fullOptions["from"] = (*options)["from"]
		fullOptions["to"] = (*options)["to"]
		fullOptions["limit"] = (*options)["limit"]
		uri = BuildParams(SPOT_ORDERS_PENDING, fullOptions)
	}

	if _, err := client.Request(GET, uri, nil, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

/*
获取订单信息
通过订单ID获取单个订单信息。

限速规则：20次/2s

HTTP请求
GET /api/spot/v3/orders/<order_id>
或者
GET /api/spot/v3/orders/<client_oid>
*/
func (client *Client) GetSpotOrdersById(instrumentId, orderOrClientId string) (*map[string]interface{}, error) {
	r := map[string]interface{}{}
	uri := strings.Replace(SPOT_ORDERS_BY_ID, "{order_client_id}", orderOrClientId, -1)
	options := NewParams()
	options["instrument_id"] = instrumentId
	uri = BuildParams(uri, options)

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
获取成交明细
获取最近的成交明细表。这个请求支持分页，并且按时间倒序排序和存储，最新的排在最前面。请参阅分页部分以获取第一页之后的其他记录。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/fills
*/
func (client *Client) GetSpotFills(order_id, instrument_id string, options *map[string]string) (*[]map[string]interface{}, error) {
	r := []map[string]interface{}{}

	fullOptions := NewParams()
	fullOptions["instrument_id"] = instrument_id
	fullOptions["order_id"] = order_id
	if options != nil && len(*options) > 0 {
		fullOptions["before"] = (*options)["before"]
		fullOptions["after"] = (*options)["after"]
		fullOptions["limit"] = (*options)["limit"]
	}

	uri := BuildParams(SPOT_FILLS, fullOptions)

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// SpotInstrumentsDesc :
type SpotInstrumentsDesc struct {
	BaseCurrency  string  `json:"base_currency"`
	InstrumentID  string  `json:"instrument_id"`
	QuoteCurrency string  `json:"quote_currency"`
	MinSize       float64 `json:"min_size,string"`
	SizeIncrement float64 `json:"size_increment,string"`
	TickSize      float64 `json:"tick_size,string"`
}

/*
获取币对信息
用于获取行情数据，这组公开接口提供了行情数据的快照，无需认证即可调用。

获取交易币对的列表，查询各币对的交易限制和价格步长等信息。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/instruments
*/
func (client *Client) GetSpotInstruments() ([]SpotInstrumentsDesc, error) {
	var r []SpotInstrumentsDesc

	if _, err := client.Request(GET, SPOT_INSTRUMENTS, nil, &r); err != nil {
		return nil, err
	}
	return r, nil
}

/*
获取深度数据
获取币对的深度列表。这个请求不支持分页，一个请求返回整个深度列表。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/instruments/<instrument_id>/book
*/
func (client *Client) GetSpotInstrumentBook(instrumentId string, optionalParams *map[string]string) (*map[string]interface{}, error) {
	r := map[string]interface{}{}
	uri := GetInstrumentIdUri(SPOT_INSTRUMENT_BOOK, instrumentId)
	if optionalParams != nil && len(*optionalParams) > 0 {
		optionals := NewParams()
		optionals["size"] = (*optionalParams)["size"]
		optionals["depth"] = (*optionalParams)["depth"]
		uri = BuildParams(uri, optionals)
	}

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
获取全部ticker信息
获取平台全部币对的最新成交价、买一价、卖一价和24小时交易量的快照信息。

限速规则：50次/2s
HTTP请求
GET /api/spot/v3/instruments/ticker
*/
func (client *Client) GetSpotInstrumentsTicker() (*[]map[string]interface{}, error) {
	r := []map[string]interface{}{}

	if _, err := client.Request(GET, SPOT_INSTRUMENTS_TICKER, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
获取某个ticker信息
获取币对的最新成交价、买一价、卖一价和24小时交易量的快照信息。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/instruments/<instrument-id>/ticker
*/
func (client *Client) GetSpotInstrumentTicker(instrument_id string) (*map[string]interface{}, error) {
	r := map[string]interface{}{}

	uri := GetInstrumentIdUri(SPOT_INSTRUMENT_TICKER, instrument_id)
	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
获取成交数据
获取币对最新的60条成交列表。这个请求支持分页，并且按时间倒序排序和存储，最新的排在最前面。请参阅分页部分以获取第一页之后的其他纪录。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/instruments/<instrument_id>/trades
*/
func (client *Client) GetSpotInstrumentTrade(instrument_id string, options *map[string]string) (*[]map[string]interface{}, error) {
	r := []map[string]interface{}{}

	uri := GetInstrumentIdUri(SPOT_INSTRUMENT_TRADES, instrument_id)
	fullOptions := NewParams()
	if options != nil && len(*options) > 0 {
		fullOptions["from"] = (*options)["from"]
		fullOptions["to"] = (*options)["to"]
		fullOptions["limit"] = (*options)["limit"]
		uri = BuildParams(uri, fullOptions)
	}

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
获取成交数据
获取币对最新的60条成交列表。这个请求支持分页，并且按时间倒序排序和存储，最新的排在最前面。请参阅分页部分以获取第一页之后的其他纪录。

限速规则：20次/2s
HTTP请求
GET /api/spot/v3/instruments/<instrument_id>/candles
*/
func (client *Client) GetSpotInstrumentCandles(instrument_id string, options *map[string]string) (*[]interface{}, error) {
	r := []interface{}{}

	uri := GetInstrumentIdUri(SPOT_INSTRUMENT_CANDLES, instrument_id)
	fullOptions := NewParams()
	if options != nil && len(*options) > 0 {
		fullOptions["start"] = (*options)["start"]
		fullOptions["end"] = (*options)["end"]
		fullOptions["granularity"] = (*options)["granularity"]
		uri = BuildParams(uri, fullOptions)
	}

	if _, err := client.Request(GET, uri, nil, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
下单
OKEx币币交易提供限价单和市价单两种下单模式(更多下单模式将会在后期支持)。只有当您的账户有足够的资金才能下单。

一旦下单，您的账户资金将在订单生命周期内被冻结。被冻结的资金以及数量取决于订单指定的类型和参数。

限速规则：100次/2s
HTTP请求
POST /api/spot/v3/orders
*/
func (client *Client) PostSpotOrders(side, instrument_id string, optionalOrderInfo *map[string]string) (result *map[string]interface{}, err error) {

	r := map[string]interface{}{}
	postParams := NewParams()
	postParams["side"] = side
	postParams["instrument_id"] = instrument_id

	if optionalOrderInfo != nil && len(*optionalOrderInfo) > 0 {
		postParams["type"] = (*optionalOrderInfo)["type"]
		if val, ok := (*optionalOrderInfo)["client_oid"]; ok {
			postParams["client_oid"] = val
		}
		if val, ok := (*optionalOrderInfo)["margin_trading"]; ok {
			postParams["margin_trading"] = val
		}

		if postParams["type"] == "limit" {
			postParams["price"] = (*optionalOrderInfo)["price"]
			postParams["size"] = (*optionalOrderInfo)["size"]

		} else if postParams["type"] == "market" {
			postParams["size"] = (*optionalOrderInfo)["size"]
			postParams["notional"] = (*optionalOrderInfo)["notional"]

		}
	}

	if _, err := client.Request(POST, SPOT_ORDERS, postParams, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

/*
批量下单
下指定币对的多个订单（每次只能下最多4个币对且每个币对可批量下10个单）。

限速规则：50次/2s
HTTP请求
POST /api/spot/v3/batch_orders
*/
func (client *Client) PostSpotBatchOrders(orderInfos *[]map[string]string) (*map[string]interface{}, error) {
	r := map[string]interface{}{}
	if _, err := client.Request(POST, SPOT_BATCH_ORDERS, orderInfos, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

/*
撤销指定订单
撤销之前下的未完成订单。

限速规则：100次/2s
HTTP请求
POST /api/spot/v3/cancel_orders/<order_id>
或者
POST /api/spot/v3/cancel_orders/<client_oid>
*/
func (client *Client) PostSpotCancelOrders(instrumentId, orderOrClientId string) (*map[string]interface{}, error) {
	r := map[string]interface{}{}

	uri := strings.Replace(SPOT_CANCEL_ORDERS_BY_ID, "{order_client_id}", orderOrClientId, -1)
	options := NewParams()
	options["instrument_id"] = instrumentId

	if _, err := client.Request(POST, uri, options, &r); err != nil {
		return nil, err
	}
	return &r, nil

}

/*
批量撤销订单
撤销指定的某一种或多种币对的所有未完成订单，每个币对可批量撤10个单。

限速规则：50次/2s
HTTP请求
POST /api/spot/v3/cancel_batch_orders
*/
func (client *Client) PostSpotCancelBatchOrders(orderInfos *[]map[string]interface{}) (*map[string]interface{}, error) {
	r := map[string]interface{}{}
	if _, err := client.Request(POST, SPOT_CANCEL_BATCH_ORDERS, orderInfos, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
