package containers

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultVersion = "v0.5.17"
	flagPort       = "8013"
)

type flagdContainer struct {
	testcontainers.Container
	port int
}

func (fc *flagdContainer) GetPort() int {
	return fc.port

}

type flagdConfig struct {
	unstable bool
	version  string
}

type FlagdContainerOption func(*flagdConfig)

func WithUnstable(unstable bool) FlagdContainerOption {
	return func(fc *flagdConfig) {
		fc.unstable = unstable
	}
}
func WithVersion(version string) FlagdContainerOption {
	return func(fc *flagdConfig) {
		fc.version = version
	}
}

func NewFlagd(ctx context.Context, opts ...FlagdContainerOption) (*flagdContainer, error) {
	c := &flagdConfig{
		unstable: false,
		version:  defaultVersion,
	}
	for _, opt := range opts {
		opt(c)
	}

	return setupContainer(ctx, c)
}

func setupContainer(ctx context.Context, cfg *flagdConfig) (*flagdContainer, error) {
	registry := "ghcr.io/open-feature"
	imgName := "flagd-testbed"

	if cfg.unstable {
		imgName += "-unstable"
	}
	fullImgName := registry + "/" + imgName + ":" + cfg.version

	req := testcontainers.ContainerRequest{
		Image:        fullImgName,
		Name:         imgName,
		ExposedPorts: []string{flagPort + "/tcp"},
		Networks:     []string{"bridge", "flagd"},
		WaitingFor:   wait.ForExposedPort(),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})

	if err != nil {
		return nil, err
	}

	mappedPort, err := c.MappedPort(ctx, flagPort)
	if err != nil {
		return nil, err
	}
	return &flagdContainer{
		Container: c,
		port:      mappedPort.Int()}, nil
}
