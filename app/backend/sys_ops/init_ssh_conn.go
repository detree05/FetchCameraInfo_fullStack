package sys_ops

import "github.com/melbahja/goph"

func InitSSHConnection(controlHost, sshUser, sshPass string) (*goph.Client, error) {
	client, err := goph.NewUnknown(sshUser, controlHost, goph.Password(sshPass)) // not quite good in terms of security but whatever, every connection will be made in local network
	return client, err
}
