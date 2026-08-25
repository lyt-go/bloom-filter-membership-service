package exportflow

import (
	"bloomfilter/internal/exportadapter"
	"bloomfilter/internal/exportrepo"
	"errors"
	"testing"
)

type scriptedClient struct {
	errs  []error
	calls int
}

func (c *scriptedClient) Send() error { err := c.errs[c.calls]; c.calls++; return err }
func TestExportRetryPreservesErrorAndCommitBoundary(t *testing.T) {
	denied := &scriptedClient{errs: []error{exportadapter.ErrDenied, exportadapter.ErrDenied}}
	repo := &exportrepo.Repository{}
	err := Export(repo, denied)
	if !errors.Is(err, exportadapter.ErrDenied) {
		t.Fatalf("denied error lost: %v", err)
	}
	if denied.calls != 1 || len(repo.Records) != 0 {
		t.Fatalf("denied calls=%d records=%v", denied.calls, repo.Records)
	}
	temporary := &scriptedClient{errs: []error{errors.New("temporary"), nil}}
	repo = &exportrepo.Repository{}
	if err := Export(repo, temporary); err != nil {
		t.Fatalf("retry failed: %v", err)
	}
	if temporary.calls != 2 || len(repo.Records) != 1 {
		t.Fatalf("retry calls=%d records=%v", temporary.calls, repo.Records)
	}
}
