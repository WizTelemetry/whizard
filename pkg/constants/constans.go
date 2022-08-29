package constants

const (
	DefaultTenantHeader    = "WHIZARD-TENANT"
	DefaultTenantId        = "default-tenant"
	DefaultTenantLabelName = "tenant_id"

	ServiceLabelKey = "monitoring.whizard.io/service"
	StorageLabelKey = "monitoring.whizard.io/storage"
	TenantLabelKey  = "monitoring.whizard.io/tenant"

	FinalizerIngester  = "finalizers.monitoring.whizard.io/ingester"
	FinalizerCompactor = "finalizers.monitoring.whizard.io/compactor"
	FinalizerDeletePVC = "finalizers.monitoring.whizard.io/deletePVC"

	LocalStorage = "default_storage.local"

	GRPCPort        = 10901
	HTTPPort        = 10902
	RemoteWritePort = 19291

	GRPCPortName        = "grpc"
	HTTPPortName        = "http"
	RemoteWritePortName = "remote-write"

	ReceiveReplicaLabelName = "receive_replica"
	RulerReplicaLabelName   = "ruler_replica"

	AppNameGateway       = "whizard-monitoring-gateway"
	AppNameQuery         = "query"
	AppNameQueryFrontend = "query-frontend"
	AppNameRouter        = "router"
	AppNameIngester      = "ingester"
	AppNameRuler         = "ruler"
	AppNameStore         = "store"
	AppNameCompactor     = "compactor"

	ServiceNameSuffix = "operated"

	LabelNameAppName      = "app.kubernetes.io/name"
	LabelNameAppManagedBy = "app.kubernetes.io/managed-by"
	LabelNameAppPartOf    = "app.kubernetes.io/part-of"

	LabelNameIngesterState        = "monitoring.whizard.io/ingester-state"
	LabelNameIngesterDeletingTime = "monitoring.whizard.io/ingester-deleting-time"

	LabelNameRulerShardSn = "monitoring.whizard.io/ruler-shard-sn"

	ConfigPath     = "/etc/whizard/"
	StorageDir     = "/whizard"
	TSDBVolumeName = "tsdb"

	LabelNameStorageHash = "monitoring.whizard.io/storage-hash"
	LabelNameTenantHash  = "monitoring.whizard.io/tenant-hash"

	TenantHash  = "TENANT_HASH"
	StorageHash = "STORAGE_HASH"
)

const (
	StorageProviderFILESYSTEM string = "FILESYSTEM"
	StorageProviderGCS        string = "GCS"
	StorageProviderS3         string = "S3"
	StorageProviderAZURE      string = "AZURE"
	StorageProviderSWIFT      string = "SWIFT"
	StorageProviderCOS        string = "COS"
	StorageProviderALIYUNOSS  string = "ALIYUNOSS"
	StorageProviderBOS        string = "BOS"
)