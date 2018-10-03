// +build !windows

package gqt_test

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"code.cloudfoundry.org/garden"
	"code.cloudfoundry.org/guardian/gqt/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Tmpfs mount", func() {
	var (
		client    *runner.RunningGarden
		container garden.Container
		path      string
	)

	BeforeEach(func() {
		path = "/tmpfs/path"
	})

	JustBeforeEach(func() {
		client = runner.Start(config)

		var err error

		container, err = client.Create(
			garden.ContainerSpec{
				Tmpfs: []garden.Tmpfs{{
					Path: path,
				}},
				Network: fmt.Sprintf("10.0.%d.0/24", GinkgoParallelNode()),
			})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(client.DestroyAndStop()).To(Succeed())
	})

	It("creates a tmpfs in the container at the given path", func() {
		Expect(containerStatFile(container, path)).To(gbytes.Say(path))
		Expect(getMountType(container, path)).To(Equal("tmpfs"))
	})

	It("the tmpfs is writable", func() {
		Expect(userWriteFile(container, path, "root").Wait()).To(Equal(0))
	})
})

func containerStatFile(container garden.Container, filePath string) *gbytes.Buffer {
	output := gbytes.NewBuffer()
	process, err := container.Run(garden.ProcessSpec{
		Path: "stat",
		Args: []string{filePath},
	}, garden.ProcessIO{
		Stdout: io.MultiWriter(GinkgoWriter, output),
		Stderr: io.MultiWriter(GinkgoWriter, output),
	})
	Expect(err).ToNot(HaveOccurred())

	code, err := process.Wait()
	Expect(err).ToNot(HaveOccurred())
	Expect(code).To(Equal(0))

	return output
}

func getMountType(container garden.Container, pathToSearchFor string) string {
	process, output := containerReadFile(container, "/proc/self/mounts", "root")
	Expect(process.Wait()).To(Equal(0))

	var mountType string
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		mountLine := scanner.Text()
		mountInfo := strings.Split(mountLine, " ")
		mountDest := mountInfo[1]
		if mountDest == pathToSearchFor {
			mountType = mountInfo[2]
		}
	}

	return mountType
}
