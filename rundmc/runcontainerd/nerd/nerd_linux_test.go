package nerd_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/oci"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	specs "github.com/opencontainers/runtime-spec/specs-go"

	"code.cloudfoundry.org/guardian/rundmc/runcontainerd/nerd"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
)

var _ = Describe("Nerd", func() {
	var (
		testLogger  lager.Logger
		cnerd       *nerd.Nerd
		containerID string
		processID   string
		processIO   func() (io.Reader, io.Writer, io.Writer)
	)

	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
		containerID = fmt.Sprintf("test-container-%s", randomString(10))
		processID = fmt.Sprintf("test-process-%s", randomString(10))
		processIO = func() (io.Reader, io.Writer, io.Writer) {
			return nil, nil, nil
		}

		testLogger = lagertest.NewTestLogger("nerd-test")
		cnerd = nerd.New(containerdClient, containerdContext)
	})

	Describe("Create", func() {
		AfterEach(func() {
			cnerd.Delete(testLogger, containerID)
		})

		It("creates the containerd container by id", func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())

			containers := listContainers(testConfig.CtrBin, testConfig.Socket)
			Expect(containers).To(ContainSubstring(containerID))
		})

		It("starts an init process in the container", func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())

			containers := listProcesses(testConfig.CtrBin, testConfig.Socket, containerID)
			Expect(containers).To(ContainSubstring(containerID))
		})
	})

	Describe("Exec", func() {
		BeforeEach(func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())
		})

		AfterEach(func() {
			cnerd.Delete(testLogger, containerID)
		})

		It("execs a process in the container", func() {
			processSpec := &specs.Process{
				Args: []string{"/bin/sleep", "30"},
				Cwd:  "/",
			}

			err := cnerd.Exec(testLogger, containerID, processID, processSpec, processIO)
			Expect(err).NotTo(HaveOccurred())

			containers := listProcesses(testConfig.CtrBin, testConfig.Socket, containerID)
			Expect(containers).To(ContainSubstring(containerID)) // init process
			Expect(containers).To(ContainSubstring(processID))   // execed process
		})

		Describe("process IO", func() {
			It("reads from stdin", func() {
				processSpec := &specs.Process{
					Args: []string{"/bin/cat", "/proc/self/fd/0"},
					Cwd:  "/",
				}

				stdout := gbytes.NewBuffer()
				processIO = func() (io.Reader, io.Writer, io.Writer) {
					stdin := gbytes.BufferWithBytes([]byte("hello nerd"))
					stdout := io.MultiWriter(stdout, GinkgoWriter)

					return stdin, stdout, nil
				}

				err := cnerd.Exec(testLogger, containerID, processID, processSpec, processIO)
				Expect(err).NotTo(HaveOccurred())
				Eventually(stdout).Should(gbytes.Say("hello nerd"))
			})

			It("writes to stdout", func() {
				processSpec := &specs.Process{
					Args: []string{"/bin/echo", "hello nerd"},
					Cwd:  "/",
				}

				stdout := gbytes.NewBuffer()
				processIO = func() (io.Reader, io.Writer, io.Writer) {
					stdout := io.MultiWriter(stdout, GinkgoWriter)
					return nil, stdout, nil
				}

				err := cnerd.Exec(testLogger, containerID, processID, processSpec, processIO)
				Expect(err).NotTo(HaveOccurred())
				Eventually(stdout).Should(gbytes.Say("hello nerd"))
			})

			It("writes to stderr", func() {
				processSpec := &specs.Process{
					Args: []string{"/bin/cat", "notafile.txt"},
					Cwd:  "/",
				}

				stderr := gbytes.NewBuffer()
				processIO = func() (io.Reader, io.Writer, io.Writer) {
					stderr := io.MultiWriter(stderr, GinkgoWriter)
					return nil, nil, stderr
				}

				err := cnerd.Exec(testLogger, containerID, processID, processSpec, processIO)
				Expect(err).NotTo(HaveOccurred())
				Eventually(stderr).Should(gbytes.Say("No such file"))
			})
		})
	})

	Describe("Wait", func() {
		var (
			exitCode int
			waitErr  error
		)

		BeforeEach(func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())

			processSpec := &specs.Process{
				Args: []string{"/bin/sh", "-c", "exit 17"},
				Cwd:  "/",
			}

			err := cnerd.Exec(testLogger, containerID, processID, processSpec, processIO)
			Expect(err).NotTo(HaveOccurred())

			exitCode, waitErr = cnerd.Wait(testLogger, containerID, processID)
		})

		AfterEach(func() {
			cnerd.Delete(testLogger, containerID)
		})

		It("succeeds", func() {
			Expect(waitErr).NotTo(HaveOccurred())
		})

		It("returns the exit code", func() {
			Expect(exitCode).To(Equal(17))
		})
	})

	Describe("Signal", func() {
		BeforeEach(func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())
			stdoutBuffer := gbytes.NewBuffer()
			processIO = func() (io.Reader, io.Writer, io.Writer) {
				return nil, stdoutBuffer, nil
			}

			processSpec := &specs.Process{
				Args: []string{"/bin/sh", "-c", `
					trap 'exit 42' TERM

					while true; do
					  echo 'sleeping'
					  sleep 1
					done
				`},
				Cwd: "/",
			}

			err := cnerd.Exec(testLogger, containerID, processID, processSpec, processIO)
			Expect(err).NotTo(HaveOccurred())

			Eventually(stdoutBuffer).Should(gbytes.Say("sleeping"))
		})

		AfterEach(func() {
			cnerd.Delete(testLogger, containerID)
		})

		It("should forward signals to the process", func() {
			Expect(cnerd.Signal(testLogger, containerID, processID, syscall.SIGTERM)).To(Succeed())

			status := make(chan int)
			go func() {
				exit, err := cnerd.Wait(testLogger, containerID, processID)
				Expect(err).NotTo(HaveOccurred())
				status <- exit
			}()

			Eventually(status, 5*time.Second).Should(Receive(BeEquivalentTo(42)))
		})

	})

	Describe("Delete", func() {
		BeforeEach(func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())
		})

		It("deletes the containerd container by id", func() {
			Expect(cnerd.Delete(testLogger, containerID)).To(Succeed())

			containers := listContainers(testConfig.CtrBin, testConfig.Socket)
			Expect(containers).NotTo(ContainSubstring(containerID))
		})
	})

	Describe("State", func() {
		BeforeEach(func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())
		})

		AfterEach(func() {
			cnerd.Delete(testLogger, containerID)
		})

		It("gets the pid and status of a running task", func() {
			pid, status, err := cnerd.State(testLogger, containerID)

			Expect(err).NotTo(HaveOccurred())
			Expect(pid).NotTo(BeZero())
			Expect(status).To(BeEquivalentTo(containerd.Running))
		})
	})

	Describe("GetContainerPID", func() {
		BeforeEach(func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())
		})

		AfterEach(func() {
			cnerd.Delete(testLogger, containerID)
		})

		It("gets the container init process pid", func() {
			procls := listProcesses(testConfig.CtrBin, testConfig.Socket, containerID)
			containerPid, err := cnerd.GetContainerPID(testLogger, containerID)
			Expect(err).ToNot(HaveOccurred())
			Expect(procls).To(ContainSubstring(strconv.Itoa(int(containerPid))))
		})
	})

	Describe("DeleteProcess", func() {
		BeforeEach(func() {
			spec := generateSpec(containerdContext, containerdClient, containerID)
			Expect(cnerd.Create(testLogger, containerID, spec)).To(Succeed())

			processSpec := &specs.Process{
				Args: []string{"/bin/sh", "-c", "exit 17"},
				Cwd:  "/",
			}

			err := cnerd.Exec(testLogger, containerID, processID, processSpec, processIO)
			Expect(err).NotTo(HaveOccurred())

			exitCode, err := cnerd.Wait(testLogger, containerID, processID)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(17))
		})

		AfterEach(func() {
			cnerd.Delete(testLogger, containerID)
		})

		It("removes process metadata", func() {
			err := cnerd.DeleteProcess(testLogger, containerID, processID)
			Expect(err).NotTo(HaveOccurred())

			err = cnerd.DeleteProcess(testLogger, containerID, processID)
			Expect(err).To(MatchError(ContainSubstring("not found")))
		})

		Context("when the container does not exist", func() {
			It("fails", func() {
				err := cnerd.DeleteProcess(testLogger, "foo", processID)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	FDescribe("GetNamespace", func() {
		It("returns the namespace", func() {
			namespace, err := cnerd.GetNamespace()
			Expect(err).NotTo(HaveOccurred())
			Expect(namespace).To(Equal(fmt.Sprintf("nerdspace%d", GinkgoParallelNode())))
		})

		Context("when we create a nerd with no context", func() {
			It("returns an error", func() {
				Expect(nerd.New(containerdClient, nil).GetNamespace()).To(MatchError("could not get namespace for container manager"))
			})
		})
	})
})

func createRootfs(modifyRootfs func(string), perm os.FileMode) string {
	var err error
	tmpDir, err := ioutil.TempDir("", "test-rootfs")
	Expect(err).NotTo(HaveOccurred())
	unpackedRootfs := filepath.Join(tmpDir, "unpacked")
	Expect(os.Mkdir(unpackedRootfs, perm)).To(Succeed())
	runCommand(exec.Command("tar", "xf", os.Getenv("GARDEN_TEST_ROOTFS"), "-C", unpackedRootfs))
	Expect(os.Chmod(tmpDir, perm)).To(Succeed())
	modifyRootfs(unpackedRootfs)
	return unpackedRootfs
}

func generateSpec(context context.Context, client *containerd.Client, containerID string) *specs.Spec {
	spec, err := oci.GenerateSpec(context, client, &containers.Container{ID: containerID})
	Expect(err).NotTo(HaveOccurred())
	spec.Process.Args = []string{"sleep", "60"}
	spec.Root = &specs.Root{
		Path: createRootfs(func(_ string) {}, 0755),
	}

	return spec
}

func listContainers(ctr, socket string) string {
	return runCtr(ctr, socket, []string{"containers", "list"})
}

func listProcesses(ctr, socket, containerID string) string {
	return runCtr(ctr, socket, []string{"tasks", "ps", containerID})
}

func runCtr(ctr, socket string, args []string) string {
	defaultArgs := []string{"--address", socket, "--namespace", fmt.Sprintf("nerdspace%d", GinkgoParallelNode())}
	cmd := exec.Command(ctr, append(defaultArgs, args...)...)

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))

	return string(session.Out.Contents())
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
