package sender

import (
	_const "aiStudio/internal/const"
	"aiStudio/internal/midj/listener"
	"aiStudio/internal/redis"
	"aiStudio/internal/repository"
	"aiStudio/internal/repository/model"
	"bytes"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handle(ctx context.Context) {
	go func() {
		for {
			// 发送队列
			data, _ := redis.DB.BRPop(ctx, time.Minute, _const.DISCORD_SEND_QUEUE).Result()
			//if err != nil {
			//	logs.Warn("sender handle error: %s", err)
			//}
			if len(data) != 2 {
				time.Sleep(30 * time.Second)
				continue
			}

			var b = QueueBody{}
			err := jsoniter.Unmarshal([]byte(data[1]), &b)
			if err != nil {
				continue
			}
			// 频道锁，一个频道只能跑一个任务
			gid, cid, token := sendJob(ctx, &b)
			b.Params = strings.NewReplacer("@@guild_id@@", gid, "@@channel_id@@", cid).Replace(b.Params)
			b.Url = strings.NewReplacer("@@guild_id@@", gid, "@@channel_id@@", cid).Replace(b.Url)
			request(&b, token)
		}
	}()
	for {
		// 结果队列
		data, _ := redis.DB.BRPop(ctx, time.Minute, _const.DISCORD_RESULTS_QUEUE).Result()
		//if err != nil {
		//	logs.Warn("listen handle error: %s", err)
		//}
		if len(data) != 2 {
			time.Sleep(10 * time.Second)
			continue
		}
		var b listener.ReqCb
		err := jsoniter.Unmarshal([]byte(data[1]), &b)
		if err != nil {
			continue
		}
		result(ctx, &b)
	}
}
func result(ctx context.Context, b *listener.ReqCb) error {
	var (
		channleID string
		guildID   string
		progress  int
		status    string
		msgId     string
		option    *string
		imageUrl  *string
		mhash     *string
		isRelease bool
		msg       = b.DiscordMsg
	)
	switch b.Type {
	case listener.FirstTrigger:
	case listener.UpdateMsg:
		var userId string
		for _, c := range confs {
			if strings.Contains(msg.Content, fmt.Sprintf("<@%s>", c.UserID)) {
				userId = c.UserID
			}
		}
		sp := strings.Split(msg.Content, fmt.Sprintf("<@%s>", userId))
		if len(sp) > 1 {
			//option = &sp[len(sp)-1]
			s := strings.Split(*option, "%>")[0]
			s = strings.Split(s, "%")[0]
			s = strings.ReplaceAll(s, "(", "")
			s = strings.TrimSpace(s)
			progress, _ = strconv.Atoi(s)
			status = "gen"
		}
	case listener.GenerateEditError:
		option = &msg.Content
		status = "error"
		isRelease = true
	case listener.GenerateEnd:
		progress = 100
		status = "done"
		isRelease = true
	default:
		return nil
	}

	msgId, mhash, imageUrl = ParseMsgIdAndHash(msg)
	channleID = msg.ChannelID
	guildID = msg.GuildID

	genId, err := redis.DB.Get(ctx, _const.DISCORD_CHANNEL_QUEUE+guildID+"-"+channleID).Result()
	if err != nil {
		return err
	}

	err = repository.UpdateRecordByGenId(genId, model.RecordTable{
		Progress:  progress,
		Status:    status,
		GuildID:   guildID,
		ChannelID: channleID,
		MsgID:     msgId,
		Option:    option,
		ImageUrl:  imageUrl,
		MHash:     mhash,
	})
	if err != nil {
		return err
	}

	if isRelease {
		// 释放频道锁
		if err := redis.DB.Del(ctx, _const.DISCORD_CHANNEL_QUEUE+guildID+"-"+channleID).Err(); err != nil {
			return err
		}
	}
	return nil
}

func request(b *QueueBody, userToken string) ([]byte, error) {
	req, err := http.NewRequest("POST", b.Url, bytes.NewBuffer([]byte(b.Params)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", userToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.35")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bod, respErr := io.ReadAll(response.Body)
	fmt.Println("response: ", string(bod), respErr, response.Status)
	return bod, respErr
}

func sendJob(ctx context.Context, b *QueueBody) (guildId, channid, userToken string) {
	for {
		for _, c := range confs {
			guildId = c.GuildID
			userToken = c.UserToken
			for _, channid = range c.ChannelID {
				if redis.DB.SetNX(ctx, _const.DISCORD_CHANNEL_QUEUE+guildId+"-"+channid, b.JobId, 10*time.Minute).Val() {
					return
				}
			}
			time.Sleep(5 * time.Second)
		}
	}
}
