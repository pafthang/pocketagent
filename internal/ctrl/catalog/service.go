package catalog

// Service describes a managed child process for ctrl.
type Service struct {
	Name       string
	Package    string
	WaitPort   int
	HealthPort int
	DependsOn  []string
	Env        map[string]string
}