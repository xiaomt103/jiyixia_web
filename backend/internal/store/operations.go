package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var validModules = map[string]bool{
	"booking": true, "order": true, "marketing": true, "finance": true, "report": true,
}

func (p *Postgres) ListOperations(ctx context.Context, module string) ([]OperationRecord, error) {
	where := ""
	if module != "" {
		value, ok := normalizeModule(module)
		if !ok {
			return nil, ErrInvalidInput
		}
		where = " WHERE module = " + sqlQuote(value)
	}
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`
		SELECT COALESCE(json_agg(row_to_json(t)), '[]'::json) FROM (
			SELECT id, module, title, member_id, amount_cents, status, due_at, created_at
			FROM operation_records%s ORDER BY due_at DESC, id DESC
		) t`, where))
	if err != nil {
		return nil, err
	}
	var records []OperationRecord
	return records, json.Unmarshal([]byte(payload), &records)
}

func (p *Postgres) CreateOperation(ctx context.Context, input CreateOperationInput) (OperationRecord, error) {
	input = trimOperation(input)
	module, ok := normalizeModule(input.Module)
	if !ok || isBlank(input.Title) || input.Amount < 0 {
		return OperationRecord{}, ErrInvalidInput
	}
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "pending"
	}
	dueAt, err := parseDueAt(input.DueAt)
	if err != nil {
		return OperationRecord{}, ErrInvalidInput
	}
	memberID := "NULL"
	if input.MemberID > 0 {
		memberID = fmt.Sprintf("%d", input.MemberID)
	}
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`
		WITH inserted AS (
			INSERT INTO operation_records (module, title, member_id, amount_cents, status, due_at)
			VALUES (%s, %s, %s, %d, %s, %s)
			RETURNING id, module, title, member_id, amount_cents, status, due_at, created_at
		) SELECT row_to_json(inserted) FROM inserted`,
		sqlQuote(module), sqlQuote(input.Title), memberID, input.Amount, sqlQuote(status), sqlQuote(dueAt.Format(time.RFC3339))))
	if err != nil {
		return OperationRecord{}, err
	}
	return decodeOperation(payload)
}

func (p *Postgres) UpdateOperationStatus(ctx context.Context, id int64, status string) (OperationRecord, error) {
	status = strings.TrimSpace(status)
	if id <= 0 || status == "" {
		return OperationRecord{}, ErrInvalidInput
	}
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`
		WITH updated AS (
			UPDATE operation_records SET status = %s WHERE id = %d
			RETURNING id, module, title, member_id, amount_cents, status, due_at, created_at
		) SELECT row_to_json(updated) FROM updated`, sqlQuote(status), id))
	if err != nil {
		return OperationRecord{}, err
	}
	return decodeOperation(payload)
}

func normalizeModule(value string) (string, bool) {
	value = strings.ToLower(strings.TrimSpace(value))
	return value, validModules[value]
}

func trimOperation(input CreateOperationInput) CreateOperationInput {
	input.Module = strings.TrimSpace(input.Module)
	input.Title = strings.TrimSpace(input.Title)
	input.Status = strings.TrimSpace(input.Status)
	input.DueAt = strings.TrimSpace(input.DueAt)
	return input
}

func parseDueAt(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Now().UTC(), nil
	}
	return time.Parse(time.RFC3339, value)
}

func decodeOperation(payload string) (OperationRecord, error) {
	if payload == "" {
		return OperationRecord{}, ErrNotFound
	}
	var record OperationRecord
	return record, json.Unmarshal([]byte(payload), &record)
}
