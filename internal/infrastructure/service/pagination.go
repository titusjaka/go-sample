package service

// Pagination is used as pagination struct for transport -> service -> storage communication
type Pagination struct {
	Limit       uint `json:"limit"`
	Offset      uint `json:"offset"`
	Total       uint `json:"total"`
	TotalPages  uint `json:"total_pages"`
	CurrentPage uint `json:"current_page"`
}
