package sys_ops

import (
	"fmt"

	"fci-backend.detree05.com/net_ops"
	"github.com/melbahja/goph"
)

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
