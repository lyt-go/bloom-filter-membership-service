package exportflow

import (
	"bloomfilter/internal/exportadapter"
	"bloomfilter/internal/exportrepo"
	"errors"
)

type Client interface{ Send() error }

func Export(repo *exportrepo.Repository, client Client) error {
	var last error
	for n := 0; n < 2; n++ {
		a := repo.Begin("snapshot")
		err := exportadapter.Normalize(client.Send())
		if err == nil {
			a.Commit()
			return nil
		}
		a.Rollback()
		last = err
		if errors.Is(err, exportadapter.ErrDenied) {
			return err
		}
	}
	return last
}
