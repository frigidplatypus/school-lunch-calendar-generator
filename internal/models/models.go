package models

type Holiday struct {
	ID    int    `json:"id"`
	Month int    `json:"month"`
	Date  int    `json:"date"`
	Name  string `json:"name"`
}

type MenuItem struct {
	ID       int    `json:"id"`
	MealType string `json:"meal_type"`
	ItemName string `json:"item_name"`
}

type CalendarData struct {
	ID         int    `json:"id"`
	Year       int    `json:"year"`
	Month      int    `json:"month"`
	SchoolName string `json:"school_name"`
	DataJSON   string `json:"data_json"`
	UpdatedAt  string `json:"updated_at"`
}

type LoginRequest struct {
	Password string `json:"password"`
}

type PasswordChangeRequest struct {
	NewPassword string `json:"new_password"`
}

type MenuItemsResponse struct {
	Breakfast []MenuItem `json:"breakfast"`
	Lunch     []MenuItem `json:"lunch"`
}

type ExportData struct {
	Holidays    []Holiday         `json:"holidays"`
	MenuItems   MenuItemsResponse `json:"menu_items"`
	CalendarData []CalendarData   `json:"calendar_data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
