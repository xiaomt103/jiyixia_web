package auth

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SessionStore struct {
	addr string
	ttl  time.Duration
}

func NewRedisSessionStore(redisURL string, ttl time.Duration) *SessionStore {
	return &SessionStore{addr: redisAddr(redisURL), ttl: ttl}
}

func (s *SessionStore) Ping(ctx context.Context) error {
	_, err := s.command(ctx, "PING")
	return err
}

func (s *SessionStore) Close() error { return nil }

func (s *SessionStore) Create(ctx context.Context, subject string) (string, error) {
	token, err := randomToken()
	if err != nil {
		return "", err
	}
	_, err = s.command(ctx, "SET", key(token), subject, "EX", strconv.Itoa(int(s.ttl.Seconds())))
	return token, err
}

func (s *SessionStore) Validate(ctx context.Context, token string) (string, bool) {
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}
	subject, err := s.command(ctx, "GET", key(token))
	return subject, err == nil && subject != ""
}

func (s *SessionStore) command(ctx context.Context, parts ...string) (string, error) {
	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", s.addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}
	if _, err := fmt.Fprint(conn, encodeRESP(parts)); err != nil {
		return "", err
	}
	return readRESP(bufio.NewReader(conn))
}

func encodeRESP(parts []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(parts))
	for _, part := range parts {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(part), part)
	}
	return b.String()
}

func readRESP(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(line, "\r\n")
	if line == "" {
		return "", errors.New("empty redis response")
	}
	switch line[0] {
	case '+':
		return line[1:], nil
	case '-':
		return "", errors.New(line[1:])
	case '$':
		return readBulkString(r, line[1:])
	default:
		return "", errors.New("unsupported redis response")
	}
}

func readBulkString(r *bufio.Reader, sizeText string) (string, error) {
	size, err := strconv.Atoi(sizeText)
	if err != nil {
		return "", err
	}
	if size < 0 {
		return "", errors.New("redis nil bulk string")
	}
	buf := make([]byte, size+2)
	if _, err := r.Read(buf); err != nil {
		return "", err
	}
	return string(buf[:size]), nil
}

func redisAddr(redisURL string) string {
	parsed, err := url.Parse(redisURL)
	if err == nil && parsed.Host != "" {
		return parsed.Host
	}
	return "localhost:6379"
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func key(token string) string {
	return "admin_session:" + token
}
