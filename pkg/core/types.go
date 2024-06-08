package core

import (
	"crazydocker/pkg/config"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type MachineStatus int

const MachineOnline = 1
const MachineOffline = 0

type MachineSSHConfig struct {
	*config.PasswordAuth
	*config.SSHAuth
	ClientConfig *ssh.ClientConfig `json:"-"`
}

type Machine struct {
	Ip        string
	parsedIp  net.IP
	Status    MachineStatus
	Os        string
	HostName  string
	Shell     string
	Error     string
	SSHConfig *MachineSSHConfig `json:"-"`
}

func (m *Machine) getSSHConn() (*ssh.Client, error) {
	if m.SSHConfig.ClientConfig == nil {
		config := &ssh.ClientConfig{
			User:            m.SSHConfig.PasswordAuth.Username,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{
				ssh.Password(m.SSHConfig.Password),
			},
			Timeout: 5 * time.Second,
		}
		if m.SSHConfig.SSHAuth.PrivateKeyFile != "" {
			pemByes, err := os.ReadFile(fmt.Sprintf("%s/%s", os.Getenv("CONFIG_FOLDER"), m.SSHConfig.SSHAuth.PrivateKeyFile))
			if err != nil {
				return nil, err
			}
			signer, err := ssh.ParsePrivateKey(pemByes)
			if err != nil {
				return nil, err
			}
			config.User = m.SSHConfig.SSHAuth.Username
			config.Auth = []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			}
		}
		m.SSHConfig.ClientConfig = config
	}
	return ssh.Dial("tcp", net.JoinHostPort(m.Ip, "22"), m.SSHConfig.ClientConfig)
}

func (m *Machine) RunCommand(cmd string) (string, error) {
	conn, err := m.getSSHConn()
	if err != nil {
		m.Status = MachineOffline
		return "", err
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	out, err := session.CombinedOutput(cmd)
	return strings.TrimSpace(string(out)), err
}

func (m *Machine) StreamCommand(writeConn *websocket.Conn, cmd string) error {
	conn, err := m.getSSHConn()
	if err != nil {
		writeConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"err:\"%s\"}", err.Error())))
		m.Status = MachineOffline
		return err
	}
	go func() {
		defer writeConn.Close()
		ticker := time.NewTicker(5 * time.Second)
		for {
			<-ticker.C
			session, err := conn.NewSession()
			if err != nil {
				writeConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"err:\"%s\"}", err.Error())))
				ticker.Stop()
				return
			}
			modes := ssh.TerminalModes{
				ssh.ECHO:          0,
				ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_OSPEED: 14400,
			}
			err = session.RequestPty("xterm", 80, 40, modes)
			if err != nil {
				writeConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"err:\"%s\"}", err.Error())))
				session.Close()
				ticker.Stop()
				return
			}
			err = writeConn.WriteMessage(websocket.PingMessage, []byte(""))
			if err != nil {
				ticker.Stop()
				session.Close()
				return
			}
			output, err := session.CombinedOutput(cmd)
			if err != nil {
				writeConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"err:\"%s\"}", err.Error())))
			} else {
				writeConn.WriteMessage(websocket.TextMessage, output)
			}
			session.Close()
		}
	}()
	return nil
}

func (m *Machine) GetShell(writeConn *websocket.Conn) error {
	return m.ExecCommand(writeConn, m.Shell)
}

func (m *Machine) ExecCommand(writeConn *websocket.Conn, cmd string) error {
	defer writeConn.Close()
	// use pty here as it combines stdout and error
	var command *exec.Cmd
	if m.SSHConfig.Password != "" {
		command = exec.Command("sshpass", "-p", "password", "ssh", "-tt", "-o StrictHostKeyChecking=no", "-o UserKnownHostsFile=/dev/null", fmt.Sprintf("%s@%s", m.SSHConfig.PasswordAuth.Username, m.Ip), cmd)
	} else {
		command = exec.Command("ssh", "-tt", "-o StrictHostKeyChecking=no", "-o UserKnownHostsFile=/dev/null", fmt.Sprintf("%s@%s", m.SSHConfig.PasswordAuth.Username, m.Ip), cmd)
	}
	ptmx, err := pty.Start(command)
	if err != nil {
		writeConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"err:\"%s\"}", err.Error())))
		return err
	}
	defer ptmx.Close()
	// go routine to read message
	go func() {
		for {
			defer ptmx.Close()
			_, data, err := writeConn.ReadMessage()
			if err != nil {
				log.Printf("reader err :%s\n", err)
				return
			}
			_, err = ptmx.Write(data)
			if err != nil {
				log.Printf("ptmx write err :%s\n", err)
				return
			}
		}
	}()
	// go routine to write message
	go func() {
		defer ptmx.Close()
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("ptmx read err :%s\n", err)
					return
				}
			}
			if err := writeConn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
				log.Printf("writer err :%s\n", err)
				return
			}
		}
	}()
	command.Wait()
	return nil
}

