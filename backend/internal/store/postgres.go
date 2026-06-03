package store

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
)

type Postgres struct {
	databaseURL string
}

func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	p := &Postgres{databaseURL: databaseURL}
	if _, err := p.run(ctx, "SELECT 1"); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Postgres) Close() {}

func (p *Postgres) Migrate(ctx context.Context) error {
	_, err := p.runWithInput(ctx, schemaSQL)
	return err
}

func (p *Postgres) queryJSON(ctx context.Context, sql string) (string, error) {
	out, err := p.run(ctx, sql)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func (p *Postgres) run(ctx context.Context, sql string) (string, error) {
	return p.runPSQL(ctx, nil, "-c", sql)
}

func (p *Postgres) runWithInput(ctx context.Context, input string) (string, error) {
	return p.runPSQL(ctx, strings.NewReader(input))
}

func (p *Postgres) runPSQL(ctx context.Context, input *strings.Reader, args ...string) (string, error) {
	base := []string{p.databaseURL, "-v", "ON_ERROR_STOP=1", "-q", "-t", "-A"}
	cmd := exec.CommandContext(ctx, "psql", append(base, args...)...)
	if input != nil {
		cmd.Stdin = input
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", errors.New(strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

func sqlQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func isBlank(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}
