package request

type CreateSubmissionRequest struct {
	SofficeID uint64 `json:"soffice_id" binding:"required"`
	VariantID uint64 `json:"variant_id" binding:"required"`
	Notes     string `json:"notes"`
	ValueType string `json:"value_type" binding:"required"`

	Customers []CreateSubmissionCustomerRequest `json:"customers" binding:"required,min=1"`
	Materials []CreateSubmissionMaterialRequest `json:"materials" binding:"required,min=1"`
}

type CreateSubmissionCustomerRequest struct {
	CustomerID uint64 `json:"customer_id" binding:"required"`
}

type CreateSubmissionMaterialRequest struct {
	MaterialID uint64  `json:"material_id" binding:"required"`
	QtyJual    float64 `json:"qty_jual" binding:"required"`
	SalesUOM   string  `json:"sales_uom" binding:"required"`
	Value      float64 `json:"value" binding:"required"`
}
