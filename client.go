// Package sdk provides a Go client for the YesPOS client API.
package sdk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL, VA, Secret, AcceptLanguage string
	HTTPClient                          *http.Client
}

func NewClient(baseURL, va, secret string) *Client {
	return &Client{BaseURL: strings.TrimRight(baseURL, "/"), VA: va, Secret: secret, HTTPClient: http.DefaultClient}
}

type HTTPError struct {
	StatusCode int
	Response   ErrorResponse
}

func (e *HTTPError) Error() string {
	if e.Response.Message != "" {
		return fmt.Sprintf("yespos: HTTP %d: %s", e.StatusCode, e.Response.Message)
	}
	return fmt.Sprintf("yespos: HTTP %d", e.StatusCode)
}

func (c *Client) signature(method string, body []byte) string {
	h := sha256.Sum256(body)
	s := method + ":" + c.VA + ":" + hex.EncodeToString(h[:]) + ":" + c.Secret
	m := hmac.New(sha256.New, []byte(c.Secret))
	_, _ = m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// Do executes a signed request and can be used for newly added server endpoints.
func (c *Client) Do(ctx context.Context, method, endpoint string, query url.Values, body, out any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	payload := []byte("{}")
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("yespos: encode request: %w", err)
		}
	}
	u, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return fmt.Errorf("yespos: invalid URL: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("yespos: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("va", c.VA)
	req.Header.Set("signature", c.signature(method, payload))
	if c.AcceptLanguage != "" {
		req.Header.Set("Accept-Language", c.AcceptLanguage)
	}
	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	res, err := hc.Do(req)
	if err != nil {
		return fmt.Errorf("yespos: request: %w", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("yespos: read response: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var apiErr ErrorResponse
		_ = json.Unmarshal(data, &apiErr)
		return &HTTPError{StatusCode: res.StatusCode, Response: apiErr}
	}
	if out != nil && len(data) > 0 {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("yespos: decode response: %w", err)
		}
	}
	return nil
}

func do[T any](c *Client, ctx context.Context, method, endpoint string, query url.Values, body any) (Response[T], error) {
	var r Response[T]
	err := c.Do(ctx, method, endpoint, query, body, &r)
	return r, err
}
func doPage[T any](c *Client, ctx context.Context, method, endpoint string, query url.Values) (PaginatedResponse[T], error) {
	var r PaginatedResponse[T]
	err := c.Do(ctx, method, endpoint, query, nil, &r)
	return r, err
}
func params(values map[string]string) url.Values {
	q := url.Values{}
	for k, v := range values {
		if v != "" {
			q.Set(k, v)
		}
	}
	return q
}
func path(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, part := range parts[1:] {
		result += "/" + url.PathEscape(part)
	}
	return result
}

func (c *Client) Test(ctx context.Context) (Response[TestResponse], error) {
	return do[TestResponse](c, ctx, http.MethodGet, "/api/client/test", nil, nil)
}
func (c *Client) Register(ctx context.Context, in ClientRegisterRequest) (Response[ClientRegisterResponse], error) {
	return do[ClientRegisterResponse](c, ctx, http.MethodPost, "/api/client/register", nil, in)
}
func (c *Client) GetToken(ctx context.Context, in GetTokenRequest) (Response[TokenResponse], error) {
	return do[TokenResponse](c, ctx, http.MethodPost, "/api/client/getToken", nil, in)
}
func (c *Client) CreateOrder(ctx context.Context, in CreateOrderRequest) (Response[OrderResponse], error) {
	return do[OrderResponse](c, ctx, http.MethodPost, "/api/client/orders", nil, in)
}
func (c *Client) ListOrders(ctx context.Context, page, limit int) (Response[OrderListResponse], error) {
	return do[OrderListResponse](c, ctx, http.MethodGet, "/api/client/orders", params(map[string]string{"page": fmt.Sprint(page), "limit": fmt.Sprint(limit)}), nil)
}
func (c *Client) GetOrder(ctx context.Context, uuid string) (Response[OrderResponse], error) {
	return do[OrderResponse](c, ctx, http.MethodGet, path("/api/client/orders", uuid), nil, nil)
}
func (c *Client) ListMembers(ctx context.Context) (Response[[]Member], error) {
	return do[[]Member](c, ctx, http.MethodGet, "/api/client/member", nil, nil)
}
func (c *Client) ListOutlets(ctx context.Context) (Response[[]Outlet], error) {
	return do[[]Outlet](c, ctx, http.MethodGet, "/api/client/outlets", nil, nil)
}
func (c *Client) CreateOutlet(ctx context.Context, in CreateOutletRequest) (Response[CreateOutletResponse], error) {
	return do[CreateOutletResponse](c, ctx, http.MethodPost, "/api/client/outlets", nil, in)
}
func (c *Client) NearbyOutlets(ctx context.Context, lat, long float64, page, limit int, deliveryOnly bool, typ, search, category string) (PaginatedResponse[[]NearbyOutlet], error) {
	return doPage[[]NearbyOutlet](c, ctx, http.MethodGet, "/api/client/outlets/nearby", params(map[string]string{"lat": fmt.Sprint(lat), "long": fmt.Sprint(long), "page": fmt.Sprint(page), "limit": fmt.Sprint(limit), "is_delivery_available": fmt.Sprint(deliveryOnly), "type": typ, "search": search, "category": category}))
}
func (c *Client) GetOutlet(ctx context.Context, uuid string) (Response[OutletDetail], error) {
	return do[OutletDetail](c, ctx, http.MethodGet, path("/api/client/outlets", uuid), nil, nil)
}
func (c *Client) ListInventories(ctx context.Context, outletUUID string) (Response[[]Inventory], error) {
	return do[[]Inventory](c, ctx, http.MethodGet, path("/api/client/outlets", outletUUID, "inventories"), nil, nil)
}
func (c *Client) CreateInventory(ctx context.Context, outletUUID string, in CreateInventoryRequest) (Response[Inventory], error) {
	return do[Inventory](c, ctx, http.MethodPost, path("/api/client/outlets", outletUUID, "inventories"), nil, in)
}
func (c *Client) UpdateInventory(ctx context.Context, outletUUID, inventoryUUID string, in CreateInventoryRequest) (Response[Inventory], error) {
	return do[Inventory](c, ctx, http.MethodPut, path("/api/client/outlets", outletUUID, "inventories", inventoryUUID), nil, in)
}
func (c *Client) DeleteInventory(ctx context.Context, outletUUID, inventoryUUID string) (Response[any], error) {
	return do[any](c, ctx, http.MethodDelete, path("/api/client/outlets", outletUUID, "inventories", inventoryUUID), nil, nil)
}
func (c *Client) ListTransactions(ctx context.Context, page, limit int, outletUUID string) (Response[TransactionList], error) {
	return do[TransactionList](c, ctx, http.MethodGet, "/api/client/transaction", params(map[string]string{"page": fmt.Sprint(page), "limit": fmt.Sprint(limit), "outlet_uuid": outletUUID}), nil)
}
func (c *Client) GetTransaction(ctx context.Context, uuid string) (Response[TransactionDetail], error) {
	return do[TransactionDetail](c, ctx, http.MethodGet, path("/api/client/transaction", uuid), nil, nil)
}
func (c *Client) ListProducts(ctx context.Context, outletUUID string) (Response[[]ProductCategory], error) {
	return do[[]ProductCategory](c, ctx, http.MethodGet, path("/api/client/outlets", outletUUID, "products"), nil, nil)
}
func (c *Client) AdjustStock(ctx context.Context, outletUUID string, in AdjustStockRequest) (Response[AdjustStockResponse], error) {
	return do[AdjustStockResponse](c, ctx, http.MethodPost, path("/api/client/outlets", outletUUID, "stocks", "adjust"), nil, in)
}
func (c *Client) GetStocks(ctx context.Context, outletUUID, inventoryUUID string) (Response[[]StockAggregate], error) {
	return do[[]StockAggregate](c, ctx, http.MethodGet, path("/api/client/outlets", outletUUID, "stocks"), params(map[string]string{"inventory_uuid": inventoryUUID}), nil)
}
func (c *Client) GetStockHistory(ctx context.Context, outletUUID, productUUID, inventoryUUID string, page, limit int) (PaginatedResponse[[]StockMovement], error) {
	return doPage[[]StockMovement](c, ctx, http.MethodGet, path("/api/client/outlets", outletUUID, "stocks", productUUID), params(map[string]string{"inventory_uuid": inventoryUUID, "page": fmt.Sprint(page), "limit": fmt.Sprint(limit)}))
}
func (c *Client) GenerateVoucherCode(ctx context.Context, in VoucherCodeRequest) (Response[VoucherCodeResponse], error) {
	return do[VoucherCodeResponse](c, ctx, http.MethodPost, "/api/client/voucher/generate-code", nil, in)
}
