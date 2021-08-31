package thanosreceive

func (r *ThanosReceive) labels() map[string]string {
	ls := map[string]string{
		"app.kubernetes.io/component": componentName,
		"app.kubernetes.io/instance":  r.Instance.Name,
	}

	switch r.GetMode() {
	case RouterOnly:
		ls["thanos.receive/router"] = "true"
	case IngestorOnly:
		ls["thanos.receive/ingestor"] = "true"
	case RouterIngestor:
		ls["thanos.receive/router"] = "true"
		ls["thanos.receive/ingestor"] = "true"
	}

	return ls
}
