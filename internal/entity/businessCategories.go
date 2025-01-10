package entity

// BusinessCategory defines the structure for the business_categories table
type BusinessCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type BusinessCategoryList struct {
	Items []BusinessCategory `json:"businesses_categories"`
	Count int                `json:"count"`
}
