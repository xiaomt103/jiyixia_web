package httpapi

import (
	"context"

	"jiyixia/member-system/internal/store"
)

type MemberStore interface {
	ListPlans(context.Context) ([]store.Plan, error)
	CreatePlan(context.Context, store.CreatePlanInput) (store.Plan, error)
	CreateMember(context.Context, store.CreateMemberInput) (store.Member, error)
	GetMember(context.Context, int64) (store.Member, error)
	ListMembers(context.Context, string) ([]store.Member, error)
	UpdateMemberStatus(context.Context, int64, string) (store.Member, error)
}

type SessionManager interface {
	Create(context.Context, string) (string, error)
	Validate(context.Context, string) (string, bool)
}
