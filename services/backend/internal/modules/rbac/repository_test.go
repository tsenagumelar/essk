package rbac

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRoleAppliesToUserTenant(t *testing.T) {
	userTenant := uuid.New()
	otherTenant := uuid.New()

	if !RoleAppliesToUserTenant(&userTenant, nil) {
		t.Fatal("expected system role to apply to tenant user")
	}
	if !RoleAppliesToUserTenant(&userTenant, &userTenant) {
		t.Fatal("expected same-tenant role to apply")
	}
	if RoleAppliesToUserTenant(&userTenant, &otherTenant) {
		t.Fatal("expected cross-tenant role to be denied")
	}
	if RoleAppliesToUserTenant(nil, &otherTenant) {
		t.Fatal("expected tenant-scoped role to be denied for user without tenant")
	}
}

func TestUserHasPermissionSQLContainsTenantBoundary(t *testing.T) {
	requiredFragments := []string{
		"FROM users u",
		"JOIN user_roles ur ON ur.user_id = u.id",
		"ro.tenant_id IS NULL OR ro.tenant_id = u.tenant_id",
		"u.is_deleted = false",
		"ro.is_deleted = false",
	}

	for _, fragment := range requiredFragments {
		if !strings.Contains(userHasPermissionSQL, fragment) {
			t.Fatalf("expected permission SQL to contain tenant boundary fragment %q", fragment)
		}
	}
}
