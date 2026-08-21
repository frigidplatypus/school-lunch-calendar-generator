package database

import (
	"testing"

	"calendar-app/internal/models"
)

func openTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestMigrateAndSeed(t *testing.T) {
	db := openTestDB(t)

	holidays, err := db.GetHolidays()
	if err != nil {
		t.Fatalf("GetHolidays: %v", err)
	}
	if len(holidays) == 0 {
		t.Fatal("expected seeded holidays, got none")
	}
}

func TestHolidayCRUD(t *testing.T) {
	db := openTestDB(t)

	id, err := db.AddHoliday(models.Holiday{Month: 2, Date: 17, Name: "St. Patrick's Day"})
	if err != nil {
		t.Fatalf("AddHoliday: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero id")
	}

	holidays, err := db.GetHolidays()
	if err != nil {
		t.Fatalf("GetHolidays: %v", err)
	}
	found := false
	for _, h := range holidays {
		if h.ID == id && h.Name == "St. Patrick's Day" {
			found = true
		}
	}
	if !found {
		t.Fatal("added holiday not found")
	}

	if err := db.DeleteHoliday(id); err != nil {
		t.Fatalf("DeleteHoliday: %v", err)
	}
	holidays, err = db.GetHolidays()
	if err != nil {
		t.Fatalf("GetHolidays: %v", err)
	}
	for _, h := range holidays {
		if h.ID == id {
			t.Fatal("deleted holiday still present")
		}
	}
}

func TestMenuItemCRUD(t *testing.T) {
	db := openTestDB(t)

	bid, err := db.AddMenuItem("breakfast", "Cereal")
	if err != nil {
		t.Fatalf("AddMenuItem: %v", err)
	}
	lid, err := db.AddMenuItem("lunch", "Pizza")
	if err != nil {
		t.Fatalf("AddMenuItem: %v", err)
	}

	items, err := db.GetMenuItems()
	if err != nil {
		t.Fatalf("GetMenuItems: %v", err)
	}
	if len(items.Breakfast) != 1 || items.Breakfast[0].ItemName != "Cereal" {
		t.Fatalf("unexpected breakfast items: %+v", items.Breakfast)
	}
	if len(items.Lunch) != 1 || items.Lunch[0].ItemName != "Pizza" {
		t.Fatalf("unexpected lunch items: %+v", items.Lunch)
	}

	if err := db.DeleteMenuItem(bid); err != nil {
		t.Fatalf("DeleteMenuItem: %v", err)
	}
	if err := db.DeleteMenuItem(lid); err != nil {
		t.Fatalf("DeleteMenuItem: %v", err)
	}
	items, err = db.GetMenuItems()
	if err != nil {
		t.Fatalf("GetMenuItems: %v", err)
	}
	if len(items.Breakfast) != 0 || len(items.Lunch) != 0 {
		t.Fatalf("expected empty items after delete: %+v", items)
	}
}

func TestCalendarDataUpsert(t *testing.T) {
	db := openTestDB(t)

	cd := models.CalendarData{
		Year: 2026, Month: 7,
		SchoolName: "Test School",
		DataJSON:   `{"menus":[]}`,
	}
	if err := db.SaveCalendarData(cd); err != nil {
		t.Fatalf("SaveCalendarData: %v", err)
	}

	got, err := db.GetCalendarData(2026, 7)
	if err != nil {
		t.Fatalf("GetCalendarData: %v", err)
	}
	if got == nil {
		t.Fatal("expected calendar data, got nil")
	}
	if got.SchoolName != "Test School" {
		t.Fatalf("unexpected school name: %s", got.SchoolName)
	}

	cd.SchoolName = "Updated School"
	cd.DataJSON = `{"menus":[{"date":10}]}`
	if err := db.SaveCalendarData(cd); err != nil {
		t.Fatalf("SaveCalendarData (upsert): %v", err)
	}

	all, err := db.GetAllCalendarData()
	if err != nil {
		t.Fatalf("GetAllCalendarData: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 row after upsert, got %d", len(all))
	}
	if all[0].SchoolName != "Updated School" {
		t.Fatalf("unexpected school name after upsert: %s", all[0].SchoolName)
	}
}

func TestGetCalendarDataMissing(t *testing.T) {
	db := openTestDB(t)

	got, err := db.GetCalendarData(1999, 0)
	if err != nil {
		t.Fatalf("GetCalendarData: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for missing data, got %+v", got)
	}
}

func TestExportImportRoundTrip(t *testing.T) {
	db := openTestDB(t)

	if _, err := db.AddMenuItem("breakfast", "Toast"); err != nil {
		t.Fatalf("AddMenuItem: %v", err)
	}
	if err := db.SaveCalendarData(models.CalendarData{Year: 2026, Month: 0, SchoolName: "S", DataJSON: `{"menus":[]}`}); err != nil {
		t.Fatalf("SaveCalendarData: %v", err)
	}

	exported, err := db.ExportAll()
	if err != nil {
		t.Fatalf("ExportAll: %v", err)
	}
	if len(exported.Holidays) == 0 || len(exported.MenuItems.Breakfast) != 1 || len(exported.CalendarData) != 1 {
		t.Fatalf("unexpected export contents")
	}

	other := openTestDB(t)
	if err := other.ImportAll(exported); err != nil {
		t.Fatalf("ImportAll: %v", err)
	}

	re, err := other.ExportAll()
	if err != nil {
		t.Fatalf("ExportAll (re): %v", err)
	}
	if len(re.Holidays) != len(exported.Holidays) {
		t.Fatalf("holiday count mismatch: %d vs %d", len(re.Holidays), len(exported.Holidays))
	}
	if len(re.MenuItems.Breakfast) != 1 || re.MenuItems.Breakfast[0].ItemName != "Toast" {
		t.Fatalf("menu items not imported: %+v", re.MenuItems.Breakfast)
	}
	if len(re.CalendarData) != 1 || re.CalendarData[0].SchoolName != "S" {
		t.Fatalf("calendar data not imported: %+v", re.CalendarData)
	}
}
