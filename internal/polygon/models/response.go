package models

type BaseResponse struct {
	PaginationHooks

	Status       string `json:"status,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	Count        int    `json:"count,omitempty"`
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

type PaginationHooks struct {
	NextURL string `json:"next_url,omitempty"`
}

func (p PaginationHooks) NextPage() string {
	return p.NextURL
}
