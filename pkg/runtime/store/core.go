package store

import (
	"github.com/Aptomi/aptomi/pkg/engine"
	"github.com/Aptomi/aptomi/pkg/engine/actual"
	"github.com/Aptomi/aptomi/pkg/engine/progress"
	"github.com/Aptomi/aptomi/pkg/engine/resolve"
	"github.com/Aptomi/aptomi/pkg/lang"
	"github.com/Aptomi/aptomi/pkg/runtime"
)

// Core represents main object store interface that covers database operations for all objects
type Core interface {
	Policy
	Revision
	ActualState
}

// Policy represents database operations for Policy object
type Policy interface {
	GetPolicy(runtime.Generation) (*lang.Policy, runtime.Generation, error)
	GetPolicyData(runtime.Generation) (*engine.PolicyData, error)
	UpdatePolicy(updated []lang.Base, deleted []runtime.Key) (changed bool, data *engine.PolicyData, err error)
}

// Revision represents database operations for Revision object
type Revision interface {
	GetRevision(gen runtime.Generation) (*engine.Revision, error)
	GetFirstRevisionForPolicy(policyGen runtime.Generation) (*engine.Revision, error)
	NewRevision(policyGen runtime.Generation) (*engine.Revision, error)
	SaveRevision(revision *engine.Revision) error
	GetRevisionProgressUpdater(revision *engine.Revision) progress.Indicator
}

// ActualState represents database operations for the actual state handling
type ActualState interface {
	GetActualState() (*resolve.PolicyResolution, error)
	GetActualStateUpdater() actual.StateUpdater
}
