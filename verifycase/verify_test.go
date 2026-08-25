package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"buscharge/internal/audit"
	"buscharge/internal/store"
)

func TestBcAuditWriteSwallow(t *testing.T) {
	dir := t.TempDir()
	st, err := store.New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "audit"), []byte("block"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc := audit.NewService(st)
	if _, err := svc.Record("depot", "charge", "b1", "p1", 1.0); err == nil {
		t.Fatal("record must propagate the store write error")
	}
}
