package filterregistry

import "bloomfilter/internal/filterpolicy"

type Service struct{ Policy *filterpolicy.Policy }

func New(p *filterpolicy.Policy) *Service { return &Service{Policy: p} }
func (s *Service) Register(name string) error {
	if s.Policy.Validator != nil {
		if err := s.Policy.Validator.Validate(name); err != nil {
			return err
		}
	}
	s.Policy.Labels[name] = "ready"
	return nil
}
