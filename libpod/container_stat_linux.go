//go:build !remote
// +build !remote

package libpod

import (
	"fmt"

	"github.com/containers/buildah/copier"
	"github.com/containers/podman/v4/libpod/define"
)

// statInsideMount stats the specified path *inside* the container's mount and PID
// namespace.  It returns the file info along with the resolved root ("/") and
// the resolved path (relative to the root).
func (c *Container) statInsideMount(containerPath string) (*copier.StatForItem, string, string, error) {
	fmt.Printf("statInsideMount\n")
	resolvedRoot := "/"
	resolvedPath := c.pathAbs(containerPath)
	var statInfo *copier.StatForItem

	err := c.joinMountAndExec(
		func() error {
			fmt.Printf("statInsideMount joinMountAndExec\n")
			var statErr error
			fmt.Printf("enter secureStat\n")
			statInfo, statErr = secureStat(resolvedRoot, resolvedPath)
			fmt.Printf("exit secureStat\n")
			return statErr
		},
	)

	return statInfo, resolvedRoot, resolvedPath, err
}

// Calls either statOnHost or statInsideMount depending on whether the
// container is running
func (c *Container) statInContainer(mountPoint string, containerPath string) (*copier.StatForItem, string, string, error) {
	fmt.Printf("statInContainer\n")
	if c.state.State == define.ContainerStateRunning {
		// If the container is running, we need to join it's mount namespace
		// and stat there.
		fmt.Printf("enter statInsideMount\n")
		return c.statInsideMount(containerPath)
	}
	// If the container is NOT running, we need to resolve the path
	// on the host.
	fmt.Printf("enter statOnHost\n")
	return c.statOnHost(mountPoint, containerPath)
}
