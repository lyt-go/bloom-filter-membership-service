package exportrepo

type Repository struct{ Records []string }
type Attempt struct {
	repo  *Repository
	value string
}

func (r *Repository) Begin(value string) *Attempt {
	r.Records = append(r.Records, value)
	return &Attempt{repo: r, value: value}
}
func (a *Attempt) Commit()   {}
func (a *Attempt) Rollback() {}
