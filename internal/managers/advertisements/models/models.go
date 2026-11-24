package models

type Photo struct {
	ID              string `json:"id"`
	AdvertisementID string `json:"advertisement_id"`
	URL             string `json:"url"`
	SortOrder       int    `json:"sort_order"`
	CreatedAt       string `json:"created_at"`
}

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserID      string  `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Photos      []Photo `json:"photos"`
}

type AdvertisementsList struct {
	Items      []Advertisement `json:"items"`
	HasNext    bool            `json:"has_next"`
	NextCursor string          `json:"next_cursor"`
	Total      int             `json:"total"`
}

// CreateAdvertisementRequest параметры формы создания объявления.
// PriceRaw и QuantityRaw нужны для негативных тестов: если заполнены, в форму уходят они вместо Price и Quantity
// (например, строка вместо числа). PhotoPaths - пути к файлам фотографий относительно папки с тестом
type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
	PriceRaw    string
	QuantityRaw string
	PhotoPaths  []string
}
