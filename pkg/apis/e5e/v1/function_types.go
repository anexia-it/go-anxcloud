package v1

type StorageBackend string

const (
	StorageBackendS3      = "s3"
	StorageBackendGit     = "git"
	StorageBackendZipFile = "zip"
)

type WorkerType string

const (
	WorkerTypeStandard = "standard"
)

// anxcloud:object

// Function represents the collection of all the metadata as well as the code itself
// that is needed to execute your application on the e5e platform.
type Function struct {
	omitResponseDecodeOnDestroy
	Identifier            string                 `json:"identifier,omitempty" anxcloud:"identifier"`
	State                 string                 `json:"state,omitempty"`
	Name                  string                 `json:"name,omitempty"`
	ApplicationIdentifier string                 `json:"application_identifier,omitempty"`
	Runtime               string                 `json:"runtime,omitempty"`
	Entrypoint            string                 `json:"entrypoint,omitempty"`
	StorageBackend        StorageBackend         `json:"storage_backend,omitempty"`
	StorageBackendMeta    *StorageBackendMeta    `json:"storage_backend_meta,omitempty"`
	EnvironmentVariables  *[]EnvironmentVariable `json:"environment_variables,omitempty"`
	Hostnames             *[]Hostname            `json:"hostnames,omitempty"`
	KeepAlive             int                    `json:"keep_alive,omitempty"`
	QuotaStorage          int                    `json:"quota_storage,omitempty"`
	QuotaMemory           int                    `json:"quota_memory,omitempty"`
	QuotaCPU              int                    `json:"quota_cpu,omitempty"`
	QuotaTimeout          int                    `json:"quota_timeout,omitempty"`
	QuotaConcurrency      int                    `json:"quota_concurrency,omitempty"`
	WorkerType            string                 `json:"worker_type,omitempty"`
}

// StorageBackendMeta is used to configure a storage backend
type StorageBackendMeta struct {
	*StorageBackendMetaGit
	*StorageBackendMetaS3
	*StorageBackendMetaZipFile `json:"zip_file,omitempty"`
}

// StorageBackendMetaGit is used to configure a git storage backend
type StorageBackendMetaGit struct {
	URL        string `json:"git_url,omitempty"`
	Branch     string `json:"git_branch,omitempty"`
	PrivateKey string `json:"git_private_key,omitempty"`
	Username   string `json:"git_username,omitempty"`
	Password   string `json:"git_password,omitempty"`
}

// StorageBackendMetaS3 is used to configure a s3 storage backend
type StorageBackendMetaS3 struct {
	Endpoint   string `json:"s3_endpoint,omitempty"`
	BucketName string `json:"s3_bucket_name,omitempty"`
	ObjectPath string `json:"s3_object_path,omitempty"`
	AccessKey  string `json:"s3_access_key,omitempty"`
	SecretKey  string `json:"s3_secret_key,omitempty"`
}

// StorageBackendMetaZipFile is used to configure a zip file storage backend
type StorageBackendMetaZipFile struct {
	// Data string containing the mime-type, encoding and encoded data
	// data:<mime type>;base64,<data>
	Content string `json:"content"`
	Name    string `json:"name"`
}

type EnvironmentVariable struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`
}

type Hostname struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}
