# gokvdb

A key-value database similar to memcached written in Go.

## Installation

### Requirements

- Go
- GNU/Make

### Build

Build the server with:

```sh
make server
```

> This will generate a `bin/server` executable

Build the client with:

```sh
make client
```

> This will generate a `bin/client` executable

### Commands

#### SET

Sets a new entry in the database with a key, value and time to live in seconds (ttl).

```sh
# SET key value ttl
?> SET name asynched 10
```

#### GET

Retrieves the value of a key from the database.

```sh
# GET key
?> GET name
```

#### DELETE

Deletes a key from the database.

```sh
# DELETE key
DELETE name
```

#### FLUSHALL

Deletes all keys from the database.

```sh
# FLUSHALL
FLUSHALL
```
