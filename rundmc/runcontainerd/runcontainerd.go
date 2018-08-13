package runcontainerd

import (
	"fmt"
	"io"
	"syscall"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gardener"
	"code.cloudfoundry.org/guardian/rundmc/goci"
	"code.cloudfoundry.org/guardian/rundmc/runrunc"
	"code.cloudfoundry.org/guardian/rundmc/users"
	"code.cloudfoundry.org/lager"
	uuid "github.com/nu7hatch/gouuid"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

//go:generate counterfeiter . ContainerManager
type ContainerManager interface {
	Create(log lager.Logger, containerID string, spec *specs.Spec) error
	Delete(log lager.Logger, containerID string) error

	Exec(log lager.Logger, containerID, processID string, spec *specs.Process, processIO func() (io.Reader, io.Writer, io.Writer)) error

	State(log lager.Logger, containerID string) (int, string, error)
	GetContainerPID(log lager.Logger, containerID string) (uint32, error)
	GetNamespace() (string, error)
}

//go:generate counterfeiter . ProcessManager
type ProcessManager interface {
	Wait(log lager.Logger, containerID, processID string) (int, error)
	Signal(log lager.Logger, containerID, processID string, signal syscall.Signal) error
	DeleteProcess(log lager.Logger, containerID, processID string) error
}

//go:generate counterfeiter . BundleLoader
type BundleLoader interface {
	Load(string) (goci.Bndl, error)
}

//go:generate counterfeiter . ProcessBuilder
type ProcessBuilder interface {
	BuildProcess(bndl goci.Bndl, spec garden.ProcessSpec, uid, gid int) *specs.Process
}

//go:generate counterfeiter . Execer
type Execer interface {
	Exec(log lager.Logger, bundlePath string, id string, spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error)
	Attach(log lager.Logger, bundlePath string, id string, processId string, io garden.ProcessIO) (garden.Process, error)
}

//go:generate counterfeiter . Statser
type Statser interface {
	Stats(log lager.Logger, id string) (gardener.ActualContainerMetrics, error)
}

type RunContainerd struct {
	containerManager          ContainerManager
	processManager            ProcessManager
	bundleLoader              BundleLoader
	processBuilder            ProcessBuilder
	execer                    Execer
	statser                   Statser
	useContainerdForProcesses bool
	userLookupper             users.UserLookupper
	cgroupManager             CgroupManager
}

func New(containerManager ContainerManager, processManager ProcessManager, bundleLoader BundleLoader, processBuilder ProcessBuilder, userLookupper users.UserLookupper, execer Execer, statser Statser, useContainerdForProcesses bool, cgroupManager CgroupManager) *RunContainerd {
	return &RunContainerd{
		containerManager:          containerManager,
		processManager:            processManager,
		bundleLoader:              bundleLoader,
		processBuilder:            processBuilder,
		execer:                    execer,
		statser:                   statser,
		useContainerdForProcesses: useContainerdForProcesses,
		userLookupper:             userLookupper,
		cgroupManager:             cgroupManager,
	}
}

// TODO io not used
func (r *RunContainerd) Create(log lager.Logger, bundlePath, id string, io garden.ProcessIO) error {
	bundle, err := r.bundleLoader.Load(bundlePath)
	if err != nil {
		return err
	}

	err = r.containerManager.Create(log, id, &bundle.Spec)
	if err != nil {
		return err
	}

	if r.useContainerdForProcesses {
		return r.cgroupManager.SetUseMemoryHierarchy(id)
	}
	return nil
}

func (r *RunContainerd) Exec(log lager.Logger, bundlePath, containerID string, gardenProcessSpec garden.ProcessSpec, gardenIO garden.ProcessIO) (garden.Process, error) {
	if !r.useContainerdForProcesses {
		return r.execer.Exec(log, bundlePath, containerID, gardenProcessSpec, gardenIO)
	}

	bundle, err := r.bundleLoader.Load(bundlePath)
	if err != nil {
		return nil, err
	}

	containerPid, err := r.containerManager.GetContainerPID(log, containerID)
	if err != nil {
		return nil, err
	}

	resolvedUser, err := r.userLookupper.Lookup(fmt.Sprintf("/proc/%d/root", containerPid), gardenProcessSpec.User)
	if err != nil {
		return nil, err
	}

	if gardenProcessSpec.Dir == "" {
		gardenProcessSpec.Dir = resolvedUser.Home
	}

	// TODO: use the uidgenerator
	if gardenProcessSpec.ID == "" {
		randomID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		gardenProcessSpec.ID = fmt.Sprintf("%s", randomID)
	}

	processIO := func() (io.Reader, io.Writer, io.Writer) {
		return gardenIO.Stdin, gardenIO.Stdout, gardenIO.Stderr
	}

	ociProcessSpec := r.processBuilder.BuildProcess(bundle, gardenProcessSpec, resolvedUser.Uid, resolvedUser.Gid)
	if err = r.containerManager.Exec(log, containerID, gardenProcessSpec.ID, ociProcessSpec, processIO); err != nil {
		return nil, err
	}

	return NewProcess(log, containerID, gardenProcessSpec.ID, r.processManager), nil
}

func (r *RunContainerd) Attach(log lager.Logger, bundlePath, id, processId string, io garden.ProcessIO) (garden.Process, error) {
	return r.execer.Attach(log, bundlePath, id, processId, io)
}

func (r *RunContainerd) Kill(log lager.Logger, bundlePath string) error {
	return fmt.Errorf("Kill is not implemented yet")
}

func (r *RunContainerd) Delete(log lager.Logger, force bool, id string) error {
	return r.containerManager.Delete(log, id)
}

func (r *RunContainerd) State(log lager.Logger, id string) (runrunc.State, error) {
	pid, status, err := r.containerManager.State(log, id)
	if err != nil {
		return runrunc.State{}, err
	}

	return runrunc.State{Pid: pid, Status: runrunc.Status(status)}, nil
}

func (r *RunContainerd) Stats(log lager.Logger, id string) (gardener.ActualContainerMetrics, error) {
	return r.statser.Stats(log, id)
}

func (r *RunContainerd) WatchEvents(log lager.Logger, id string, eventsNotifier runrunc.EventsNotifier) error {
	return fmt.Errorf("WatchEvents is not implemented yet")
}
