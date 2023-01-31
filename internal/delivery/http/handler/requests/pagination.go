package requests

type Pagination struct {
	Page         int    `json:"page"`
	FirstPageUrl string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	NextPageUrl  string `json:"next_page_url"`
	PrevPageUrl  string `json:"prev_page_url"`
	To           int    `json:"to"`
	Total        int    `json:"total"`
	Links        []Link `json:"links"`
}

type Link struct {
	Url    string `json:"url"`
	Label  string    `json:"label"`
	Active bool   `json:"active"`
	Page   int    `json:"page"`
}
