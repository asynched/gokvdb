package commands

import (
	"testing"

	"github.com/asynched/gokvdb/internal/database"
)

func TestSetCommand(t *testing.T) {
	t.Run("Type()", func(t *testing.T) {
		c := &SetCommand{}
		if c.Type() != CommandTypeSet {
			t.Error("Expected CommandTypeSet")
		}
	})

	t.Run("Apply()", func(t *testing.T) {
		db := database.New()
		c := &SetCommand{
			Key:   "foo",
			Value: "bar",
			Ttl:   0,
		}
		_, err := c.Apply(db)

		if err != nil {
			t.Error("Expected nil error")
		}

		val, ok := db.Get("foo")

		if !ok {
			t.Error("Expected key to be found")
		}

		if val.Value != "bar" {
			t.Error("Expected value to be 'bar'")
		}
	})
}

func TestGetCommand(t *testing.T) {
	t.Run("Type()", func(t *testing.T) {
		c := &GetCommand{}
		if c.Type() != CommandTypeGet {
			t.Error("Expected CommandTypeGet")
		}
	})

	t.Run("Apply()", func(t *testing.T) {
		db := database.New()
		db.Set("foo", database.Record{
			Value: "bar",
			Ttl:   0,
		})

		c := &GetCommand{
			Key: "foo",
		}

		val, err := c.Apply(db)

		if err != nil {
			t.Error("Expected nil error")
		}

		if val != "bar" {
			t.Error("Expected value to be 'bar'")
		}
	})
}

func TestDeleteCommand(t *testing.T) {
	t.Run("Type()", func(t *testing.T) {
		c := &DeleteCommand{}
		if c.Type() != CommandTypeDelete {
			t.Error("Expected CommandTypeDelete")
		}
	})

	t.Run("Apply()", func(t *testing.T) {
		db := database.New()
		db.Set("foo", database.Record{
			Value: "bar",
			Ttl:   0,
		})

		c := &DeleteCommand{
			Key: "foo",
		}

		_, err := c.Apply(db)

		if err != nil {
			t.Error("Expected nil error")
		}

		_, ok := db.Get("foo")

		if ok {
			t.Error("Expected key to be deleted")
		}
	})
}

func TestFlushAll(t *testing.T) {
	t.Run("Type()", func(t *testing.T) {
		c := &FlushAllCommand{}
		if c.Type() != CommandTypeFlushAll {
			t.Error("Expected CommandTypeFlushAll")
		}
	})

	t.Run("Apply()", func(t *testing.T) {
		db := database.New()
		db.Set("foo", database.Record{
			Value: "bar",
			Ttl:   0,
		})

		c := &FlushAllCommand{}

		_, err := c.Apply(db)

		if err != nil {
			t.Error("Expected nil error")
		}

		snapshot := db.Snapshot()

		if len(snapshot) != 0 {
			t.Error("Expected database to be empty")
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("Parse()", func(t *testing.T) {
		cmd, err := Parse("SET foo bar 0")

		if err != nil {
			t.Error("Expected nil error")
		}

		if cmd.Type() != CommandTypeSet {
			t.Error("Expected CommandTypeSet")
		}
	})

	t.Run("Parse()", func(t *testing.T) {
		cmd, err := Parse("GET foo")

		if err != nil {
			t.Error("Expected nil error")
		}

		if cmd.Type() != CommandTypeGet {
			t.Error("Expected CommandTypeGet")
		}
	})

	t.Run("Parse()", func(t *testing.T) {
		cmd, err := Parse("DEL foo")

		if err != nil {
			t.Error("Expected nil error")
		}

		if cmd.Type() != CommandTypeDelete {
			t.Error("Expected CommandTypeDelete")
		}
	})

	t.Run("Parse()", func(t *testing.T) {
		cmd, err := Parse("FLUSHALL")

		if err != nil {
			t.Error("Expected nil error")
		}

		if cmd.Type() != CommandTypeFlushAll {
			t.Error("Expected CommandTypeFlushAll")
		}
	})

	t.Run("Parse()", func(t *testing.T) {
		_, err := Parse("INVALID")

		if err != ErrInvalidCommand {
			t.Error("Expected ErrInvalidCommand")
		}
	})
}
