package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/asynched/gokvdb/internal/database"
)

type CommandType int

const (
	CommandTypeSet CommandType = iota
	CommandTypeGet
	CommandTypeDelete
	CommandTypeFlushAll
)

type CommandError struct {
	Message string
}

func (e CommandError) Error() string {
	return e.Message
}

var (
	ErrInvalidCommand error = CommandError{"ERR_INVALID_COMMAND"}
	ErrKeyNotFoud           = CommandError{"ERR_KEY_NOT_FOUND"}
	ErrInvalidValue         = CommandError{"ERR_INVALID_VALUE"}
)

type Command interface {
	Type() CommandType
	Apply(*database.Database) (string, error)
}

type SetCommand struct {
	Key   string
	Value string
	Ttl   int64
}

func (c *SetCommand) Type() CommandType {
	return CommandTypeSet
}

func (c *SetCommand) Apply(db *database.Database) (string, error) {
	db.Set(c.Key, database.Record{
		Value: c.Value,
		Ttl:   c.Ttl,
	})

	return "OK", nil
}

type GetCommand struct {
	Key string
}

func (c *GetCommand) Type() CommandType {
	return CommandTypeGet
}

func (c *GetCommand) Apply(db *database.Database) (string, error) {
	record, ok := db.Get(c.Key)

	if !ok {
		return "", ErrKeyNotFoud
	}

	return record.Value, nil
}

type DeleteCommand struct {
	Key string
}

func (c *DeleteCommand) Type() CommandType {
	return CommandTypeDelete
}

func (c *DeleteCommand) Apply(db *database.Database) (string, error) {
	_, ok := db.Get(c.Key)

	if !ok {
		return "", ErrKeyNotFoud
	}

	db.Delete(c.Key)

	return "OK", nil
}

type FlushAllCommand struct{}

func (c *FlushAllCommand) Type() CommandType {
	return CommandTypeFlushAll
}

func (c *FlushAllCommand) Apply(db *database.Database) (string, error) {
	db.FlushAll()
	return "OK", nil
}

func Parse(cmd string) (Command, error) {
	if strings.HasPrefix(cmd, "GET") {
		parts := strings.Fields(cmd)

		if len(parts) != 2 {
			return nil, ErrInvalidCommand
		}

		return &GetCommand{Key: parts[1]}, nil
	}

	if strings.HasPrefix(cmd, "SET") {
		parts := strings.Fields(cmd)

		if len(parts) < 4 {
			return nil, ErrInvalidCommand
		}

		value := strings.Join(parts[2:len(parts)-1], " ")
		ttl, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)

		if err != nil {
			return nil, ErrInvalidValue
		}

		if ttl <= 0 {
			return &SetCommand{Key: parts[1], Value: value, Ttl: 0}, nil
		}

		exp := time.Now().Add(time.Duration(ttl) * time.Second).Unix()

		return &SetCommand{Key: parts[1], Value: value, Ttl: exp}, nil
	}

	if strings.HasPrefix(cmd, "DEL") {
		parts := strings.Fields(cmd)

		if len(parts) != 2 {
			return nil, ErrInvalidCommand
		}

		return &DeleteCommand{Key: parts[1]}, nil
	}

	if strings.HasPrefix(cmd, "FLUSHALL") {
		return &FlushAllCommand{}, nil
	}

	return nil, ErrInvalidCommand
}
