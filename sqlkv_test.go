package sqlkv

import (
	"database/sql"
	"os"
	"testing"
	"time"
	
	_ "github.com/mattn/go-sqlite3"
)

func getStore() *SqlKv {
	var err error
	var db *sql.DB
	
	os.Mkdir("test", 0777)

	os.Remove("test/database.db")
	db, err = sql.Open("sqlite3", "test/database.db")
	if err != nil {
		panic(err)
	}
	
	return New(db, "kvstore")	
}

func clearStore(store *SqlKv) {
	store.db.Close()
	os.RemoveAll("test")
}

func panicHandler(t *testing.T, message string) {
	if r := recover(); r == nil {
		t.Error("%s: Expected call to panic, but it didn't", message)
	}
}

func noPanicHandler(t *testing.T, message string) {
	if r := recover(); r != nil {
		t.Error("%s: Expected call to not panic, but it did", message)
	}
}

func Test_createTable(t *testing.T) {
	store := getStore()
	
	store.db.Close()
		
	err := store.createTable()
	if err == nil {
		t.Error("Expected an error, got nil.")
	}
	
	clearStore(store)
	store = getStore()
	
	err = store.createTable()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	
	clearStore(store)
}

func Test_GetSetString(t *testing.T) {
	store := getStore()
	defer clearStore(store)
	
	store.SetString("test", "abcd")
	value := store.GetString("test")
	if value != "abcd" {
		t.Errorf("Expected 'abcd', got '%s'", value)
	}
	
	store.SetString("test", "1234")
	value = store.GetString("test")
	if value != "1234" {
		t.Errorf("Expected '1234', got '%s'", value)
	}

	store.db.Close()
	
	defer panicHandler(t, "GetString: database is closed")
	store.GetString("test")

	defer panicHandler(t, "SetString: database is closed")
	store.SetString("test", "panic")
}

func Test_GetSetInt(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	store.SetInt("test", 1234)
	i := store.GetInt("test")
	if i != 1234 {
		t.Errorf("Expected 1234, got %d", i)
	}
	
	store.SetString("test", "abcd")

	defer panicHandler(t, "GetInt: not a number")
	store.GetInt("test")
}

func Test_GetSetFloat(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	store.SetFloat("test", 1234.567)
	f := store.GetFloat("test")
	if f != 1234.567 {
		t.Errorf("Expected 1234.567, got %f", f)
	}
	
	store.SetString("test", "abcd")

	defer panicHandler(t, "GetFloat: not a number")
	store.GetFloat("test")
}

func TestGetSetBool(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	b := store.GetBool("nothere")
	if b {
		t.Error("Expected false, got true")
	}
	
	store.SetBool("test", true)
	if !store.GetBool("test") {
		t.Error("Expected true, got false")
	}
	
	store.SetBool("test", false)
	if store.GetBool("test") {
		t.Error("Expected false, got true")
	}
}

func Test_GetSetTime(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	v := store.GetTime("nothere")
	if !v.IsZero() {
		t.Errorf("Expected zero value, got %s", t)
	}

	now := time.Now()
	store.SetTime("test", now)
	v = store.GetTime("test")
	if v != now {
		t.Errorf("Expected %s, got %s", now, v)
	}
	
	store.SetString("test", "not a date")

	defer panicHandler(t, "GetTime: not a date")
	store.GetTime("test")
}

func Test_Del(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	defer noPanicHandler(t, "Del: deleting non-existant key")
	store.Del("blabla")
	
	store.SetString("test", "abcd")
	store.Del("test")
	
	value := store.GetString("test")
	if value != "" {
		t.Errorf("Expected '', got '%s'", value)
	}

	store.db.Close()
	defer panicHandler(t, "Del: database is closed")
	store.Del("test")
}

func Test_HasKey(t *testing.T) {
	store := getStore()
	defer clearStore(store)

	if store.HasKey("test") {
		t.Error("Expected false, got true")
	}
	
	store.SetString("test", "abcd")
	if !store.HasKey("test") {
		t.Error("Expected true, got false")		
	}

	store.SetString("test", "")
	if !store.HasKey("test") {
		t.Error("Expected true, got false")		
	}
	
	store.Del("test")
	if store.HasKey("test") {
		t.Error("Expected false, got true")		
	}
}