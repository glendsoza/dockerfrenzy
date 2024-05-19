package api

import (
	"crazydocker/pkg/core"

	"docker.io/go-docker/api/types"
)

type Response struct {
	Error string
}

type MachineListResponse struct {
	Response
	Machines []*core.Machine
}

type ContainerListResponse struct {
	Response
	Containers core.Containers
}

type ImageListResponse struct {
	Response
	Images core.Images
}

type ContainerResponse struct {
	Response
	Container *types.ContainerJSON
}

type ImageResponse struct {
	Response
	Image *types.ImageInspect
}

type ConfigResponse struct {
	Response
	HTML string
	Raw  string
}

type CreateContainerPayload struct {
	Ip    string
	Image string
	Args  string
}

type CreateContainerResponse struct {
	Response
	Msg string
}
