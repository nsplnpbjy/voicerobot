package main

import (
	"context"
	"errors"
	"log"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ModelName string

const (
	AIAccessKey = "0b69c9516e4e439b8a1c2809525727e9"
	AISecretKey = "2a769031f34b4d7e9f093677e024efe9"
	//收费的
	ERNIE_4_0_8K       ModelName = "ERNIE-4.0-8K"
	ERNIE_4_0_Turbo_8K ModelName = "ERNIE-4.0-Turbo-8K"
	//免费的
	ERNIE_Speed_128K ModelName = "ERNIE-Speed-128K"
)

const (
	ERNIE_4_0_8K_CHOICE       string = "ERNIE-4.0-8K"
	ERNIE_4_0_Turbo_8K_CHOICE string = "ERNIE-4.0-Turbo-8K"
	ERNIE_Speed_128K_CHOICE   string = "ERNIE-Speed-128K"
)

type MultiChatBody struct {
	IsChanged  bool
	IsChatting bool
	Msg        []qianfan.ChatCompletionMessage
	Chat       *qianfan.ChatCompletion
}

var (
	multiChatBody = MultiChatBody{false, false, []qianfan.ChatCompletionMessage{}, nil}
)

func GetIsChanged() bool {
	return multiChatBody.IsChanged
}

func AIInit() {
	// 初始化配置
	qianfan.GetConfig().AccessKey = AIAccessKey
	qianfan.GetConfig().SecretKey = AISecretKey
	qianfan.WithLLMRetryTimeout(10000)
}

func GetQianfanChat(modelName ModelName) *qianfan.ChatCompletion {
	// 指定特定模型
	return qianfan.NewChatCompletion(
		qianfan.WithModel(string(modelName)),
	)
}

func DoChat(text string, chat *qianfan.ChatCompletion) string {
	resp, err := chat.Do(
		context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages: []qianfan.ChatCompletionMessage{
				qianfan.ChatCompletionUserMessage(text),
			},
		},
	)
	if err != nil {
		log.Printf("Error during chat completion: %v", err)
		return ""
	}
	return resp.Result
}

func MultiChatNewText(text string) error {
	if text == "" {
		return errors.New("输入不能为空")
	}
	multiChatBody.IsChanged = true
	multiChatBody.IsChatting = true
	multiChatBody.Msg = append(multiChatBody.Msg, qianfan.ChatCompletionUserMessage(text))
	return nil
}

func MultiChatClose() {
	multiChatBody = MultiChatBody{false, false, nil, nil}
}

func StartMultiChat(modelName ModelName, ctx context.Context) {
	multiChatBody.IsChatting = true
	multiChatBody.Chat = GetQianfanChat(modelName)
	for {
		if !multiChatBody.IsChatting {
			break
		}
		if multiChatBody.IsChanged {
			DoMultiChat(modelName, ctx)
		}
	}
}

func DoMultiChat(modelName ModelName, ctx context.Context) {
	resp, chatError := multiChatBody.Chat.Stream(context.TODO(), &qianfan.ChatCompletionRequest{Messages: multiChatBody.Msg})
	if chatError != nil {
		println(chatError.Error())
		return
	}
	defer resp.Close()
	recv := ""
	runtime.EventsEmit(ctx, "output", "\n"+modelName+":")
	for {
		r, err := resp.Recv()
		if err != nil {
			println(err.Error())
			break
		}
		if r.IsEnd {
			break
		} else {
			recv = recv + r.Result
			runtime.EventsEmit(ctx, "output", r.Result)
		}
	}
	runtime.EventsEmit(ctx, "output", "\n")
	multiChatBody.IsChanged = false
	multiChatBody.Msg = append(multiChatBody.Msg, qianfan.ChatCompletionAssistantMessage(recv))
}

type VoiceMultiTalkBody struct {
	Msg  []qianfan.ChatCompletionMessage
	Chat *qianfan.ChatCompletion
}

func (v *VoiceMultiTalkBody) InitVoiceMultiTalkBody(modelname ModelName) {
	v.Msg = []qianfan.ChatCompletionMessage{}
	v.Chat = GetQianfanChat(modelname)
}

func (v *VoiceMultiTalkBody) VoiceMultiTalk(text string) string {
	v.Msg = append(v.Msg, qianfan.ChatCompletionUserMessage(text))
	resp, chatError := v.Chat.Do(context.TODO(), &qianfan.ChatCompletionRequest{Messages: v.Msg})
	if chatError != nil {
		println(chatError.Error())
		return ""
	}
	recv := resp.Result
	v.Msg = append(v.Msg, qianfan.ChatCompletionAssistantMessage(recv))
	return recv
}
