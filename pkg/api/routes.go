package api

import (
	"crazydocker/pkg/core"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var ce *core.CommandExecutor

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func listMachines(c *gin.Context) {
	machines := ce.ListMachines()
	c.JSON(200, &MachineListResponse{Machines: machines})
}

func listContainers(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	if ip == "" {
		c.JSON(500, &ContainerListResponse{Containers: nil, Response: Response{Error: "Please provide valid ip"}})
		return
	}
	containers, err := ce.ListContainers(ip)
	if err != nil {
		c.JSON(500, &ContainerListResponse{Containers: nil, Response: Response{Error: err.Error()}})
		return
	}
	c.JSON(200, &ContainerListResponse{Containers: containers})
}

func performActionOnContainer(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	containerId := c.Request.URL.Query().Get("containerID")
	action := c.Request.URL.Query().Get("action")
	if ip == "" || containerId == "" || action == "" {
		c.JSON(500, &Response{Error: "Please provide both ip and containerID"})
		return
	}
	_, err := ce.PerformAction(ip, containerId, action)
	if err != nil {
		c.JSON(500, &Response{Error: err.Error()})
		return
	}
	c.JSON(200, &Response{Error: ""})
}

func streamContainer(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	containerId := c.Request.URL.Query().Get("containerID")
	if ip == "" || containerId == "" {
		c.JSON(500, &ContainerResponse{Container: nil, Response: Response{Error: "Please provide both ip and containerID"}})
		return
	}
	// upgrade to websocket and start streaming container info
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, &Response{Error: "unable to upgrade the connection to ws"})
		return
	}
	err = ce.StreamContainer(conn, ip, containerId)
	if err != nil {
		c.JSON(500, &Response{Error: err.Error()})
		return
	}
}

func streamImage(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	imageId := c.Request.URL.Query().Get("imageID")
	if ip == "" || imageId == "" {
		c.JSON(500, &ImageResponse{Image: nil, Response: Response{Error: "Please provide both ip and containerID"}})
		return
	}
	// upgrade to websocket and start streaming container info
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, &Response{Error: "unable to upgrade the connection to ws"})
		return
	}
	err = ce.StreamImage(conn, ip, imageId)
	if err != nil {
		c.JSON(500, &Response{Error: err.Error()})
		return
	}
}

func listImages(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	if ip == "" {
		c.JSON(500, &ImageListResponse{Images: nil, Response: Response{Error: "Please provide valid ip"}})
		return
	}
	Images, err := ce.ListImages(ip)
	if err != nil {
		c.JSON(500, &ImageListResponse{Images: nil, Response: Response{Error: err.Error()}})
		return
	}

	c.JSON(200, &ImageListResponse{Images: Images})
}

func getConfig(c *gin.Context) {
	data, err := os.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		c.JSON(500, &Response{Error: err.Error()})
		return
	}
	s := &strings.Builder{}
	quick.Highlight(s, string(data), "yaml", "html", "dracula")
	c.JSON(200, &ConfigResponse{HTML: s.String(), Raw: string(data)})
}

func createContainer(c *gin.Context) {
	var payload *CreateContainerPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(500, &Response{Error: err.Error()})
		return
	}
	out, err := ce.CreateContainer(payload.Ip, payload.Image, payload.Args)
	if err != nil {
		fmt.Println()
		c.JSON(500, &CreateContainerResponse{Msg: out, Response: Response{Error: err.Error()}})
		return
	}
	c.JSON(200, &CreateContainerResponse{Msg: out})
}

func execIntoMachine(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	if ip == "" {
		c.JSON(500, &Response{Error: "please provide valid ip"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, &Response{Error: "unable to upgrade the connection to ws"})
		return
	}
	// get the machine
	if err := ce.ExecIntoMachine(conn, ip); err != nil {
		c.JSON(500, &Response{Error: fmt.Sprintf("got error %s", err)})
		return
	}
}

func execIntoContainer(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	containerID := c.Request.URL.Query().Get("containerID")
	if ip == "" || containerID == "" {
		c.JSON(500, &Response{Error: "please provide valid ip and container id"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, &Response{Error: "unable to upgrade the connection to ws"})
		return
	}
	if err := ce.ExecIntoContainer(conn, ip, containerID); err != nil {
		c.JSON(500, &Response{Error: fmt.Sprintf("got error %s", err)})
		return
	}
}

func streamContainerLogs(c *gin.Context) {
	ip := c.Request.URL.Query().Get("ip")
	containerID := c.Request.URL.Query().Get("containerID")
	if ip == "" || containerID == "" {
		c.JSON(500, &Response{Error: "please provide valid ip and container id"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, &Response{Error: "unable to upgrade the connection to ws"})
		return
	}
	if err := ce.StreamContainerLogs(conn, ip, containerID); err != nil {
		c.JSON(500, &Response{Error: fmt.Sprintf("got error %s", err)})
		return
	}
}
