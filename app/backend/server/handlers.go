package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"fci-backend.detree05.com/cfg"
	"fci-backend.detree05.com/sql_ops"
	"fci-backend.detree05.com/sys_ops"
	"github.com/melbahja/goph"
)

type GetChannelRequest struct {
	Verbose     bool
	ProjectName string
	CameraName  string
	Username    string
	Password    string
}

type GetChannelResponse struct {
	Name               string
	ExtId              string
	ChannelId          string
	ControlHost        string
	CameraAvailability string
	ChannelConfigLink  string
	ChannelStatusLink  string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", cfg.Config.CorsPolicy.Allow)
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
}

func genVerboseResponse(locLogLine uint, controlHost string, gcRequest GetChannelRequest, channels *map[string][]sql_ops.ChannelInfo, gcResponse *[]GetChannelResponse) {
	log.Printf("[%d] Opening new connection to %s", locLogLine, controlHost)

	var (
		pathToChannels     string
		channelConfigLink  string
		channelStatusLink  string
		cameraAvailability string
		sshStatus          bool
	)

	client, err := sys_ops.InitSSHConnection(controlHost, gcRequest.Username, gcRequest.Password)
	if err != nil {
		log.Printf("[%d] Couldn't connect to %s!", locLogLine, controlHost)
		sshStatus = false
		channelConfigLink, channelStatusLink = "Err", "Err"
	} else {
		log.Printf("[%d] Successfully connected to %s", locLogLine, controlHost)
		sshStatus = true
		// var pathToChannels is needed to speed up the fetching data process, since original method implemented searching channel directory in global /opt/netris dir every time GetChannelSysInfo is called
		// now it searches in channel directory directly (e.g. /opt/netris/storage-meta) which is way faster since there is no need to look around in global directory
		// usually getting this variable takes 2-3 seconds, although all next GetChannelSysInfo calls will take only 600-900 ms on average
		pathToChannelsRaw, _ := client.Run(fmt.Sprintf("dirname $(find /opt/netris -type d -name %s -print -quit 2>/dev/null)", (*channels)[controlHost][0].ChannelId)) // also it turns out that this command output returns []bytes with trailing \n
		pathToChannels = strings.TrimSuffix(string(pathToChannelsRaw), "\n")                                                                                            // obliged to remove that little bastard so it works correctly
	}

	var innerWait sync.WaitGroup

	for _, channel := range (*channels)[controlHost] {
		innerWait.Add(1)

		time.Sleep(250 * time.Millisecond) // PREVENTING SSH FROM EXPLODING

		go func(client *goph.Client, pathToChannels string, channel sql_ops.ChannelInfo, gcResponse *[]GetChannelResponse) {
			defer innerWait.Done()

			if sshStatus {
				log.Printf("[%d] Fetching data on %s and creating links...", locLogLine, channel.ChannelId)

				channelConfigLink, channelStatusLink, err = sys_ops.GetChannelSysInfo(client, pathToChannels, channel.ChannelId)
				if err != nil {
					log.Printf("[%d] Couldn't fetch data on %s with error: %s", locLogLine, channel.ChannelId, err)
				}

				cameraAvailability, err = sys_ops.GetCameraAvailability(client, pathToChannels, channel.ChannelId)
				if err != nil {
					log.Printf("[%d] Couldn't check camera availability on %s with error: %s", locLogLine, channel.ChannelId, err)
				}
			}

			gcResponseObj := GetChannelResponse{
				Name:               channel.Name,
				ExtId:              channel.ExtId,
				ChannelId:          channel.ChannelId,
				ControlHost:        channel.ControlHost,
				CameraAvailability: cameraAvailability,
				ChannelConfigLink:  channelConfigLink,
				ChannelStatusLink:  channelStatusLink,
			}

			*gcResponse = append((*gcResponse), gcResponseObj)
		}(client, pathToChannels, channel, gcResponse)
	}

	innerWait.Wait()

	if sshStatus {
		log.Printf("[%d] Closing connection to %s", locLogLine, controlHost)
		client.Close()
	}
}

