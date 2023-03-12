package AI

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"yourbackend/internal/svc"
	"yourbackend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AILogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAILogic(ctx context.Context, svcCtx *svc.ServiceContext) *AILogic {
	return &AILogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type AI struct {
	Text string `json:"text"`
}

func (l *AILogic) AI(req *types.AIreq) (resp *types.AIresp, err error) {
	re, e := l.ToTartget(req)
	if e != nil {
		fmt.Println("错误", e)
		return &types.AIresp{
			Status: 1,
		}, nil
	}
	p, _ := ioutil.ReadAll(re.Body)
	/*
		ans := l.TransformBody(string(p))
		Text := ans["choices"].([]map[string]string)
	*/
	a := AI{}
	json.Unmarshal(p, &a)

	return &types.AIresp{
		Status:       0,
		Chat:         a.Text,
		FinishReason: "none",
	}, nil
}
func (l *AILogic) ToTartget(req *types.AIreq) (*http.Response, error) {
	bo := map[string]interface{}{}
	bo["model"] = "text-davinci-003"
	bo["temperature"] = req.Temerature
	bo["max_tokens"] = req.MaxLength
	bo["top_p"] = req.TopP
	bo["frequency_penalty"] = req.FrequencyPenalty
	bo["presence_penalty"] = req.PresencePenalty
	bo["stop"] = []string{" Human:", " AI:"}
	bo["prompt"] = req.Prompt[1:]
	body, _ := json.Marshal(bo)
	request, _ := http.NewRequest("POST", l.svcCtx.Config.ChatGPT.TargetServer, bytes.NewBuffer(body))
	request.Header.Set("Authorization", "Bearer "+l.svcCtx.Config.ChatGPT.Key)
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{
	}
	resp, e := client.Do(request)
	return resp, e
}
func (l *AILogic) TransformBody(body string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	return result
}
