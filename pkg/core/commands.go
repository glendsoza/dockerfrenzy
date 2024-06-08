package core

import (
	"crazydocker/pkg/config"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

func NewCommandExecutor() (*CommandExecutor, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	ce := &CommandExecutor{machines: make(Machines), config: config}
	ce.loadMachines()
	return ce, nil
}

type CommandExecutor struct {
	machines Machines
	lock     sync.RWMutex
	config   *config.Config
}

func (ce *CommandExecutor) ReloadConfig() error {
	if err := ce.config.Reload(); err != nil {
		return err
	}
	ce.loadMachines()
	return nil
}

func (ce *CommandExecutor) UpdateConfig(data []byte) error {
	if err := ce.config.Update(data); err != nil {
		return err
	}
	ce.loadMachines()
	return nil
}

func (ce *CommandExecutor) loadMachines() {
	var wg sync.WaitGroup
	maxWorkers := make(chan int, 10)
	ce.clearMachines()
	for _, c := range ce.config.Get() {
		{
			for _, ip := range c.Ips {
				maxWorkers <- 1
				wg.Add(1)
				go func(ip string) {
					defer wg.Done()
					parsedIp := net.ParseIP(ip)
					m := &Machine{Ip: ip, Status: MachineOnline, parsedIp: parsedIp, SSHConfig: &MachineSSHConfig{
						PasswordAuth: &config.PasswordAuth{
							Username: c.PasswordAuth.Username,
							Password: c.Password,
						},
						SSHAuth: &config.SSHAuth{Username: c.SSHAuth.Username, PrivateKeyFile: c.PrivateKeyFile},
					}}
					ce.addMachine(ip, m)
					if parsedIp == nil {
						// see if its valid host
						ips, _ := net.LookupIP(ip)
						if len(ips) > 0 {
							m.parsedIp = ips[0]
						} else {
							m.Status = MachineOffline
							m.Error = fmt.Sprintf("unable to parse %s as valid ip", ip)
							log.Println(m.Error)
							return
						}
					}
					data, err := ce.GetHostNameAndOs(m.Ip)
					if err != nil {
						m.Error = err.Error()
					} else {
						splitData := strings.Split(data, "\n")
						if len(splitData) != 2 {
							m.Error = fmt.Sprintf("unable to determine os or hostname from %s", data)
						} else {
							m.HostName = strings.TrimSpace(splitData[0])
							m.Os = strings.TrimSpace(splitData[1])
						}
					}
					data, err = ce.GetShell(m.Ip)
					if err != nil {
						// combine error
						m.Error = fmt.Sprintf("%s\n%s", err.Error(), m.Error)
					} else {
						m.Shell = strings.TrimSpace(data)
					}
					<-maxWorkers
				}(ip)
			}
		}
		wg.Wait()
	}
}

func (ce *CommandExecutor) getMachines() []*Machine {
	ce.lock.RLock()
	defer ce.lock.RUnlock()
	machines := []*Machine{}
	for _, m := range ce.machines {
		data := *m
		machines = append(machines, &data)
	}
	return machines
}

func (ce *CommandExecutor) getMachine(ip string) *Machine {
	ce.lock.RLock()
	defer ce.lock.RUnlock()
	// panic if the machine is not found for now
	return ce.machines[ip]
}

func (ce *CommandExecutor) addMachine(ip string, m *Machine) {
	ce.lock.Lock()
	defer ce.lock.Unlock()
	ce.machines[ip] = m
}

func (ce *CommandExecutor) clearMachines() {
	ce.lock.Lock()
	defer ce.lock.Unlock()
	ce.machines = make(Machines)
}

func (ce *CommandExecutor) ListMachines() []*Machine {
	return ce.getMachines()
}

func marshalOut[T Container | Image](out string) []*T {
	var dataList []*T
	for _, c := range strings.Split(out, "\n") {
		var data *T
		err := json.Unmarshal([]byte(c), &data)
		if err == nil {
			dataList = append(dataList, data)
		}
	}
	return dataList
}

func (ce *CommandExecutor) StreamContainer(conn *websocket.Conn, ip, containerID string) error {
	return ce.getMachine(ip).StreamCommand(conn, fmt.Sprintf("docker container inspect %s", containerID))
}

func (ce *CommandExecutor) StreamImage(conn *websocket.Conn, ip string, imageId string) error {
	return ce.getMachine(ip).StreamCommand(conn, fmt.Sprintf("docker image inspect %s", imageId))
}

func (ce *CommandExecutor) ListImages(ip string) (Images, error) {
	machine := ce.getMachine(ip)
	out, err := machine.RunCommand(`docker image ls --all --format "{{json . }}" --no-trunc`)
	if err != nil {
		return nil, err
	}
	return marshalOut[Image](out), nil
}

func (ce *CommandExecutor) ListContainers(ip string) (Containers, error) {
	machine := ce.getMachine(ip)
	out, err := machine.RunCommand(`docker container ls --all --format "{{json . }}" --no-trunc`)
	if err != nil {
		return nil, err
	}
	return marshalOut[Container](out), nil
}

func (ce *CommandExecutor) GetHostNameAndOs(ip string) (string, error) {
	return ce.getMachine(ip).RunCommand(`hostnamectl | grep -Ei "Static hostname|Operating System" | cut -f2 -d ":"`)
}

func (ce *CommandExecutor) GetShell(ip string) (string, error) {
	return ce.getMachine(ip).RunCommand("echo $SHELL")
}

func (ce *CommandExecutor) PerformAction(ip string, containerID string, action string) (string, error) {
	return ce.getMachine(ip).RunCommand(fmt.Sprintf("docker %s %s", action, containerID))
}

func (ce *CommandExecutor) CreateContainer(ip string, image string, args string) (string, error) {
	return ce.getMachine(ip).RunCommand(fmt.Sprintf(`docker run -d %s $(docker inspect %s --format "{{ index .RepoTags 0 }}")`, args, image))
}

func (ce *CommandExecutor) ExecIntoMachine(conn *websocket.Conn, ip string) error {
	return ce.getMachine(ip).GetShell(conn)
}

func (ce *CommandExecutor) ExecIntoContainer(conn *websocket.Conn, ip string, containerID string) error {
	return ce.getMachine(ip).ExecCommand(conn, fmt.Sprintf("docker exec -it %s sh", containerID))
}

func (ce *CommandExecutor) StreamContainerLogs(conn *websocket.Conn, ip string, containerID string) error {
	return ce.getMachine(ip).ExecCommand(conn, fmt.Sprintf("docker logs --follow %s", containerID))
}
