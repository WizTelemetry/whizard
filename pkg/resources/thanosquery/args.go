package thanosquery

import (
	"fmt"
	"path/filepath"
)

func (q *ThanosQuery) args() []string {
	var args = []string{"query"}

	args = append(args, fmt.Sprintf("--grpc-address=[$(POD_IP)]:%d", grpcPort))
	args = append(args, fmt.Sprintf("--http-address=[$(POD_IP)]:%d", httpPort))
	args = append(args, "--store.sd-files="+filepath.Join(thanosConfigDir, storeSDFileName))
	args = append(args, "--store.sd-interval=2m")

	if kvs := q.Instance.Spec.SelectorLabels; len(kvs) > 0 {
		for k, v := range kvs {
			args = append(args, fmt.Sprintf(`--selector-label=%s="%s"`, k, v))
		}
	}

	if q.Instance.Spec.LogLevel != "" {
		args = append(args, "--log.level="+q.Instance.Spec.LogLevel)
	}
	if q.Instance.Spec.LogFormat != "" {
		args = append(args, "--log.format="+q.Instance.Spec.LogFormat)
	}
	return args
}
