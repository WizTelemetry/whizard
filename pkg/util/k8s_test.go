package util

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestDecodeRawToContainers(t *testing.T) {
	containersBody := `
- args:
  - --listen-address=:8080
  - --reload-url=http://127.0.0.1:9090/-/reload
  - --config-file=/etc/prometheus/config/prometheus.yaml.gz
  - --config-envsubst-file=/etc/prometheus/config_out/prometheus.env.yaml
  - --watched-dir=/etc/prometheus/rules/prometheus-k8s-rulefiles-0
  command:
  - /bin/prometheus-config-reloader
  env:
  - name: POD_NAME
    valueFrom:
      fieldRef:
        apiVersion: v1
        fieldPath: metadata.name
  - name: SHARD
    value: "0"
  image: quay.io/prometheus-operator/prometheus-config-reloader:v0.68.0
  imagePullPolicy: IfNotPresent
  name: config-reloader
  ports:
  - containerPort: 8080
    name: reloader-web
    protocol: TCP
  resources:
    limits:
      cpu: 500m
      memory: 50Mi
    requests:
      cpu: 200m
      memory: 50Mi
  securityContext:
    allowPrivilegeEscalation: false
    capabilities:
      drop:
      - ALL
    readOnlyRootFilesystem: true
  terminationMessagePath: /dev/termination-log
  terminationMessagePolicy: FallbackToLogsOnError
  volumeMounts:
  - mountPath: /etc/prometheus/config
    name: config
  - mountPath: /etc/prometheus/config_out
    name: config-out
  - mountPath: /etc/prometheus/rules/prometheus-k8s-rulefiles-0
    name: prometheus-k8s-rulefiles-0
  - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
    name: kube-api-access-tpdrp
    readOnly: true
`

	containers, err := DecodeRawToContainers(runtime.RawExtension{Raw: []byte(containersBody)})
	if err != nil {
		t.Error(err)
	}
	if len(containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(containers))
	}
	if containers[0].Name != "config-reloader" {
		t.Errorf("Expected container name to be config-reloader, got %s", containers[0].Name)
	}
}
