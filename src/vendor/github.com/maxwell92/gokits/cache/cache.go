package cache

type Cache interface {
	// ToDo: support multi-tenancy： namespace like k8s
	Create(db string) error
	Index(columns []string) error
	Update(key, value string) error
	Delete(key string) error
	Search(keyword string) (*map[string]string, error)
	Watch(name string) error
	Find(key string) ([]string, error)
}
