package sender

import (
	"aiStudio/internal/conf"
	_const "aiStudio/internal/const"
	"aiStudio/internal/redis"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"path/filepath"
)

const (
	url       = "https://discord.com/api/v9/interactions"
	uploadUrl = "https://discord.com/api/v9/channels/@@channel_id@@/attachments"
)

var (
	confs []conf.Midj
)

func Init(c []conf.Midj) {
	confs = c
	go handle(context.Background()) // 队列
}

func GenerateImage(jobId, prompt string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       "@@guild_id@@",
		ChannelID:     "@@channel_id@@",
		ApplicationId: "936929561302675456",
		SessionId:     "cb06f61453064c0983f2adae2a88c223",
		Data: DSCommand{
			Version: "1077969938624553050",
			Id:      "938956540159881230",
			Name:    "imagine",
			Type:    1,
			Options: []DSOption{{Type: 3, Name: "prompt", Value: prompt}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       "938956540159881230",
				ApplicationId:            "936929561302675456",
				Version:                  "1077969938624553050",
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "imagine",
				Description:              "Lucky you!",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 3, Name: "prompt", Description: "The prompt to imagine", Required: true}},
			},
			Attachments: []ReqCommandAttachments{},
		},
	}
	r, _ := jsoniter.MarshalToString(requestBody)
	err := push(r, url, jobId)
	return err
}

func Upscale(jobId string, index int64, messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       "@@guild_id@@",
		ChannelId:     "@@channel_id@@",
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, messageHash),
		},
	}
	r, _ := jsoniter.MarshalToString(requestBody)
	err := push(r, url, jobId)
	return err
}

func MaxUpscale(jobId, messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       "@@guild_id@@",
		ChannelId:     "@@channel_id@@",
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "1f3dbdf09efdf93d81a3a6420882c92c",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::1::%s::SOLO", messageHash),
		},
	}

	r, _ := jsoniter.MarshalToString(requestBody)
	err := push(r, url, jobId)
	return err
}

func Variate(jobId string, index int64, messageId string, messageHash string) error {
	requestBody := ReqVariationDiscord{
		Type:          3,
		GuildId:       "@@guild_id@@",
		ChannelId:     "@@channel_id@@",
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, messageHash),
		},
	}
	r, _ := jsoniter.MarshalToString(requestBody)
	err := push(r, url, jobId)
	return err
}

func Reset(jobId, messageId string, messageHash string) error {
	requestBody := ReqResetDiscord{
		Type:          3,
		GuildId:       "@@guild_id@@",
		ChannelId:     "@@channel_id@@",
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", messageHash),
		},
	}
	r, _ := jsoniter.MarshalToString(requestBody)
	err := push(r, url, jobId)
	return err
}

func Describe(jobId, uploadName string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       "@@guild_id@@",
		ChannelID:     "@@channel_id@@",
		ApplicationId: "936929561302675456",
		SessionId:     "0033db636f7ce1a951e54cdac7044de3",
		Data: DSCommand{
			Version: "1092492867185950853",
			Id:      "1092492867185950852",
			Name:    "describe",
			Type:    1,
			Options: []DSOption{{Type: 11, Name: "image", Value: 0}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       "1092492867185950852",
				ApplicationId:            "936929561302675456",
				Version:                  "1092492867185950853",
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "describe",
				Description:              "Writes a prompt based on your image.",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 11, Name: "image", Description: "The image to describe", Required: true}},
			},
			Attachments: []ReqCommandAttachments{{
				Id:             "0",
				Filename:       filepath.Base(uploadName),
				UploadFilename: uploadName,
			}},
		},
	}
	r, _ := jsoniter.MarshalToString(requestBody)
	err := push(r, url, jobId)
	return err
}

//func Attachments(name string, size int64) (ResAttachments, error) {
//	requestBody := ReqAttachments{
//		Files: []ReqFile{{
//			Filename: name,
//			FileSize: size,
//			Id:       "1",
//		}},
//	}
//	r,_:=jsoniter.Marshal(requestBody)
//	body, err := push(r, uploadUrl)
//	var data ResAttachments
//	json.Unmarshal(body, &data)
//	return data, err
//}

func push(p, u, jobId string) error {
	data, err := jsoniter.MarshalToString(&QueueBody{
		Params: p,
		Url:    u,
		JobId:  jobId,
	})
	if err != nil {
		return err
	}
	return redis.DB.LPush(context.Background(), _const.DISCORD_SEND_QUEUE, data).Err()
}
