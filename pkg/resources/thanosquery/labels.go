package thanosquery

func (q *ThanosQuery) labels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/component": componentName,
		"app.kubernetes.io/instance":  q.Instance.Name,
	}
}
