package store

import "strings"

func normalizeStatus(status string) (string, bool) {
	value := strings.ToLower(strings.TrimSpace(status))
	return value, validStatuses[value]
}

func trimMemberInput(input CreateMemberInput) CreateMemberInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(input.Email)
	return input
}

func trimPlanInput(input CreatePlanInput) CreatePlanInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	return input
}
