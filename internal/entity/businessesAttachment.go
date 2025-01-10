package entity

// BusinessAttachment defines the structure for the businesses_attachment table
type BusinessAttachment struct {
	ID          string `json:"id"`
	BusinessID  string `json:"business_id"`
	Filepath    string `json:"filepath"`
	ContentType string `json:"content_type"` // Should match the ENUM values: "photo", "video"
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
