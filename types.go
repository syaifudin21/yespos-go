package sdk

import "time"

// Response is the common YesPOS API response envelope.
type Response[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// PaginatedResponse is used by endpoints that return pagination metadata.
type PaginatedResponse[T any] struct {
	Message    string `json:"message"`
	Data       T      `json:"data"`
	TotalCount int64  `json:"total_count,omitempty"`
	Total      int64  `json:"total,omitempty"`
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type ErrorResponse struct {
	Message     string         `json:"message"`
	Details     map[string]any `json:"details,omitempty"`
	RedirectURL string         `json:"redirect_url,omitempty"`
}

type ClientRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type ClientRegisterResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	VAIPaymu string `json:"va_ipaymu"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
}

type GetTokenRequest struct {
	Email string `json:"email"`
}
type TokenResponse struct {
	Token     string    `json:"token"`
	User      User      `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
	URL       string    `json:"url"`
}
type User struct {
	ID    uint   `json:"id"`
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type TestResponse struct {
	Name string `json:"name"`
}

type OrderItem struct {
	Name     string  `json:"product_name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Total    float64 `json:"total,omitempty"`
}
type SplitItem struct {
	VA     string `json:"va"`
	Amount int64  `json:"amount"`
	Status string `json:"status,omitempty"`
}
type CreateOrderRequest struct {
	TotalPayment  int64       `json:"total_payment"`
	PaymentMethod string      `json:"payment_method"`
	Name          string      `json:"name"`
	Email         string      `json:"email"`
	Phone         string      `json:"phone"`
	NotifyURL     string      `json:"notify_url,omitempty"`
	Split         []SplitItem `json:"split,omitempty"`
	Items         []OrderItem `json:"items,omitempty"`
	AppCode       string      `json:"app_code,omitempty"`
	TerminalCode  string      `json:"terminal_code,omitempty"`
	MerchantCode  string      `json:"merchant_code,omitempty"`
}
type OrderResponse struct {
	UUID            string      `json:"uuid"`
	TotalPayment    int64       `json:"total_payment"`
	PaymentMethod   string      `json:"payment_method"`
	Status          string      `json:"status"`
	QRString        string      `json:"qr_string,omitempty"`
	AppLink         string      `json:"app_link,omitempty"`
	NotifyURL       string      `json:"notify_url,omitempty"`
	Split           []SplitItem `json:"split,omitempty"`
	ReferenceIPaymu string      `json:"reference_ipaymu,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
}
type OrderListResponse struct {
	Data  []OrderResponse `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}

type Member struct {
	IPaymuVA string `json:"ipaymu_va"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	UserUUID string `json:"user_uuid"`
}
type Outlet struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Contact  string `json:"contact"`
	Type     string `json:"type"`
	Category string `json:"category"`
}
type NearbyOutlet struct {
	UUID                string  `json:"uuid"`
	Name                string  `json:"name"`
	Address             string  `json:"address"`
	Type                string  `json:"type"`
	Category            string  `json:"category"`
	Lat                 float64 `json:"lat"`
	Long                float64 `json:"long"`
	DistanceKM          float64 `json:"distance_km"`
	URL                 string  `json:"url"`
	IconURL             *string `json:"icon_url,omitempty"`
	HasActiveShift      bool    `json:"has_active_shift"`
	IsDeliveryAvailable bool    `json:"is_delivery_available"`
}

type CreateOutletRequest struct {
	Name                string   `json:"name"`
	Address             string   `json:"address"`
	Contact             string   `json:"contact,omitempty"`
	Type                string   `json:"type"`
	Category            string   `json:"category,omitempty"`
	Icon                *string  `json:"icon,omitempty"`
	BackgroundImage     *string  `json:"background_image,omitempty"`
	Timezone            *string  `json:"timezone,omitempty"`
	Lat                 *float64 `json:"lat,omitempty"`
	Long                *float64 `json:"long,omitempty"`
	PostalCode          *string  `json:"postal_code,omitempty"`
	IsDeliveryAvailable *bool    `json:"is_delivery_available,omitempty"`
}
type CreateOutletResponse struct {
	Outlet    map[string]any `json:"outlet"`
	Inventory Inventory      `json:"inventory"`
}
type OutletDetail map[string]any

type Inventory struct {
	ID                  uint      `json:"id"`
	UUID                string    `json:"uuid"`
	OutletID            uint      `json:"outlet_id"`
	Name                string    `json:"name"`
	Lat                 *float64  `json:"lat,omitempty"`
	Long                *float64  `json:"long,omitempty"`
	PostalCode          *string   `json:"postal_code,omitempty"`
	BiteshipID          *string   `json:"biteship_id,omitempty"`
	KoperasiHubGudangID *int      `json:"koperasi_hub_gudang_id"`
	UserID              uint      `json:"user_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
type CreateInventoryRequest struct {
	Name                string  `json:"name"`
	Lat                 float64 `json:"lat,omitempty"`
	Long                float64 `json:"long,omitempty"`
	PostalCode          string  `json:"postal_code,omitempty"`
	BiteshipID          string  `json:"biteship_id,omitempty"`
	KoperasiHubGudangID *int    `json:"koperasi_hub_gudang_id,omitempty"`
}
type UpdateInventoryRequest struct {
	Name                *string  `json:"name,omitempty"`
	Lat                 *float64 `json:"lat,omitempty"`
	Long                *float64 `json:"long,omitempty"`
	PostalCode          *string  `json:"postal_code,omitempty"`
	BiteshipID          *string  `json:"biteship_id,omitempty"`
	KoperasiHubGudangID *int     `json:"koperasi_hub_gudang_id,omitempty"`
}

type Transaction struct {
	UUID              string    `json:"transaction_uuid"`
	Code              string    `json:"transaction_code"`
	Type              string    `json:"transaction_type"`
	Status            string    `json:"status"`
	Amount            float64   `json:"amount"`
	TransactionDate   time.Time `json:"transaction_date"`
	PaymentMethodName string    `json:"payment_method_name,omitempty"`
	OutletUUID        string    `json:"outlet_uuid,omitempty"`
	OutletName        string    `json:"outlet_name,omitempty"`
	CustomerName      string    `json:"customer_name,omitempty"`
	CustomerPhone     string    `json:"customer_phone,omitempty"`
}
type TransactionList struct {
	Data  []Transaction `json:"data"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}
type TransactionDetail map[string]any

type ProductInventoryStock struct {
	InventoryUUID string  `json:"inventory_uuid"`
	InventoryName string  `json:"inventory_name"`
	CurrentStock  float64 `json:"current_stock"`
}
type ProductVariant struct {
	UUID           string                  `json:"uuid"`
	Name           string                  `json:"name"`
	Price          float64                 `json:"price"`
	CurrentStock   float64                 `json:"current_stock"`
	BarcodeData    string                  `json:"barcode_data,omitempty"`
	MainImageToken string                  `json:"main_image_token,omitempty"`
	Inventories    []ProductInventoryStock `json:"inventories,omitempty"`
}
type ProductAddOn struct {
	UUID  string  `json:"uuid"`
	Name  string  `json:"add_on_name"`
	Price float64 `json:"price"`
}
type Product struct {
	UUID         string                  `json:"uuid"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description,omitempty"`
	Price        float64                 `json:"price"`
	Unit         string                  `json:"unit"`
	Type         string                  `json:"type"`
	CurrentStock float64                 `json:"current_stock"`
	TrackStock   bool                    `json:"track_stock"`
	Variants     []ProductVariant        `json:"variants"`
	AddOns       []ProductAddOn          `json:"add_ons,omitempty"`
	Inventories  []ProductInventoryStock `json:"inventories,omitempty"`
}
type ProductCategory struct {
	CategoryUUID string    `json:"category_uuid"`
	CategoryName string    `json:"category_name"`
	Products     []Product `json:"products"`
}

type AdjustStockRequest struct {
	ProductUUID        string  `json:"product_uuid,omitempty"`
	ProductVariantUUID string  `json:"product_variant_uuid,omitempty"`
	InventoryUUID      string  `json:"inventory_uuid"`
	Quantity           float64 `json:"quantity"`
	Notes              string  `json:"notes,omitempty"`
}
type AdjustStockResponse struct {
	ProductUUID    string  `json:"product_uuid"`
	ProductName    string  `json:"product_name"`
	ProductSKU     string  `json:"product_sku,omitempty"`
	VariantUUID    string  `json:"variant_uuid,omitempty"`
	VariantName    string  `json:"variant_name,omitempty"`
	InventoryUUID  string  `json:"inventory_uuid"`
	InventoryName  string  `json:"inventory_name"`
	QuantityBefore float64 `json:"quantity_before"`
	QuantityAfter  float64 `json:"quantity_after"`
	QuantityChange float64 `json:"quantity_change"`
	Unit           string  `json:"unit"`
}
type StockVariant struct {
	VariantUUID  string  `json:"variant_uuid"`
	VariantName  string  `json:"variant_name"`
	CurrentStock float64 `json:"current_stock"`
	Unit         string  `json:"unit"`
}
type StockProduct struct {
	ProductUUID  string         `json:"product_uuid"`
	ProductName  string         `json:"product_name"`
	ProductSKU   string         `json:"product_sku,omitempty"`
	CurrentStock float64        `json:"current_stock"`
	Unit         string         `json:"unit"`
	TrackStock   bool           `json:"track_stock"`
	Type         string         `json:"type"`
	Variants     []StockVariant `json:"variants,omitempty"`
}
type StockAggregate struct {
	InventoryUUID string         `json:"inventory_uuid"`
	InventoryName string         `json:"inventory_name"`
	Products      []StockProduct `json:"products"`
}
type StockMovement struct {
	UUID           string  `json:"uuid"`
	MovementType   string  `json:"movement_type"`
	QuantityChange float64 `json:"quantity_change"`
	Description    string  `json:"description,omitempty"`
	InventoryUUID  string  `json:"inventory_uuid,omitempty"`
	InventoryName  string  `json:"inventory_name,omitempty"`
	VariantUUID    string  `json:"variant_uuid,omitempty"`
	VariantName    string  `json:"variant_name,omitempty"`
	CreatedAt      string  `json:"created_at"`
}
type VoucherCodeRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
type VoucherCodeResponse struct {
	Code      string     `json:"code"`
	ExpiresAt *time.Time `json:"expired_at"`
}
