package common

import (
	"kubernetes-ingress-controller/conf"
	"kubernetes-ingress-controller/utils"
)

type IFactory interface {
	Config() *conf.Config
	File() utils.IFile
	Utils() utils.IUtils
}

type factory struct {
	config *conf.Config
	file   utils.IFile
	utils  utils.IUtils
}

func (f *factory) Config() *conf.Config {
	return f.config
}

func (f *factory) File() utils.IFile {
	return f.file
}

func (f *factory) Utils() utils.IUtils {
	return f.utils
}

var _factory *factory

func init() {
	_factory = new(factory)
	_factory.config = conf.LoadConfig()
	_factory.file = utils.NewFile()
}

// GetFactory 获取工厂
func GetFactory() IFactory {
	return _factory
}
