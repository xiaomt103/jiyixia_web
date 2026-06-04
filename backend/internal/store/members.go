package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

var validStatuses = map[string]bool{
	"pending": true, "active": true, "paused": true, "expired": true, "cancelled": true,
}

func (p *Postgres) CreateMember(ctx context.Context, input CreateMemberInput) (Member, error) {
	input = trimMemberInput(input)
	if isBlank(input.Name, input.Phone, input.Email) || input.PlanID <= 0 {
		return Member{}, ErrInvalidInput
	}
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`
		WITH inserted AS (
			INSERT INTO members (name, phone, email, plan_id)
			VALUES (%s, %s, %s, %d) RETURNING id
		)
		SELECT row_to_json(t) FROM (%s WHERE m.id = (SELECT id FROM inserted)) t`,
		sqlQuote(input.Name), sqlQuote(input.Phone), sqlQuote(input.Email), input.PlanID, memberSelect()))
	if err != nil {
		return Member{}, err
	}
	return decodeMember(payload)
}

func (p *Postgres) GetMember(ctx context.Context, id int64) (Member, error) {
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`SELECT row_to_json(t) FROM (%s WHERE m.id = %d) t`, memberSelect(), id))
	if err != nil {
		return Member{}, err
	}
	return decodeMember(payload)
}

func (p *Postgres) ListMembers(ctx context.Context, status string) ([]Member, error) {
	where := ""
	if status != "" {
		value, ok := normalizeStatus(status)
		if !ok {
			return nil, ErrInvalidInput
		}
		where = " WHERE m.status = " + sqlQuote(value)
	}
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`
		SELECT COALESCE(json_agg(row_to_json(t)), '[]'::json)
		FROM (%s%s ORDER BY m.id DESC) t`, memberSelect(), where))
	if err != nil {
		return nil, err
	}
	var members []Member
	return members, json.Unmarshal([]byte(payload), &members)
}

func (p *Postgres) UpdateMemberStatus(ctx context.Context, id int64, status string) (Member, error) {
	value, ok := normalizeStatus(status)
	if !ok {
		return Member{}, ErrInvalidInput
	}
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`
		WITH updated AS (
			UPDATE members SET status = %s,
			start_at = CASE WHEN %s = 'active' AND start_at IS NULL THEN now() ELSE start_at END,
			expire_at = CASE WHEN %s = 'active' THEN now() + (plans.duration_days || ' days')::interval ELSE expire_at END,
			updated_at = now() FROM plans WHERE members.plan_id = plans.id AND members.id = %d RETURNING members.id
		)
		SELECT row_to_json(t) FROM (%s WHERE m.id = (SELECT id FROM updated)) t`,
		sqlQuote(value), sqlQuote(value), sqlQuote(value), id, memberSelect()))
	if err != nil {
		return Member{}, err
	}
	return decodeMember(payload)
}

func memberSelect() string {
	return `SELECT m.id, m.name, m.phone, m.email, m.plan_id, p.name AS plan_name, m.status,
		m.start_at, m.expire_at, m.created_at, m.updated_at FROM members m JOIN plans p ON p.id = m.plan_id`
}

func decodeMember(payload string) (Member, error) {
	if payload == "" {
		return Member{}, ErrNotFound
	}
	var member Member
	if err := json.Unmarshal([]byte(payload), &member); err != nil {
		return Member{}, err
	}
	if member.ExpireAt != nil && member.ExpireAt.Before(time.Now()) && member.Status == "active" {
		member.Status = "expired"
	}
	return member, nil
}
