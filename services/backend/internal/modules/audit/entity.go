package audit

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID           uuid.UUID
	TenantID     *uuid.UUID
	ActorUserID  *uuid.UUID
	Action       string
	ResourceType string
	ResourceID   *string
	IPAddress    *string
	UserAgent    *string
	Metadata     json.RawMessage
	CreatedDate  time.Time
}

type Event struct {
	TenantID     *uuid.UUID
	ActorUserID  *uuid.UUID
	Action       string
	ResourceType string
	ResourceID   *string
	IPAddress    *string
	UserAgent    *string
	Metadata     map[string]any
}
