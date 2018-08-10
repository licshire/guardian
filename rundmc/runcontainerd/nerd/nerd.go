package nerd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"syscall"

	"code.cloudfoundry.org/lager"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/linux/runctypes"
	"github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type Nerd struct {
	client  *containerd.Client
	context context.Context
}

func New(client *containerd.Client, context context.Context) *Nerd {
	return &Nerd{
		client:  client,
		context: context,
	}
}

// TODO: didn't we PR this?
func WithNoNewKeyring(ctx context.Context, c *containerd.Client, ti *containerd.TaskInfo) error {
	ti.Options = &runctypes.CreateOptions{NoNewKeyring: true}
	return nil
}

func (n *Nerd) Create(log lager.Logger, containerID string, spec *specs.Spec) error {
	log.Debug("creating-container", lager.Data{"containerID": containerID})
	container, err := n.client.NewContainer(n.context, containerID, containerd.WithSpec(spec))
	if err != nil {
		return err
	}

	log.Debug("creating-task", lager.Data{"containerID": containerID})
	task, err := container.NewTask(n.context, cio.NullIO, WithNoNewKeyring)
	if err != nil {
		return err
	}

	log.Debug("starting-task", lager.Data{"containerID": containerID})
	return task.Start(n.context)
}

func (n *Nerd) Delete(log lager.Logger, containerID string) error {
	container, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return err
	}

	log.Debug("deleting-task", lager.Data{"containerID": containerID})
	_, err = task.Delete(n.context, withProcessKill)
	if err != nil {
		return err
	}

	log.Debug("deleting-container", lager.Data{"containerID": containerID})
	return container.Delete(n.context)
}

// WithProcessKill will forcefully kill and delete a process
func withProcessKill(ctx context.Context, p containerd.Process) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// ignore errors to wait and kill as we are forcefully killing
	// the process and don't care about the exit status
	s, err := p.Wait(ctx)
	if err != nil {
		return err
	}
	if err := p.Kill(ctx, syscall.SIGKILL); err != nil {
		if errdefs.IsFailedPrecondition(err) || errdefs.IsNotFound(err) {
			return nil
		}
		return err
	}
	// wait for the process to fully stop before letting the rest of the deletion complete
	<-s
	return nil
}

func (n *Nerd) State(log lager.Logger, containerID string) (int, string, error) {
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return 0, "", err
	}

	log.Debug("getting-task-status", lager.Data{"containerID": containerID})
	status, err := task.Status(n.context)
	if err != nil {
		return 0, "", err
	}

	log.Debug("task-result", lager.Data{"containerID": containerID, "pid": strconv.Itoa(int(task.Pid())), "status": string(status.Status)})
	return int(task.Pid()), string(status.Status), nil
}

func (n *Nerd) Exec(log lager.Logger, containerID, processID string, spec *specs.Process, processIO func() (io.Reader, io.Writer, io.Writer)) error {
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return err
	}

	log.Debug("execing-task", lager.Data{"containerID": containerID, "processID": processID})
	process, err := task.Exec(n.context, processID, spec, cio.NewCreator(withProcessIO(processIO)))
	if err != nil {
		return err
	}

	log.Debug("starting-task", lager.Data{"containerID": containerID, "processID": processID})
	if err := process.Start(n.context); err != nil {
		return err
	}

	log.Debug("closing-stdin", lager.Data{"containerID": containerID, "processID": processID})
	return process.CloseIO(n.context, containerd.WithStdinCloser)
}

func (n *Nerd) DeleteProcess(log lager.Logger, containerID, processID string) error {
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return err
	}

	process, err := task.LoadProcess(n.context, processID, nil)
	if err != nil {
		return err
	}

	_, err = process.Delete(n.context)
	if err != nil {
		return err
	}

	return nil
}

func withProcessIO(processIO func() (io.Reader, io.Writer, io.Writer)) cio.Opt {
	return func(opt *cio.Streams) {
		stdIn, stdOut, stdErr := processIO()
		if stdIn != nil {
			opt.Stdin = stdIn
		} else {
			opt.Stdin = bytes.NewBuffer(nil)
		}

		if stdOut != nil {
			opt.Stdout = stdOut
		} else {
			opt.Stdout = ioutil.Discard
		}

		if stdErr != nil {
			opt.Stderr = stdErr
		} else {
			opt.Stderr = ioutil.Discard
		}
	}
}

func (n *Nerd) GetContainerPID(log lager.Logger, containerID string) (uint32, error) {
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return 0, err
	}

	return task.Pid(), nil
}

func (n *Nerd) loadContainerAndTask(log lager.Logger, containerID string) (containerd.Container, containerd.Task, error) {
	log.Debug("loading-container", lager.Data{"containerID": containerID})
	container, err := n.client.LoadContainer(n.context, containerID)
	if err != nil {
		return nil, nil, err
	}

	log.Debug("loading-task", lager.Data{"containerID": containerID})
	task, err := container.Task(n.context, nil)
	if err != nil {
		return nil, nil, err
	}

	return container, task, nil
}

func (n *Nerd) Wait(log lager.Logger, containerID, processID string) (int, error) {
	log.Debug("waiting-on-process", lager.Data{"containerID": containerID, "processID": processID})
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return 0, err
	}

	process, err := task.LoadProcess(n.context, processID, nil)
	if err != nil {
		return 0, err
	}

	exitCh, err := process.Wait(n.context)
	if err != nil {
		return 0, err
	}

	// Containerd might fail to retrieve the ExitCode for non-process related reasons
	exitStatus := <-exitCh
	if exitStatus.Error() != nil {
		return 0, exitStatus.Error()
	}

	return int(exitStatus.ExitCode()), nil
}

func (n *Nerd) Signal(log lager.Logger, containerID, processID string, signal syscall.Signal) error {
	log.Debug("signalling-process", lager.Data{"containerID": containerID, "processID": processID, "signal": signal})
	_, task, err := n.loadContainerAndTask(log, containerID)
	if err != nil {
		return err
	}

	process, err := task.LoadProcess(n.context, processID, nil)
	if err != nil {
		return err
	}

	return process.Kill(n.context, signal)
}

func (n *Nerd) GetNamespace() (string, error) {
	namespace, ok := namespaces.Namespace(n.context)
	if !ok {
		return "", errors.New("could not get namespace for container manager")
	}

	return namespace, nil
}
