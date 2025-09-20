package container

import (
	"dniprom-cli/internal/model"
	"dniprom-cli/pkg/logger"
)

type Container interface {
	GetLogger() logger.Logger
	GetConfig() *model.Config
}

type container struct {
	logger logger.Logger
	config *model.Config
}

func NewContainer(logger logger.Logger, config *model.Config) Container {
	return &container{
		logger: logger,
		config: config,
	}
}

func (c *container) GetLogger() logger.Logger {
	return c.logger
}

func (c *container) GetConfig() *model.Config {
	return c.config
}
