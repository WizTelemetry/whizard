package thanosreceive

import (
	"fmt"
	"path/filepath"
)

func (r *ThanosReceive) args() []string {
	var args = []string{
		"receive",
		fmt.Sprintf("--grpc-address=[$(POD_IP)]:%d", grpcPort),
		fmt.Sprintf("--http-address=[$(POD_IP)]:%d", httpPort),
		fmt.Sprintf("--remote-write.address=[$(POD_IP)]:%d", remoteWritePort),
		fmt.Sprintf(`--label=receive_replica="$(NAME)"`),
		fmt.Sprintf(`--tsdb.path=%s`, storageDir),
	}

	if r.Instance.Spec.LogLevel != "" {
		args = append(args, "--log.level="+r.Instance.Spec.LogLevel)
	}
	if r.Instance.Spec.LogFormat != "" {
		args = append(args, "--log.format="+r.Instance.Spec.LogFormat)
	}

	var m = r.GetMode()

	if m == IngestorOnly || m == RouterIngestor {
		args = append(args, fmt.Sprintf("--receive.local-endpoint=[$(POD_IP)]:%d", remoteWritePort))
		if ingestorSpec := r.Instance.Spec.Ingestor; ingestorSpec != nil {
			if ingestorSpec.LocalTSDBRetention != "" {
				args = append(args, "--tsdb.retention="+r.Instance.Spec.Ingestor.LocalTSDBRetention)
			}
			if configSecret := ingestorSpec.ObjectStorageConfig; configSecret != nil && configSecret.Name != "" {
				args = append(args, "--objstore.config-file="+filepath.Join(mountDirSecrets, configSecret.Name, configSecret.Key))
			} else {
				// TODO set max-block-duration to enable compact
			}
		}
	}
	if m == RouterOnly || m == RouterIngestor {
		args = append(args, "--receive.hashrings-file="+filepath.Join(thanosConfigDir, hashringsFileName))
		if routerSpec := r.Instance.Spec.Router; routerSpec != nil {
			if routerSpec.HashringsRefreshInterval != "" {
				args = append(args, "--receive.hashrings-file-refresh-interval="+routerSpec.HashringsRefreshInterval)
			}

			var replicationFactor uint64 = 1
			if routerSpec.ReplicationFactor != nil && *routerSpec.ReplicationFactor > 1 {
				replicationFactor = *routerSpec.ReplicationFactor
			}
			args = append(args, fmt.Sprintf("--receive.replication-factor=%d", replicationFactor))
		}
	}

	return args
}
