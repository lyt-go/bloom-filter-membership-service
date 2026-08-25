package exportrepo

type Repository struct{ Records []string }
type Attempt struct {
	repo  *Repository
	value string
}

func (r *Repository) Begin(value string) *Attempt {
	return &Attempt{repo: r, value: value}
}
func (a *Attempt) Commit() {
	a.repo.Records = append(a.repo.Records, a.value)
}
func (a *Attempt) Rollback() {}