func genNonVerboseResponse(controlHost string, channels *map[string][]sql_ops.ChannelInfo, gcResponse *[]GetChannelResponse) {
	for _, channel := range (*channels)[controlHost] {
		gcResponseObj := GetChannelResponse{
			Name:        channel.Name,
			ExtId:       channel.ExtId,
			ChannelId:   channel.ChannelId,
			ControlHost: channel.ControlHost,
		}
		*gcResponse = append((*gcResponse), gcResponseObj)
	}
}

func getChannel(w http.ResponseWriter, r *http.Request) {
	var gcRequest GetChannelRequest

	start := time.Now()   // for checking time of query execution
	locLogLine := LogLine // for checking requests as req-ans dict
	LogLine += 1

	log.Printf("[%d] Got request! >> Handler: getChannel [Method: %s] [URL: %s] [Protocol: %s] [Content-Type: %s]",
		locLogLine, r.Method, r.URL, r.Proto, r.Header.Get("Content-Type")) // basic logging

	enableCors(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&gcRequest)
	if err != nil {
		log.Printf("[%d] Return 400 ERR! Couldn't decode JSON! %s", locLogLine, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if gcRequest.CameraName == "" {
		log.Printf("[%d] Return 400 ERR! Camera name field is empty!", locLogLine)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("[%d] Getting channels info in sql_ops block.", locLogLine)
	channels, channelsCount, err := sql_ops.GetChannelInfo(gcRequest.ProjectName, gcRequest.CameraName)
	if err != nil {
		log.Printf("[%d] Return 400 ERR! Couldn't get channels info! %s", locLogLine, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// returns map with control hosts as keys and slices with channels as values
	// where 'name' - camera name; 'ext_id', 'channel_id' - respectively to their names; 'control_host' - IP of a recorder where camera is located

	if cfg.Config.Server.Verbose {
		log.Printf("[%d] channels: %s", locLogLine, channels)
	}

	var gcResponse []GetChannelResponse
	var waitGroup sync.WaitGroup

	if gcRequest.Verbose {
		log.Printf("[%d] Verbose answer enabled.", locLogLine)

		if channelsCount > 100 {
			log.Printf("[%d] Achtung! Return 400 ERR! More than 100 channels with enabled verbose!", locLogLine)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for controlHost := range channels {
			waitGroup.Add(1)

			go func(channels *map[string][]sql_ops.ChannelInfo) {
				defer waitGroup.Done()
				genVerboseResponse(locLogLine, controlHost, gcRequest, channels, &gcResponse)
			}(&channels)
		}
	} else {
		log.Printf("[%d] Verbose answer disabled.", locLogLine)

		for controlHost := range channels {
			waitGroup.Add(1)

			go func(channels *map[string][]sql_ops.ChannelInfo) {
				defer waitGroup.Done()
				genNonVerboseResponse(controlHost, channels, &gcResponse)
			}(&channels)
		}
	}

	waitGroup.Wait()

	jsonResponse, err := json.Marshal(gcResponse)
	if err != nil {
		log.Printf("[%d] Return 400 ERR! Couldn't form JSON response! %s", locLogLine, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if cfg.Config.Server.Verbose {
		log.Printf("[%d] jsonResponse: %s", locLogLine, string(jsonResponse))
	}

	log.Printf("[%d] Return 200 OK! Sent request back to sender in %s", locLogLine, time.Since(start)) // basic logging
	w.Write(jsonResponse)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	locLogLine := LogLine
	LogLine += 1

	log.Printf("[%d] Got request! >> Handler: ping [Method: %s] [URL: %s] [Protocol: %s] [Content-Type: %s]",
		locLogLine, r.Method, r.URL, r.Proto, r.Header.Get("Content-Type"))
	log.Printf("[%d] Healthcheck. Return 200 OK!", locLogLine)

	w.WriteHeader(http.StatusOK)
}

func invalidReq(w http.ResponseWriter, r *http.Request) {
	locLogLine := LogLine
	LogLine += 1

	log.Printf("[%d] Got request! >> Handler: invalidReq [Method: %s] [URL: %s] [Protocol: %s] [Content-Type: %s]",
		locLogLine, r.Method, r.URL, r.Proto, r.Header.Get("Content-Type"))
	log.Printf("[%d] Return 400 ERR!", locLogLine)

	w.WriteHeader(http.StatusBadRequest)
}
