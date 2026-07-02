package audit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestToResponseParsesMetadata(t *testing.T) {
	tenantID := uuid.New()
	actorID := uuid.New()
	resourceID := "resource-1"
	ip := "127.0.0.1"
	userAgent := "test-agent"
	metadata, err := json.Marshal(map[string]any{"status": "active"})
	if err != nil {
		t.Fatalf("marshal metadata: %v", err)
	}

	response := toResponse(Log{
		ID:           uuid.New(),
		TenantID:     &tenantID,
		ActorUserID:  &actorID,
		Action:       "tenant.update",
		ResourceType: "tenant",
		ResourceID:   &resourceID,
		IPAddress:    &ip,
		UserAgent:    &userAgent,
		Metadata:     metadata,
		CreatedDate:  time.Now().UTC(),
	})

	if response.TenantID == nil || *response.TenantID != tenantID.String() {
		t.Fatal("expected tenant id")
	}
	if response.ActorUserID == nil || *response.ActorUserID != actorID.String() {
		t.Fatal("expected actor user id")
	}
	if response.Metadata["status"] != "active" {
		t.Fatalf("expected metadata status active, got %v", response.Metadata["status"])
	}
}
