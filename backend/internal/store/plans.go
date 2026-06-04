package store

import (
	"context"
	"encoding/json"
	"fmt"
)

func (p *Postgres) ListPlans(ctx context.Context) ([]Plan, error) {
	payload, err := p.queryJSON(ctx, `
		SELECT COALESCE(json_agg(row_to_json(t)), '[]'::json) FROM (
			SELECT id, name, price_cents, duration_days, description, created_at
			FROM plans ORDER BY id
		) t`)
	if err != nil {
		return nil, err
	}
	var plans []Plan
	return plans, json.Unmarshal([]byte(payload), &plans)
}

func (p *Postgres) CreatePlan(ctx context.Context, input CreatePlanInput) (Plan, error) {
	input = trimPlanInput(input)
	if isBlank(input.Name) || input.PriceCents < 0 || input.DurationDay <= 0 {
		return Plan{}, ErrInvalidInput
	}
	payload, err := p.queryJSON(ctx, fmt.Sprintf(`
		WITH inserted AS (
			INSERT INTO plans (name, price_cents, duration_days, description)
			VALUES (%s, %d, %d, %s)
			RETURNING id, name, price_cents, duration_days, description, created_at
		)
		SELECT row_to_json(inserted) FROM inserted`,
		sqlQuote(input.Name), input.PriceCents, input.DurationDay, sqlQuote(input.Description)))
	if err != nil {
		return Plan{}, err
	}
	var plan Plan
	return plan, json.Unmarshal([]byte(payload), &plan)
}
