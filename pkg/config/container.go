package config

// TypeContainer is the resource string for a Container resource
const TypeContainer ResourceType = "container"

// Container defines a structure for creating Docker containers
type Container struct {
	// embedded type holding name, etc
	ResourceInfo

	Depends []string `hcl:"depends_on,optional" json:"depends,omitempty"`

	Networks []NetworkAttachment `hcl:"network,block" json:"networks,omitempty"` // Attach to the correct network // only when Image is specified

	Image       *Image            `hcl:"image,block" json:"image"`                        // Image to use for the container
	Build       *Build            `hcl:"build,block" json:"build"`                        // Enables containers to be built on the fly
	Entrypoint  []string          `hcl:"entrypoint,optional" json:"entrypoint,omitempty"` // entrypoint to use when starting the container
	Command     []string          `hcl:"command,optional" json:"command,omitempty"`       // command to use when starting the container
	Environment []KV              `hcl:"env,block" json:"environment,omitempty"`          // environment variables to set when starting the container, // Depricated field
	EnvVar      map[string]string `hcl:"env_var,optional" json:"env_var,omitempty"`       // environment variables to set when starting the container
	Volumes     []Volume          `hcl:"volume,block" json:"volumes,omitempty"`           // volumes to attach to the container
	Ports       []Port            `hcl:"port,block" json:"ports,omitempty"`               // ports to expose
	PortRanges  []PortRange       `hcl:"port_range,block" json:"port_ranges,omitempty"`   // range of ports to expose

	Privileged bool `hcl:"privileged,optional" json:"privileged,omitempty"` // run the container in privileged mode?

	// resource constraints
	Resources *Resources `hcl:"resources,block" json:"resources,omitempty"` // resource constraints for the container

	// health checks for the container
	HealthCheck *HealthCheck `hcl:"health_check,block" json:"health_check,omitempty"`
}

// NewContainer returns a new Container resource with the correct default options
func NewContainer(name string) *Container {
	return &Container{ResourceInfo: ResourceInfo{Name: name, Type: TypeContainer, Status: PendingCreation}}
}

type NetworkAttachment struct {
	Name      string   `hcl:"name" json:"name"`
	IPAddress string   `hcl:"ip_address,optional" json:"ip_address,omitempty"`
	Aliases   []string `hcl:"aliases,optional" json:"aliases,omitempty"` // Network aliases for the resource
}

// Resources allows the setting of resource constraints for the Container
type Resources struct {
	CPU    int   `hcl:"cpu,optional" json:"cpu,omitempty"`         // cpu limit for the container where 1 CPU = 1000
	CPUPin []int `hcl:"cpu_pin,optional" json:"cpu_pin,omitempty"` // pin the container to one or more cpu cores
	Memory int   `hcl:"memory,optional" json:"memory,omitempty"`   // max memory the container can consume in MB
}

// Volume defines a folder, Docker volume, or temp folder to mount to the Container
type Volume struct {
	Source      string `hcl:"source" json:"source"`                          // source path on the local machine for the volume
	Destination string `hcl:"destination" json:"destination"`                // path to mount the volume inside the container
	Type        string `hcl:"type,optional" json:"type,omitempty"`           // type of the volume to mount [bind, volume, tmpfs]
	ReadOnly    bool   `hcl:"read_only,optional" json:"read_only,omitempty"` // specify that the volume is mounted read only
}

// KV is a key/value type
type KV struct {
	Key   string `hcl:"key" json:"key"`
	Value string `hcl:"value" json:"value"`
}

// Build allows you to define the conditions for building a container
// on run from a Dockerfile
type Build struct {
	File    string `hcl:"file,optional" json:"file,omitempty"` // Location of build file inside build context defaults to ./Dockerfile
	Context string `hcl:"context" json:"context"`              // Path to build context
}

// Validate the config
func (c *Container) Validate() error {
	return nil
}
