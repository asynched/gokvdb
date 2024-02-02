package database

import "testing"

func TestDatabase(t *testing.T) {
	t.Run("New()", func(t *testing.T) {
		db := New()

		if db == nil {
			t.Error("Expected database")
		}
	})

	t.Run("Set()", func(t *testing.T) {
		db := New()
		db.Set("foo", Record{
			Value: "bar",
			Ttl:   0,
		})

		val, ok := db.Get("foo")

		if !ok {
			t.Error("Expected value to be set")
		}

		if val.Value != "bar" {
			t.Error("Expected value to be 'bar'")
		}
	})

	t.Run("Get()", func(t *testing.T) {
		db := New()
		db.Set("foo", Record{
			Value: "bar",
			Ttl:   0,
		})

		val, ok := db.Get("foo")

		if !ok {
			t.Error("Expected value to be set")
		}

		if val.Value != "bar" {
			t.Error("Expected value to be 'bar'")
		}
	})

	t.Run("Delete()", func(t *testing.T) {
		db := New()
		db.Set("foo", Record{
			Value: "bar",
			Ttl:   0,
		})

		db.Delete("foo")

		_, ok := db.Get("foo")

		if ok {
			t.Error("Expected value to be deleted")
		}
	})

	t.Run("FlushAll()", func(t *testing.T) {
		db := New()
		db.Set("foo", Record{
			Value: "bar",
			Ttl:   0,
		})

		db.FlushAll()

		_, ok := db.Get("foo")

		if ok {
			t.Error("Expected value to be deleted")
		}
	})

	t.Run("Snapshot()", func(t *testing.T) {
		db := New()
		db.Set("foo", Record{
			Value: "bar",
			Ttl:   0,
		})

		snapshot := db.Snapshot()

		if len(snapshot) != 1 {
			t.Error("Expected snapshot to have 1 record")
		}
	})
}
