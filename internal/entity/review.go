package entity

// Review represents the reviewes table
type Review struct {
	ID          string             `json:"id"`
	BusinessID  string             `json:"business_id"`
	UserID      string             `json:"user_id"`
	Rating      uint8              `json:"rating"`
	Comment     string             `json:"comment"`
	Attachments []ReviewAttachment `json:"attachments"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
}

type ReviewList struct {
	Items []Review `json:"reviewes"`
	Count int      `json:"count"`
}

type ReviewSingleRequest struct {
	ID         string `json:"id"`
	OwnerID    string `json:"owner_id"`
	CategoryID string `json:"category_id"`
}

// ReviewAttachmentList defines the structure for the list of review attachments
type ReviewAttachment struct {
	Id          string `json:"id"`
	ReviewId    string `json:"-"`
	FilePath    string `json:"filepath"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ReviewAttachmentList struct {
	Items []ReviewAttachment `json:"items"`
	Count int64              `json:"count"`
}

type ReviewAttachmentMultipleInsertRequest struct {
	ReviewId    string             `json:"review_id"`
	Attachments []ReviewAttachment `json:"attachments"`
}
