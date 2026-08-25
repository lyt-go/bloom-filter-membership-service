package hashdecode

type Decoder struct{ scratch []string }

func (d *Decoder) Decode(raw string) []string {
	d.scratch = append(d.scratch[:0], raw)
	return d.scratch
}
