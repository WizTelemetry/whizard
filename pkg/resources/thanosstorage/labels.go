package thanosstorage

func (s *ThanosStorage) gatewayLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/component":    componentName,
		"app.kubernetes.io/instance":     s.Instance.Name,
		"app.kubernetes.io/subcomponent": "gateway",
	}
}

func (s *ThanosStorage) compactLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/component":    componentName,
		"app.kubernetes.io/instance":     s.Instance.Name,
		"app.kubernetes.io/subcomponent": "compact",
	}
}
