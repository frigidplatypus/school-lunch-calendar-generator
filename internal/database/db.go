package database

import (
	"database/sql"
	"fmt"

	"calendar-app/internal/models"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func Open(path string) (*DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	d := &DB{db}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return d, nil
}

func (d *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS holidays (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			month INTEGER NOT NULL,
			date INTEGER NOT NULL,
			name TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS menu_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			meal_type TEXT NOT NULL CHECK(meal_type IN ('breakfast', 'lunch')),
			item_name TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS calendar_data (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			year INTEGER NOT NULL,
			month INTEGER NOT NULL,
			school_name TEXT DEFAULT '',
			data_json TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(year, month)
		)`,
	}

	for _, q := range queries {
		if _, err := d.Exec(q); err != nil {
			return err
		}
	}

	return d.seedDefaults()
}

func (d *DB) seedDefaults() error {
	var count int
	if err := d.QueryRow("SELECT COUNT(*) FROM holidays").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	defaults := []models.Holiday{
		{Month: 0, Date: 1, Name: "New Year's Day"},
		{Month: 0, Date: 14, Name: "MLK Day"},
		{Month: 1, Date: 14, Name: "Valentine's Day"},
		{Month: 4, Date: 26, Name: "Memorial Day"},
		{Month: 6, Date: 4, Name: "Independence Day"},
		{Month: 8, Date: 1, Name: "Labor Day"},
		{Month: 9, Date: 14, Name: "Columbus Day"},
		{Month: 9, Date: 31, Name: "Halloween"},
		{Month: 10, Date: 11, Name: "Veterans Day"},
		{Month: 10, Date: 27, Name: "Thanksgiving"},
		{Month: 10, Date: 28, Name: "Thanksgiving Break"},
		{Month: 11, Date: 24, Name: "Christmas Eve"},
		{Month: 11, Date: 25, Name: "Christmas Day"},
		{Month: 11, Date: 31, Name: "New Year's Eve"},
	}

	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO holidays (month, date, name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, h := range defaults {
		if _, err := stmt.Exec(h.Month, h.Date, h.Name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *DB) GetHolidays() ([]models.Holiday, error) {
	rows, err := d.Query("SELECT id, month, date, name FROM holidays ORDER BY month, date")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []models.Holiday
	for rows.Next() {
		var h models.Holiday
		if err := rows.Scan(&h.ID, &h.Month, &h.Date, &h.Name); err != nil {
			return nil, err
		}
		holidays = append(holidays, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return holidays, nil
}

func (d *DB) AddHoliday(h models.Holiday) (int, error) {
	res, err := d.Exec("INSERT INTO holidays (month, date, name) VALUES (?, ?, ?)", h.Month, h.Date, h.Name)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (d *DB) DeleteHoliday(id int) error {
	_, err := d.Exec("DELETE FROM holidays WHERE id = ?", id)
	return err
}

func (d *DB) GetMenuItems() (*models.MenuItemsResponse, error) {
	rows, err := d.Query("SELECT id, meal_type, item_name FROM menu_items ORDER BY item_name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.MenuItemsResponse{
		Breakfast: []models.MenuItem{},
		Lunch:     []models.MenuItem{},
	}

	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.ID, &item.MealType, &item.ItemName); err != nil {
			return nil, err
		}
		switch item.MealType {
		case "breakfast":
			resp.Breakfast = append(resp.Breakfast, item)
		case "lunch":
			resp.Lunch = append(resp.Lunch, item)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DB) AddMenuItem(mealType, itemName string) (int, error) {
	res, err := d.Exec("INSERT INTO menu_items (meal_type, item_name) VALUES (?, ?)", mealType, itemName)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (d *DB) DeleteMenuItem(id int) error {
	_, err := d.Exec("DELETE FROM menu_items WHERE id = ?", id)
	return err
}

func (d *DB) GetCalendarData(year, month int) (*models.CalendarData, error) {
	var cd models.CalendarData
	err := d.QueryRow(
		"SELECT id, year, month, school_name, data_json, updated_at FROM calendar_data WHERE year = ? AND month = ?",
		year, month,
	).Scan(&cd.ID, &cd.Year, &cd.Month, &cd.SchoolName, &cd.DataJSON, &cd.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cd, nil
}

func (d *DB) SaveCalendarData(cd models.CalendarData) error {
	_, err := d.Exec(
		`INSERT INTO calendar_data (year, month, school_name, data_json, updated_at)
		 VALUES (?, ?, ?, ?, datetime('now'))
		 ON CONFLICT(year, month) DO UPDATE SET
		   school_name = excluded.school_name,
		   data_json = excluded.data_json,
		   updated_at = excluded.updated_at`,
		cd.Year, cd.Month, cd.SchoolName, cd.DataJSON,
	)
	return err
}

func (d *DB) GetAllCalendarData() ([]models.CalendarData, error) {
	rows, err := d.Query("SELECT id, year, month, school_name, data_json, updated_at FROM calendar_data ORDER BY year, month")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.CalendarData
	for rows.Next() {
		var cd models.CalendarData
		if err := rows.Scan(&cd.ID, &cd.Year, &cd.Month, &cd.SchoolName, &cd.DataJSON, &cd.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, cd)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DB) ExportAll() (*models.ExportData, error) {
	holidays, err := d.GetHolidays()
	if err != nil {
		return nil, err
	}
	if holidays == nil {
		holidays = []models.Holiday{}
	}

	menuItems, err := d.GetMenuItems()
	if err != nil {
		return nil, err
	}

	calData, err := d.GetAllCalendarData()
	if err != nil {
		return nil, err
	}
	if calData == nil {
		calData = []models.CalendarData{}
	}

	return &models.ExportData{
		Holidays:     holidays,
		MenuItems:    *menuItems,
		CalendarData: calData,
	}, nil
}

func (d *DB) ImportAll(data *models.ExportData) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM holidays"); err != nil {
		return err
	}
	for _, h := range data.Holidays {
		if _, err := tx.Exec("INSERT INTO holidays (month, date, name) VALUES (?, ?, ?)", h.Month, h.Date, h.Name); err != nil {
			return err
		}
	}

	if _, err := tx.Exec("DELETE FROM menu_items"); err != nil {
		return err
	}
	for _, item := range data.MenuItems.Breakfast {
		if _, err := tx.Exec("INSERT INTO menu_items (meal_type, item_name) VALUES ('breakfast', ?)", item.ItemName); err != nil {
			return err
		}
	}
	for _, item := range data.MenuItems.Lunch {
		if _, err := tx.Exec("INSERT INTO menu_items (meal_type, item_name) VALUES ('lunch', ?)", item.ItemName); err != nil {
			return err
		}
	}

	if _, err := tx.Exec("DELETE FROM calendar_data"); err != nil {
		return err
	}
	for _, cd := range data.CalendarData {
		if _, err := tx.Exec(
			"INSERT INTO calendar_data (year, month, school_name, data_json) VALUES (?, ?, ?, ?)",
			cd.Year, cd.Month, cd.SchoolName, cd.DataJSON,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *DB) ImportHolidays(holidays []models.Holiday) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM holidays"); err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO holidays (month, date, name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, holiday := range holidays {
		if _, err := stmt.Exec(holiday.Month, holiday.Date, holiday.Name); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) ImportMenuItems(items []models.MenuItem) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM menu_items"); err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO menu_items (meal_type, item_name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, item := range items {
		if _, err := stmt.Exec(item.MealType, item.ItemName); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) ImportCalendarData(calendarData []models.CalendarData) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM calendar_data"); err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO calendar_data (year, month, school_name, data_json) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, data := range calendarData {
		if _, err := stmt.Exec(data.Year, data.Month, data.SchoolName, data.DataJSON); err != nil {
			return err
		}
	}
	return tx.Commit()
}
