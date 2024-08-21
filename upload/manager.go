package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

var mgr *UploadMgr

type UploadMgr struct {
	options *UploadOption

	// ssh
	ssh     *ssh.Client
	sftp    *sftp.Client
	session *ssh.Session
}

func init() {
	mgr = &UploadMgr{
		options: NewUploadOption(),
	}
}

func Init(opts ...Option) (err error) {
	for _, opt := range opts {
		opt(mgr.options)
	}

	flag := false
	for i := 0; i < 3; i++ {
		if err = mgr.connSSH(); err != nil {
			fmt.Println("ssh init error. try again")
			time.Sleep(500 * time.Millisecond)
			continue
		}
		flag = true
		break
	}
	if !flag {
		return
	}

	return
}

func (mgr *UploadMgr) connSSH() (err error) {
	// ssh config
	config := &ssh.ClientConfig{
		User:            mgr.options.sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	switch mgr.options.connType {
	case CONN_TYPE_PASSWORD:
		config.Auth = []ssh.AuthMethod{ssh.Password(mgr.options.sshPass)}
	case CONN_TYPE_PUBLIC_KEY:
		// 读取私钥
		privateKeyBytes, err := os.ReadFile(mgr.options.sshPass)
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}
		// 解析私钥
		privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(privateKey)}
	}

	// connection
	mgr.ssh, err = ssh.Dial("tcp", mgr.options.sshAddress, config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}

	// create session
	mgr.session, err = mgr.ssh.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %v", err)
	}

	// 创建一个SFTP会话
	mgr.sftp, err = sftp.NewClient(mgr.ssh)
	if err != nil {
		log.Fatalf("unable to create sftp client: %v", err)
	}

	return
}

func (mgr *UploadMgr) run() {
	var err error
	if mgr.options.up {
		if err = mgr.uploadDir(mgr.options.localPath, mgr.options.remotePath); err != nil {
			log.Fatalf("upload dir err: %v", err)
		}
	}

	if mgr.options.cmd {
		mgr.CmdDo()
	}

	mgr.ssh.Close()
	mgr.sftp.Close()
	mgr.session.Close()
}
