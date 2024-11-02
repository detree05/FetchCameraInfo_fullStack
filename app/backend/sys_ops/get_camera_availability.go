package sys_ops

import (
	"fmt"

	"github.com/melbahja/goph"
)

func GetCameraAvailability(client *goph.Client, pathToChannels, channel_id string) (string, error) {
	cameraIp, _ := client.Run(fmt.Sprintf("jq '.url' %s/%s/channel.config | grep -o '[0-9]\\+[.][0-9]\\+[.][0-9]\\+[.][0-9]\\+'", pathToChannels, channel_id))
	availabilityStatus, _ := client.Run(fmt.Sprintf("bash -c 'ping -c 1 %s' | grep loss", cameraIp))

	return string(availabilityStatus), nil
}
