package main

import (
	"strings"
)

type UploadOption struct {
	up     bool   // 是否上传
	upType string // 上传分类

	localPath  string // 本地路径
	remotePath string // 远程路径

	connType   string // ssh连接方式
	sshAddress string // ssh地址
	sshUser    string // ssh用户
	sshPass    string // ssh密码

	cmd     bool   // 是否执行指令
	cmdInfo string // 指令

	filter map[string]struct{} // 过滤
}

type Option func(opts *UploadOption)

func NewUploadOption() *UploadOption {
	o := &UploadOption{
		filter: make(map[string]struct{}),
	}

	return o
}

func WithUp(up bool) Option {
	return func(opts *UploadOption) {
		opts.up = up
	}
}
func WithUpType(upType string) Option {
	return func(opts *UploadOption) {
		opts.upType = upType
	}
}

func WithLocalPath(localPath string) Option {
	return func(opts *UploadOption) {
		opts.localPath = localPath
	}
}

func WithRemotePath(remotePath string) Option {
	return func(opts *UploadOption) {
		opts.remotePath = remotePath
	}
}

func WithSshAddress(sshAddress string) Option {
	return func(opts *UploadOption) {
		opts.sshAddress = sshAddress
	}
}

func WithSshUser(sshUser string) Option {
	return func(opts *UploadOption) {
		opts.sshUser = sshUser
	}
}

func WithSshPass(sshPass string) Option {
	return func(opts *UploadOption) {
		opts.sshPass = sshPass
	}
}

func WithCmd(cmd bool) Option {
	return func(opts *UploadOption) {
		opts.cmd = cmd
	}
}
func WithCmdInfo(cmdInfo string) Option {
	return func(opts *UploadOption) {
		opts.cmdInfo = cmdInfo
	}
}

func WithFilter(filter string) Option {
	return func(opts *UploadOption) {
		filters := strings.Split(filter, ";")
		for _, s := range filters {
			opts.filter[s] = struct{}{}
		}
		//fmt.Println(opts.filter)
	}
}

func WithConnType(connType string) Option {
	return func(opts *UploadOption) {
		opts.connType = connType
	}
}
