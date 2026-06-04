package store

import "time"

type Plan struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	PriceCents  int       `json:"price_cents"`
	DurationDay int       `json:"duration_days"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Member struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	PlanID    int64      `json:"plan_id"`
	PlanName  string     `json:"plan_name"`
	Status    string     `json:"status"`
	StartAt   *time.Time `json:"start_at,omitempty"`
	ExpireAt  *time.Time `json:"expire_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateMemberInput struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	PlanID int64  `json:"plan_id"`
}

type CreatePlanInput struct {
	Name        string `json:"name"`
	PriceCents  int    `json:"price_cents"`
	DurationDay int    `json:"duration_days"`
	Description string `json:"description"`
}

type OperationRecord struct {
	ID        int64     `json:"id"`
	Module    string    `json:"module"`
	Title     string    `json:"title"`
	MemberID  int64     `json:"member_id,omitempty"`
	Amount    int       `json:"amount_cents"`
	Status    string    `json:"status"`
	DueAt     time.Time `json:"due_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateOperationInput struct {
	Module   string `json:"module"`
	Title    string `json:"title"`
	MemberID int64  `json:"member_id"`
	Amount   int    `json:"amount_cents"`
	Status   string `json:"status"`
	DueAt    string `json:"due_at"`
}
