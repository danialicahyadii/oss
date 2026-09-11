package response

type SubmissionResponse struct {
	UUID        string `json:"uuid"`
	BusinessID  string `json:"business_id"`
	SofficeID   uint64 `json:"soffice_id"`
	VariantID   uint64 `json:"variant_id"`
	Notes       string `json:"notes"`
	ValueType   string `json:"value_type"`
	CurrentRole string `json:"current_role"`
	Status      string `json:"status"`
	CreatedBy   uint64 `json:"created_by"`

	Customers []SubmissionCustomerResponse `json:"customers"`
	Materials []SubmissionMaterialResponse `json:"materials"`
}

type SubmissionCustomerResponse struct {
	UUID       string `json:"uuid"`
	CustomerID uint64 `json:"customer_id"`
}

type SubmissionMaterialResponse struct {
	UUID       string  `json:"uuid"`
	MaterialID uint64  `json:"material_id"`
	QtyJual    float64 `json:"qty_jual"`
	SalesUOM   string  `json:"sales_uom"`
	Value      float64 `json:"value"`
}
