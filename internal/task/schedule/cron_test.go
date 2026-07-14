package schedule

import (
	"testing"
	"time"
)

func TestValidateCronExpr(t *testing.T) {
	if err := ValidateCronExpr("0 9 * * *"); err != nil {
		t.Fatalf("valid cron rejected: %v", err)
	}
	if err := ValidateCronExpr("not cron"); err == nil {
		t.Fatal("expected invalid cron error")
	}
}

func TestNextCronRun(t *testing.T) {
	from := time.Date(2026, 7, 14, 8, 0, 0, 0, time.UTC)
	next, err := NextCronRun("0 9 * * *", from)
	if err != nil {
		t.Fatal(err)
	}
	if next.Hour() != 9 {
		t.Fatalf("expected 9:00 next run, got %v", next)
	}
}