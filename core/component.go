package core

import (
	"time"
)

// Component represents a independent workflow component, handling only 1 input set.
type Component struct {

	// Identity specifies the unique identity of this component.
	Identity Identity

	// Jobs represent specific set of jobs which should be executed on this component.
	Jobs []*Job

	// CommonParameters represents the common env (and TODO: input, output) for all the workflow containers.
	CommonParameters *Parameters

	// SharedData ...
	SharedData *SharedData

	// Machine represents the spec of machines on which each job is executed.
	Machine *Machine

	/* You Ain't Gonna Need It!! */
	// // Nest can specify nested components.
	// // If "Nest" is provided not nil, all the "Jobs" are ignored.
	// // If neither "Parallel" nor "Serial" is provided, it results in an error.
	// Nest *struct {
	// 	Parallel []*Component
	// 	Serial   []*Component
	// }

	/* You Ain't Gonna Need It!! */
	// // Parent
	// Parent *Component

	Runtime struct {
		Image  *Image
		Script *Script
	}

	// TODO: Not yet used
	// RootLog is an application logger ONLY FOR ROOT COMPONENT.
	// RootLog *log.Logger

	// JobLoggerFactory is an interface to specify logger for each job.
	JobLoggerFactory LoggerFactory

	// Concurrency for creating machines.
	// We estimate that creating machines is the most costly process in job lifecycle,
	// therefore, this "Concurrency" prop should throttle the concurrency of them.
	Concurrency int64
}

// RootComponentTemplate ...
func RootComponentTemplate(name string) *Component {
	return &Component{
		Identity: Identity{Name: name, Timestamp: time.Now().UnixNano()},
		Machine:  &Machine{},
		Runtime: struct {
			Image  *Image
			Script *Script
		}{Image: &Image{}, Script: &Script{}},
		CommonParameters: &Parameters{},
		SharedData:       &SharedData{},
		Concurrency:      8,
	}
}
