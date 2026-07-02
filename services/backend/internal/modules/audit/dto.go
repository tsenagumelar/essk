package audit

import "time"

type LogResponse struct {
	ID           string         `json:"id"`
	TenantID     *string        `json:"tenant_id,omitempty"`
	ActorUserID  *string        `json:"actor_user_id,omitempty"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resource_type"`
	ResourceID   *string        `json:"resource_id,omitempty"`
	IPAddress    *string        `json:"ip_address,omitempty"`
	UserAgent    *string        `json:"user_agent,omitempty"`
	Metadata     map[string]any `json:"metadata"`
	CreatedDate  time.Time      `json:"created_date"`
}

type ListQuery struct {
	TenantID     *string
	ActorUserID  *string
	Action       string
	ResourceType string
	Page         int
	PageSize     int
}
