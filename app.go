package main

import (
	"context"
	"encoding/base64"
)

var (
	VMTB = VoiceMultiTalkBody{}
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	AIInit()
	VMTB.InitVoiceMultiTalkBody(ERNIE_Speed_128K)
}

// Greet returns a greeting for the given name
func (a *App) V2T(audioData []byte) (string, error) {
	return Voice2Text(audioData)
}

func (a *App) T2V(text string) string {
	mp3data, _ := Text2Voice(text)
	base64Str := base64.StdEncoding.EncodeToString(mp3data)
	return base64Str
}
func (a *App) VMT(text string) string {
	return VMTB.VoiceMultiTalk(text)
}
