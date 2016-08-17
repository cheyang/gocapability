package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/cheyang/gocapability/capability"
	docker "github.com/fsouza/go-dockerclient"
)

const allCapabilityTypes = capability.CAPS | capability.BOUNDS

var (
	container string
	capStr    string
	pid       int

	localDockerEndpoint = "unix:///var/run/docker.sock"
	capabilityMap       map[string]capability.Cap
)

func init() {
	flag.StringVar(&container, "name", "", "The name of container")
	flag.StringVar(&capStr, "cap-add", "", "Capablities separated by comma, like NET_ADMIN,SYS_ADMIN")
	flag.IntVar(&pid, "pid", 0, "The process id")
	initCapMap()
}

func main() {
	flag.Parse()

	var err error

	if err = validate(); err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		return
	}

	if pid == 0 {
		pid, err = getPidFromContainer(container)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	caps := strings.Split(capStr, ",")

	err = addCaps(pid, caps)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func validate() error {
	if capStr == "" {
		return fmt.Errorf("Please set capablities.")
	}

	if pid == 0 && container == "" {
		return fmt.Errorf("Please set one of the container name and pid.")
	}

	if pid > 0 && container != "" {
		return fmt.Errorf("Please set one of the container name and pid.")
	}

	return nil
}

func getPidFromContainer(name string) (pid int, err error) {
	client, err := docker.NewClient(localDockerEndpoint)

	if err != nil {
		return pid, err
	}

	container, err := client.InspectContainer(name)

	if err != nil {
		return pid, err
	}

	if container.State.Pid == 0 {
		return pid, fmt.Errorf("The container is not running!")
	}

	return pid, nil
}

func addCaps(pid int, caps []string) error {
	l := []capability.Cap{}
	for _, c := range caps {

		c = strings.ToUpper(c)

		if !strings.HasPrefix(c, "CAP_") {
			c = strings.Join([]string{"CAP_", c}, "")
		}

		v, ok := capabilityMap[c]
		if !ok {
			return fmt.Errorf("unknown capability %q", c)
		}
		l = append(l, v)
	}
	process, err := capability.NewPid(pid)
	if err != nil {
		return err
	}

	fmt.Println(process.String())
	fmt.Println(process.StringCap(allCapabilityTypes))
	// process.

	// process.Clear(allCapabilityTypes)
	// process.Set(allCapabilityTypes, l)
	// // pid.Set(which, ...)
	// process.Apply(allCapabilityTypes)

	return nil
}

func initCapMap() {
	capabilityMap = make(map[string]capability.Cap)
	last := capability.CAP_LAST_CAP
	// workaround for RHEL6 which has no /proc/sys/kernel/cap_last_cap
	if last == capability.Cap(63) {
		last = capability.CAP_BLOCK_SUSPEND
	}
	for _, cap := range capability.List() {
		if cap > last {
			continue
		}
		capKey := fmt.Sprintf("CAP_%s", strings.ToUpper(cap.String()))
		capabilityMap[capKey] = cap
	}

	// fmt.Println(capabilityMap)
}
