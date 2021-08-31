package thanosstorage

import (
	"fmt"
	"path/filepath"
)

func (s *ThanosStorage) gatewayArgs() []string {
	var (
		gatewaySpec      = s.Instance.Spec.Gateway
		objStorageConfig = s.Instance.Spec.ObjectStorageConfig

		args = []string{
			"store",
			fmt.Sprintf("--grpc-address=[$(POD_IP)]:%d", grpcPort),
			fmt.Sprintf("--http-address=[$(POD_IP)]:%d", httpPort),
			fmt.Sprintf(`--data-dir="%s"`, storageDir),
		}
	)

	if objStorageConfig != nil && objStorageConfig.Name != "" {
		args = append(args, "--objstore.config-file="+filepath.Join(mountDirSecrets, objStorageConfig.Name, objStorageConfig.Key))
	}

	if gatewaySpec.MinTime != "" {
		args = append(args, "--min-time="+gatewaySpec.MinTime)
	}
	if gatewaySpec.MaxTime != "" {
		args = append(args, "--max-time="+gatewaySpec.MaxTime)
	}

	if s.Instance.Spec.LogLevel != "" {
		args = append(args, "--log.level="+s.Instance.Spec.LogLevel)
	}
	if s.Instance.Spec.LogFormat != "" {
		args = append(args, "--log.format="+s.Instance.Spec.LogFormat)
	}

	return args
}

func (s *ThanosStorage) compactArgs() []string {
	var (
		compactSpec      = s.Instance.Spec.Compact
		objStorageConfig = s.Instance.Spec.ObjectStorageConfig

		args = []string{
			"compact",
			"--wait",
			fmt.Sprintf("--http-address=[$(POD_IP)]:%d", httpPort),
			fmt.Sprintf(`--data-dir="%s"`, storageDir),
		}
	)

	if compactSpec.DownsamplingDisable != nil {
		args = append(args, fmt.Sprintf("--downsampling.disable=%v", compactSpec.DownsamplingDisable))
	}

	if objStorageConfig != nil && objStorageConfig.Name != "" {
		args = append(args, "--objstore.config-file="+filepath.Join(mountDirSecrets, objStorageConfig.Name, objStorageConfig.Key))
	}

	if s.Instance.Spec.LogLevel != "" {
		args = append(args, "--log.level="+s.Instance.Spec.LogLevel)
	}
	if s.Instance.Spec.LogFormat != "" {
		args = append(args, "--log.format="+s.Instance.Spec.LogFormat)
	}

	if retention := compactSpec.Retention; retention != nil {
		if retention.RetentionRaw != "" {
			args = append(args, "--retention.resolution-raw"+retention.RetentionRaw)
		}
		if retention.Retention5m != "" {
			args = append(args, "--retention.resolution-5m"+retention.Retention5m)
		}
		if retention.Retention5m != "" {
			args = append(args, "--retention.resolution-1h"+retention.Retention5m)
		}
	}

	return args
}
