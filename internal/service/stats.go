package service

// Overview 返回全局统计总览：每个实体的总数与按状态分布。
func (s *Service) Overview() map[string]interface{} {
	return map[string]interface{}{
		"filters":    s.countFilter(),
		"strategies": s.countStrategy(),
		"elements":   s.countElement(),
		"probes":     s.countProbe(),
		"groups":     s.countGroup(),
		"tags":       s.countTag(),
	}
}

func (s *Service) countFilter() map[string]interface{} {
	items := s.store.ListFilters()
	m := map[string]int{}
	for _, it := range items {
		m[it.Status]++
	}
	return map[string]interface{}{"total": len(items), "by_status": m}
}

func (s *Service) countStrategy() map[string]interface{} {
	items := s.store.ListStrategies()
	m := map[string]int{}
	for _, it := range items {
		m[it.Status]++
	}
	return map[string]interface{}{"total": len(items), "by_status": m}
}

func (s *Service) countElement() map[string]interface{} {
	items := s.store.ListElements()
	m := map[string]int{}
	for _, it := range items {
		m[it.Category]++
	}
	return map[string]interface{}{"total": len(items), "by_category": m}
}

func (s *Service) countProbe() map[string]interface{} {
	items := s.store.ListProbes()
	m := map[string]int{}
	fp := 0
	for _, it := range items {
		m[it.Result]++
		if it.FalsePositive {
			fp++
		}
	}
	return map[string]interface{}{"total": len(items), "by_result": m, "false_positive": fp}
}

func (s *Service) countGroup() map[string]interface{} {
	items := s.store.ListGroups()
	return map[string]interface{}{"total": len(items)}
}

func (s *Service) countTag() map[string]interface{} {
	items := s.store.ListTags()
	return map[string]interface{}{"total": len(items)}
}
