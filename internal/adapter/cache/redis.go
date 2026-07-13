package cache

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// RedisCache implements the domain cache port using only the standard library.
// This avoids a second Redis dependency and keeps the binary small.
type RedisCache struct {
	address, password string
	db                int
	ttl, timeout      time.Duration
}

func NewRedis(address, password string, db int, ttl time.Duration) *RedisCache {
	if address == "" {
		return nil
	}
	return &RedisCache{address: address, password: password, db: db, ttl: ttl, timeout: 2 * time.Second}
}

func (r *RedisCache) Get(ctx context.Context, key string, target any) bool {
	raw, err := r.command(ctx, "GET", key)
	return err == nil && json.Unmarshal(raw, target) == nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value any) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = r.command(ctx, "SET", key, string(raw), "EX", strconv.FormatInt(max(1, int64(r.ttl.Seconds())), 10))
	return err
}

func (r *RedisCache) Ping(ctx context.Context) error { _, err := r.command(ctx, "PING"); return err }

func (r *RedisCache) command(ctx context.Context, args ...string) ([]byte, error) {
	conn, err := (&net.Dialer{Timeout: r.timeout}).DialContext(ctx, "tcp", r.address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(r.timeout))
	reader := bufio.NewReader(conn)
	if r.password != "" {
		if err := exchange(conn, reader, "AUTH", r.password); err != nil {
			return nil, err
		}
	}
	if r.db != 0 {
		if err := exchange(conn, reader, "SELECT", strconv.Itoa(r.db)); err != nil {
			return nil, err
		}
	}
	if err := writeCommand(conn, args...); err != nil {
		return nil, err
	}
	return readResponse(reader)
}

func exchange(w io.Writer, r *bufio.Reader, args ...string) error {
	if err := writeCommand(w, args...); err != nil {
		return err
	}
	_, err := readResponse(r)
	return err
}

func writeCommand(w io.Writer, args ...string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(args))
	for _, arg := range args {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(arg), arg)
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func readResponse(r *bufio.Reader) ([]byte, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
	switch prefix {
	case '+':
		return []byte(line), nil
	case '-':
		return nil, fmt.Errorf("redis: %s", line)
	case '$':
		n, err := strconv.Atoi(line)
		if err != nil || n < 0 {
			return nil, fmt.Errorf("redis: cache miss")
		}
		buf := make([]byte, n+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		return buf[:n], nil
	default:
		return nil, fmt.Errorf("redis: unsupported response")
	}
}
