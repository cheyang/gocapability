package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cheyang/gocapability/capability"
	docker "github.com/fsouza/go-dockerclient"
)

// const allCapabilityTypes = capability.EFFECTIVE | capability.BOUNDS | capability.PERMITTED
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

		fmt.Printf("PID: %d\n", pid)
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
		pid = os.Getpid()
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
	pid = container.State.Pid
	if pid == 0 {
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

	showCapabilities(process, pid)
	showGivenCaps(process, caps)
	process.Set(allCapabilityTypes, l...)
	err = process.Apply(allCapabilityTypes)
	if err != nil {
		return err
	}
	showGivenCaps(process, caps)

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

func showCapabilities(process capability.Capabilities, pid int) {

	fmt.Printf("Process %d was give capabilities:", pid)

	for k, v := range capabilityMap {

		// fmt.Printf("%s: EFFECTIVE-%t,  PERMITTED-%t, INHERITABLE-%t\n",
		// 	k,
		// 	process.Get(capability.EFFECTIVE, v),
		// 	process.Get(capability.PERMITTED, v),
		// 	process.Get(capability.INHERITABLE, v))

		//allCapabilityTypes

		fmt.Printf("%s: %t\n",
			k,
			process.Get(capability.BOUNDS, v))
		// if process.Get(capability.BOUNDS, v) {
		// 	fmt.Printf("%s ", k)
		// }
	}

	fmt.Println("")
}

func showGivenCaps(process capability.Capabilities, caps []string) {

	for _, v := range caps {
		fmt.Printf("%s: %t\n",
			v,
			process.Get(allCapabilityTypes, capabilityMap[v]))
	}
}
