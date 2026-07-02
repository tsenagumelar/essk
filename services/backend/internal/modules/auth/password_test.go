package auth

import "testing"

func TestPasswordHasher(t *testing.T) {
	hasher := NewPasswordHasher()

	hash, err := hasher.Hash("Admin123!")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	ok, err := hasher.Verify("Admin123!", hash)
	if err != nil {
		t.Fatalf("verify password: %v", err)
	}
	if !ok {
		t.Fatal("expected password verification to pass")
	}

	ok, err = hasher.Verify("wrong-password", hash)
	if err != nil {
		t.Fatalf("verify wrong password: %v", err)
	}
	if ok {
		t.Fatal("expected wrong password verification to fail")
	}
}

func TestRefreshTokenHash(t *testing.T) {
	raw, hash, err := NewRefreshToken()
	if err != nil {
		t.Fatalf("new refresh token: %v", err)
	}
	if raw == "" {
		t.Fatal("expected raw token")
	}
	if hash == "" {
		t.Fatal("expected token hash")
	}
	if HashRefreshToken(raw) != hash {
		t.Fatal("expected deterministic refresh token hash")
	}
}