type Machines map[string]*Machine

type Container struct {
	Command      string `json:"Command"`
	CreatedAt    string `json:"CreatedAt"`
	ID           string `json:"ID"`
	Image        string `json:"Image"`
	Labels       string `json:"Labels"`
	LocalVolumes string `json:"LocalVolumes"`
	Mounts       string `json:"Mounts"`
	Names        string `json:"Names"`
	Networks     string `json:"Networks"`
	Ports        string `json:"Ports"`
	RunningFor   string `json:"RunningFor"`
	Size         string `json:"Size"`
	State        string `json:"State"`
	Status       string `json:"Status"`
}

type Containers []*Container

type Image struct {
	Containers   string `json:"Containers"`
	CreatedAt    string `json:"CreatedAt"`
	CreatedSince string `json:"CreatedSince"`
	Digest       string `json:"Digest"`
	ID           string `json:"ID"`
	Repository   string `json:"Repository"`
	SharedSize   string `json:"SharedSize"`
	Size         string `json:"Size"`
	Tag          string `json:"Tag"`
	UniqueSize   string `json:"UniqueSize"`
	VirtualSize  string `json:"VirtualSize"`
}

type Images []*Image

type AutoGenerated []struct {
	ID      string    `yaml:"Id"`
	Created time.Time `yaml:"Created"`
	Path    string    `yaml:"Path"`
	Args    []string  `yaml:"Args"`
	State   struct {
		Status     string    `yaml:"Status"`
		Running    bool      `yaml:"Running"`
		Paused     bool      `yaml:"Paused"`
		Restarting bool      `yaml:"Restarting"`
		OOMKilled  bool      `yaml:"OOMKilled"`
		Dead       bool      `yaml:"Dead"`
		Pid        int       `yaml:"Pid"`
		ExitCode   int       `yaml:"ExitCode"`
		Error      string    `yaml:"Error"`
		StartedAt  time.Time `yaml:"StartedAt"`
		FinishedAt time.Time `yaml:"FinishedAt"`
	} `yaml:"State"`
	Image           string      `yaml:"Image"`
	ResolvConfPath  string      `yaml:"ResolvConfPath"`
	HostnamePath    string      `yaml:"HostnamePath"`
	HostsPath       string      `yaml:"HostsPath"`
	LogPath         string      `yaml:"LogPath"`
	Name            string      `yaml:"Name"`
	RestartCount    int         `yaml:"RestartCount"`
	Driver          string      `yaml:"Driver"`
	Platform        string      `yaml:"Platform"`
	MountLabel      string      `yaml:"MountLabel"`
	ProcessLabel    string      `yaml:"ProcessLabel"`
	AppArmorProfile string      `yaml:"AppArmorProfile"`
	ExecIDs         interface{} `yaml:"ExecIDs"`
	HostConfig      struct {
		Binds           []string `yaml:"Binds"`
		ContainerIDFile string   `yaml:"ContainerIDFile"`
		LogConfig       struct {
			Type   string `yaml:"Type"`
			Config struct {
			} `yaml:"Config"`
		} `yaml:"LogConfig"`
		NetworkMode  string `yaml:"NetworkMode"`
		PortBindings struct {
			One0000TCP []struct {
				HostIP   string `yaml:"HostIp"`
				HostPort string `yaml:"HostPort"`
			} `yaml:"10000/tcp"`
			Nine901TCP []struct {
				HostIP   string `yaml:"HostIp"`
				HostPort string `yaml:"HostPort"`
			} `yaml:"9901/tcp"`
		} `yaml:"PortBindings"`
		RestartPolicy struct {
			Name              string `yaml:"Name"`
			MaximumRetryCount int    `yaml:"MaximumRetryCount"`
		} `yaml:"RestartPolicy"`
		AutoRemove           bool          `yaml:"AutoRemove"`
		VolumeDriver         string        `yaml:"VolumeDriver"`
		VolumesFrom          interface{}   `yaml:"VolumesFrom"`
		ConsoleSize          []int         `yaml:"ConsoleSize"`
		CapAdd               interface{}   `yaml:"CapAdd"`
		CapDrop              interface{}   `yaml:"CapDrop"`
		CgroupnsMode         string        `yaml:"CgroupnsMode"`
		DNS                  []interface{} `yaml:"Dns"`
		DNSOptions           []interface{} `yaml:"DnsOptions"`
		DNSSearch            []interface{} `yaml:"DnsSearch"`
		ExtraHosts           interface{}   `yaml:"ExtraHosts"`
		GroupAdd             interface{}   `yaml:"GroupAdd"`
		IpcMode              string        `yaml:"IpcMode"`
		Cgroup               string        `yaml:"Cgroup"`
		Links                interface{}   `yaml:"Links"`
		OomScoreAdj          int           `yaml:"OomScoreAdj"`
		PidMode              string        `yaml:"PidMode"`
		Privileged           bool          `yaml:"Privileged"`
		PublishAllPorts      bool          `yaml:"PublishAllPorts"`
		ReadonlyRootfs       bool          `yaml:"ReadonlyRootfs"`
		SecurityOpt          interface{}   `yaml:"SecurityOpt"`
		UTSMode              string        `yaml:"UTSMode"`
		UsernsMode           string        `yaml:"UsernsMode"`
		ShmSize              int           `yaml:"ShmSize"`
		Runtime              string        `yaml:"Runtime"`
		Isolation            string        `yaml:"Isolation"`
		CPUShares            int           `yaml:"CpuShares"`
		Memory               int           `yaml:"Memory"`
		NanoCpus             int           `yaml:"NanoCpus"`
		CgroupParent         string        `yaml:"CgroupParent"`
		BlkioWeight          int           `yaml:"BlkioWeight"`
		BlkioWeightDevice    []interface{} `yaml:"BlkioWeightDevice"`
		BlkioDeviceReadBps   []interface{} `yaml:"BlkioDeviceReadBps"`
		BlkioDeviceWriteBps  []interface{} `yaml:"BlkioDeviceWriteBps"`
		BlkioDeviceReadIOps  []interface{} `yaml:"BlkioDeviceReadIOps"`
		BlkioDeviceWriteIOps []interface{} `yaml:"BlkioDeviceWriteIOps"`
		CPUPeriod            int           `yaml:"CpuPeriod"`
		CPUQuota             int           `yaml:"CpuQuota"`
		CPURealtimePeriod    int           `yaml:"CpuRealtimePeriod"`
		CPURealtimeRuntime   int           `yaml:"CpuRealtimeRuntime"`
		CpusetCpus           string        `yaml:"CpusetCpus"`
		CpusetMems           string        `yaml:"CpusetMems"`
		Devices              []interface{} `yaml:"Devices"`
		DeviceCgroupRules    interface{}   `yaml:"DeviceCgroupRules"`
		DeviceRequests       interface{}   `yaml:"DeviceRequests"`
		MemoryReservation    int           `yaml:"MemoryReservation"`
		MemorySwap           int           `yaml:"MemorySwap"`
		MemorySwappiness     interface{}   `yaml:"MemorySwappiness"`
		OomKillDisable       bool          `yaml:"OomKillDisable"`
		PidsLimit            interface{}   `yaml:"PidsLimit"`
		Ulimits              interface{}   `yaml:"Ulimits"`
		CPUCount             int           `yaml:"CpuCount"`
		CPUPercent           int           `yaml:"CpuPercent"`
		IOMaximumIOps        int           `yaml:"IOMaximumIOps"`
		IOMaximumBandwidth   int           `yaml:"IOMaximumBandwidth"`
		MaskedPaths          []string      `yaml:"MaskedPaths"`
		ReadonlyPaths        []string      `yaml:"ReadonlyPaths"`
	} `yaml:"HostConfig"`
	GraphDriver struct {
		Data struct {
			LowerDir  string `yaml:"LowerDir"`
			MergedDir string `yaml:"MergedDir"`
			UpperDir  string `yaml:"UpperDir"`
			WorkDir   string `yaml:"WorkDir"`
		} `yaml:"Data"`
		Name string `yaml:"Name"`
	} `yaml:"GraphDriver"`
	Mounts []struct {
		Type        string `yaml:"Type"`
		Source      string `yaml:"Source"`
		Destination string `yaml:"Destination"`
		Mode        string `yaml:"Mode"`
		RW          bool   `yaml:"RW"`
		Propagation string `yaml:"Propagation"`
	} `yaml:"Mounts"`
	Config struct {
		Hostname     string `yaml:"Hostname"`
		Domainname   string `yaml:"Domainname"`
		User         string `yaml:"User"`
		AttachStdin  bool   `yaml:"AttachStdin"`
		AttachStdout bool   `yaml:"AttachStdout"`
		AttachStderr bool   `yaml:"AttachStderr"`
		ExposedPorts struct {
			One0000TCP struct {
			} `yaml:"10000/tcp"`
			Nine901TCP struct {
			} `yaml:"9901/tcp"`
		} `yaml:"ExposedPorts"`
		Tty        bool        `yaml:"Tty"`
		OpenStdin  bool        `yaml:"OpenStdin"`
		StdinOnce  bool        `yaml:"StdinOnce"`
		Env        []string    `yaml:"Env"`
		Cmd        []string    `yaml:"Cmd"`
		Image      string      `yaml:"Image"`
		Volumes    interface{} `yaml:"Volumes"`
		WorkingDir string      `yaml:"WorkingDir"`
		Entrypoint []string    `yaml:"Entrypoint"`
		OnBuild    interface{} `yaml:"OnBuild"`
		Labels     struct {
			OrgOpencontainersImageRefName string `yaml:"org.opencontainers.image.ref.name"`
			OrgOpencontainersImageVersion string `yaml:"org.opencontainers.image.version"`
		} `yaml:"Labels"`
	} `yaml:"Config"`
	NetworkSettings struct {
		Bridge                 string `yaml:"Bridge"`
		SandboxID              string `yaml:"SandboxID"`
		HairpinMode            bool   `yaml:"HairpinMode"`
		LinkLocalIPv6Address   string `yaml:"LinkLocalIPv6Address"`
		LinkLocalIPv6PrefixLen int    `yaml:"LinkLocalIPv6PrefixLen"`
		Ports                  struct {
			One0000TCP []struct {
				HostIP   string `yaml:"HostIp"`
				HostPort string `yaml:"HostPort"`
			} `yaml:"10000/tcp"`
			Nine901TCP []struct {
				HostIP   string `yaml:"HostIp"`
				HostPort string `yaml:"HostPort"`
			} `yaml:"9901/tcp"`
		} `yaml:"Ports"`
		SandboxKey             string      `yaml:"SandboxKey"`
		SecondaryIPAddresses   interface{} `yaml:"SecondaryIPAddresses"`
		SecondaryIPv6Addresses interface{} `yaml:"SecondaryIPv6Addresses"`
		EndpointID             string      `yaml:"EndpointID"`
		Gateway                string      `yaml:"Gateway"`
		GlobalIPv6Address      string      `yaml:"GlobalIPv6Address"`
		GlobalIPv6PrefixLen    int         `yaml:"GlobalIPv6PrefixLen"`
		IPAddress              string      `yaml:"IPAddress"`
		IPPrefixLen            int         `yaml:"IPPrefixLen"`
		IPv6Gateway            string      `yaml:"IPv6Gateway"`
		MacAddress             string      `yaml:"MacAddress"`
		Networks               struct {
			Bridge struct {
				IPAMConfig          interface{} `yaml:"IPAMConfig"`
				Links               interface{} `yaml:"Links"`
				Aliases             interface{} `yaml:"Aliases"`
				NetworkID           string      `yaml:"NetworkID"`
				EndpointID          string      `yaml:"EndpointID"`
				Gateway             string      `yaml:"Gateway"`
				IPAddress           string      `yaml:"IPAddress"`
				IPPrefixLen         int         `yaml:"IPPrefixLen"`
				IPv6Gateway         string      `yaml:"IPv6Gateway"`
				GlobalIPv6Address   string      `yaml:"GlobalIPv6Address"`
				GlobalIPv6PrefixLen int         `yaml:"GlobalIPv6PrefixLen"`
				MacAddress          string      `yaml:"MacAddress"`
				DriverOpts          interface{} `yaml:"DriverOpts"`
			} `yaml:"bridge"`
		} `yaml:"Networks"`
	} `yaml:"NetworkSettings"`
}
