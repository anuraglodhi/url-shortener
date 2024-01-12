package database

import (
	"testing"
)

func TestShortUrlExists(t *testing.T) {
    db := New()

    exists, err := db.ShortUrlExists("pb0zf6a8q")
    if err != nil {
        t.Fatalf("DB error: %s", err)
    }
    if !exists {
        t.Fatalf("Expected short URL to exist, but it doesn't")
    }

    exists, err = db.ShortUrlExists("nonexistent")
    if err != nil {
        t.Fatalf("DB error: %s", err)
    }
    if exists {
        t.Fatalf("Expected short URL to not exist, but it does")
    }

		db.Close()
}

func TestLongUrlExists(t *testing.T) {
		db := New()

		exists, err := db.LongUrlExists("www.google.com")
		if err != nil {
				t.Fatalf("DB error: %s", err)
		}
		if !exists {
				t.Fatalf("Expected long URL to exist, but it doesn't")
		}

		exists, err = db.LongUrlExists("nonexistent")
		if err != nil {
				t.Fatalf("DB error: %s", err)
		}
		if exists {
				t.Fatalf("Expected long URL to not exist, but it does")
		}

		db.Close()
}

func TestGetLongUrl(t *testing.T) {
		db := New()

		longUrl, err := db.GetLongUrl("m805f8faq")
		if err != nil {
				t.Fatalf("DB error: %s", err)
		}
		if longUrl != "www.google.com" {
				t.Fatalf("Expected long URL to be www.google.com, but it is %s", longUrl)
		}

		longUrl, err = db.GetLongUrl("nonexistent")
		if err == nil && longUrl != "" {
				t.Fatalf("Expected long URL to be empty, but it is %s", longUrl)
		}

		db.Close()
}