# YesPOS Go SDK

SDK Go untuk seluruh endpoint `/api/client` YesPOS. Package ini mandiri dan tidak bergantung pada package internal aplikasi YesPOS.

## Instalasi

```bash
go get github.com/syaifudin21/yespos-go
```

```go
import yespos "github.com/syaifudin21/yespos-go"
```

SDK membutuhkan Go 1.24 atau lebih baru.

## Inisialisasi

```go
client := yespos.NewClient(
	"https://api.example.com", // Base URL tanpa trailing slash
	"your-va",                 // header va
	"your-secret",             // secret credential client
)
```

Semua request otomatis menyertakan header `va`, `signature`, `Content-Type`, dan `Accept`. Signature mengikuti middleware server:

```text
SHA256(body)
string_to_sign = METHOD:VA:SHA256(body):SECRET
signature = HMAC-SHA256(string_to_sign, SECRET)
```

Gunakan bahasa response tertentu dengan `client.AcceptLanguage = "id"`. Untuk timeout, tracing, proxy, atau transport khusus, isi `client.HTTPClient` dengan `*http.Client` milik aplikasi.

## Response dan error

Endpoint biasa mengembalikan `yespos.Response[T]`, sedangkan endpoint pagination mengembalikan `yespos.PaginatedResponse[T]`.

```go
result, err := client.Test(ctx)
if err != nil {
	var apiErr *yespos.HTTPError
	if errors.As(err, &apiErr) {
		log.Printf("status=%d message=%s", apiErr.StatusCode, apiErr.Response.Message)
	}
	return err
}
fmt.Println(result.Data.Name)
```

`HTTPError` berisi `StatusCode` dan response error server. Error network, encode request, dan decode response dikembalikan sebagai error dengan prefix `yespos:`. Semua function menerima `context.Context`; gunakan `context.WithTimeout` untuk membatasi durasi request.

## Daftar function lengkap

| Function | Method dan endpoint | Return |
|---|---|---|
| `Test(ctx)` | `GET /api/client/test` | `Response[TestResponse]` |
| `Register(ctx, request)` | `POST /api/client/register` | `Response[ClientRegisterResponse]` |
| `GetToken(ctx, request)` | `POST /api/client/getToken` | `Response[TokenResponse]` |
| `CreateOrder(ctx, request)` | `POST /api/client/orders` | `Response[OrderResponse]` |
| `ListOrders(ctx, page, limit)` | `GET /api/client/orders` | `Response[OrderListResponse]` |
| `GetOrder(ctx, uuid)` | `GET /api/client/orders/{uuid}` | `Response[OrderResponse]` |
| `ListMembers(ctx)` | `GET /api/client/member` | `Response[[]Member]` |
| `ListOutlets(ctx)` | `GET /api/client/outlets` | `Response[[]Outlet]` |
| `CreateOutlet(ctx, request)` | `POST /api/client/outlets` | `Response[CreateOutletResponse]` |
| `NearbyOutlets(ctx, lat, long, page, limit, deliveryOnly, type, search, category)` | `GET /api/client/outlets/nearby` | `PaginatedResponse[[]NearbyOutlet]` |
| `GetOutlet(ctx, uuid)` | `GET /api/client/outlets/{outlet_uuid}` | `Response[OutletDetail]` |
| `ListInventories(ctx, outletUUID)` | `GET /api/client/outlets/{outlet_uuid}/inventories` | `Response[[]Inventory]` |
| `CreateInventory(ctx, outletUUID, request)` | `POST /api/client/outlets/{outlet_uuid}/inventories` | `Response[Inventory]` |
| `UpdateInventory(ctx, outletUUID, inventoryUUID, request)` | `PUT /api/client/outlets/{outlet_uuid}/inventories/{inventory_uuid}` | `Response[Inventory]` |
| `DeleteInventory(ctx, outletUUID, inventoryUUID)` | `DELETE /api/client/outlets/{outlet_uuid}/inventories/{inventory_uuid}` | `Response[any]` |
| `ListTransactions(ctx, page, limit, outletUUID)` | `GET /api/client/transaction` | `Response[TransactionList]` |
| `GetTransaction(ctx, uuid)` | `GET /api/client/transaction/{uuid}` | `Response[TransactionDetail]` |
| `ListProducts(ctx, outletUUID)` | `GET /api/client/outlets/{outlet_uuid}/products` | `Response[[]ProductCategory]` |
| `AdjustStock(ctx, outletUUID, request)` | `POST /api/client/outlets/{outlet_uuid}/stocks/adjust` | `Response[AdjustStockResponse]` |
| `GetStocks(ctx, outletUUID, inventoryUUID)` | `GET /api/client/outlets/{outlet_uuid}/stocks` | `Response[[]StockAggregate]` |
| `GetStockHistory(ctx, outletUUID, productUUID, inventoryUUID, page, limit)` | `GET /api/client/outlets/{outlet_uuid}/stocks/{product_uuid}` | `PaginatedResponse[[]StockMovement]` |
| `GenerateVoucherCode(ctx, request)` | `POST /api/client/voucher/generate-code` | `Response[VoucherCodeResponse]` |

Untuk query optional, kirim string kosong. Nilai `page` dan `limit` mengikuti default server jika nol atau tidak valid.

## Contoh penggunaan

### Register dan token

```go
registered, err := client.Register(ctx, yespos.ClientRegisterRequest{
	Name: "Demo User", Email: "user@example.com",
	Password: "secret-password", Phone: "081234567890",
})

token, err := client.GetToken(ctx, yespos.GetTokenRequest{Email: registered.Data.Email})
fmt.Println(token.Data.Token)
```

### Membuat order

```go
result, err := client.CreateOrder(ctx, yespos.CreateOrderRequest{
	TotalPayment: 100000, PaymentMethod: "qris",
	Name: "Customer", Email: "customer@example.com", Phone: "081234567890",
	Items: []yespos.OrderItem{{Name: "Coffee", Quantity: 2, Price: 25000, Total: 50000}},
})
```

### Pagination

```go
orders, err := client.ListOrders(ctx, 1, 20)
if err == nil {
	fmt.Printf("page=%d total=%d\n", orders.Data.Page, orders.Data.Total)
}
```

### Voucher Laris Manis

```go
result, err := client.GenerateVoucherCode(ctx, yespos.VoucherCodeRequest{
	Phone: "081234567890", Code: "kodelarismanis",
})
fmt.Println(result.Data.Code)
```

### Endpoint baru yang belum memiliki method typed

```go
var result yespos.Response[map[string]any]
err := client.Do(ctx, http.MethodGet, "/api/client/new-endpoint", nil, nil, &result)
```

## Testing

```bash
GOCACHE=/tmp/yespos-go-cache go test ./...
```

Unit test mencakup algoritma signature dan request `GenerateVoucherCode`.
