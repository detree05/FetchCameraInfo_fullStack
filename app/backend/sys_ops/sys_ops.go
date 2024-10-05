package sys_ops

import (
	"fmt"

	"fci-backend.detree05.com/net_ops"
	"github.com/melbahja/goph"
)

func InitSSHConnection(controlHost, sshUser, sshPass string) (*goph.Client, error) {
	client, err := goph.NewUnknown(sshUser, controlHost, goph.Password(sshPass)) // not quite good in terms of security but whatever, every connection will be made in local network
	return client, err
}

func GetChannelSysInfo(client *goph.Client, pathToChannels, channel_id string) (string, string, error) {
	channelConfig, _ := client.Run(fmt.Sprintf("jq . %s/%s/channel.config", pathToChannels, channel_id))
	channelStatus, _ := client.Run(fmt.Sprintf("curl localhost:2032/%s/info | jq .", channel_id))

	channelConfigLink, err := net_ops.PutDataOnPaste(string(channelConfig))
	if err != nil {
		return "", "", err
	}

	channelStatusLink, err := net_ops.PutDataOnPaste(string(channelStatus))
	if err != nil {
		return "", "", err
	}

	return channelConfigLink, channelStatusLink, err
}
