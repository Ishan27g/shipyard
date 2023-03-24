package resources

import "github.com/shipyard-run/hclconfig/types"

// TypeContainer is the resource string for a Container resource
const TypeContainer string = "container"

// Container defines a structure for creating Docker containers
type Container struct {
	// embedded type holding name, etc
	types.ResourceMetadata `hcl:",remain"`

	Depends []string `hcl:"depends_on,optional" json:"depends,omitempty"`

	Networks []NetworkAttachment `hcl:"network,block" json:"networks,omitempty"` // Attach to the correct network // only when Image is specified

	Image      *Image            `hcl:"image,block" json:"image"`                        // Image to use for the container
	Build      *Build            `hcl:"build,block" json:"build"`                        // Enables containers to be built on the fly
	Entrypoint []string          `hcl:"entrypoint,optional" json:"entrypoint,omitempty"` // entrypoint to use when starting the container
	Command    []string          `hcl:"command,optional" json:"command,omitempty"`       // command to use when starting the container
	Env        map[string]string `hcl:"env,optional" json:"env,omitempty"`               // environment variables to set when starting the container
	Volumes    []Volume          `hcl:"volume,block" json:"volumes,omitempty"`           // volumes to attach to the container
	Ports      []Port            `hcl:"port,block" json:"ports,omitempty"`               // ports to expose
	PortRanges []PortRange       `hcl:"port_range,block" json:"port_ranges,omitempty"`   // range of ports to expose
	DNS        []string          `hcl:"dns,optional" json:"dns,omitempty"`               // Add custom DNS servers to the container

	Privileged bool `hcl:"privileged,optional" json:"privileged,omitempty"` // run the container in privileged mode?

	// resource constraints
	Resources *Resources `hcl:"resources,block" json:"resources,omitempty"` // resource constraints for the container

	// health checks for the container
	HealthCheck *HealthCheck `hcl:"health_check,block" json:"health_check,omitempty"`

	MaxRestartCount int `hcl:"max_restart_count,optional" json:"max_restart_count,omitempty"`

	// User block for mapping the user id and group id inside the container
	RunAs *User `hcl:"run_as,block" json:"run_as,omitempty"`
}

type User struct {
	// Username or UserID of the user to run the container as
	User string `hcl:"user" json:"user,omitempty"`
	// Group is the GroupID of the user to run the container as
	Group string `hcl:"group" json:"group,omitempty"`
}

type NetworkAttachment struct {
	ID        string   `hcl:"id" json:"id"`
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
	Source                      string `hcl:"source" json:"source"`                                                                    // source path on the local machine for the volume
	Destination                 string `hcl:"destination" json:"destination"`                                                          // path to mount the volume inside the container
	Type                        string `hcl:"type,optional" json:"type,omitempty"`                                                     // type of the volume to mount [bind, volume, tmpfs]
	ReadOnly                    bool   `hcl:"read_only,optional" json:"read_only,omitempty"`                                           // specify that the volume is mounted read only
	BindPropagation             string `hcl:"bind_propagation,optional" json:"bind_propagation,omitempty"`                             // propagation mode for bind mounts [shared, private, slave, rslave, rprivate]
	BindPropagationNonRecursive bool   `hcl:"bind_propagation_non_recursive,optional" json:"bind_propagation_non_recursive,omitempty"` // recursive bind mount, default true
}

// Build allows you to define the conditions for building a container
// on run from a Dockerfile
type Build struct {
	File    string `hcl:"file,optional" json:"file,omitempty"` // Location of build file inside build context defaults to ./Dockerfile
	Context string `hcl:"context" json:"context"`              // Path to build context
	Tag     string `hcl:"tag,optional" json:"tag,omitempty"`   // Image tag, defaults to latest
}

func (c *Container) Process() error {
	// process volumes
	for i, v := range c.Volumes {
		// make sure mount paths are absolute when type is bind
		if v.Type == "" || v.Type == "bind" {
			c.Volumes[i].Source = ensureAbsolute(v.Source, c.File)
		}
	}

	// make sure build paths are absolute
	if c.Build != nil {
		c.Build.File = ensureAbsolute(c.Build.File, c.File)
		c.Build.Context = ensureAbsolute(c.Build.Context, c.File)
	}

	return nil
}