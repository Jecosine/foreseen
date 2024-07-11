package service

import (
	"github.com/alioth-center/foreseen/app/entity"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/thirdparty/lark"
	"strings"
)

func NewNotifyService() *NotifyService {
	return &NotifyService{client: lark.NewClient(lark.Config{AppID: entity.GlobalConfig.AppID, AppSecret: entity.GlobalConfig.AppSecret})}
}

type NotifyService struct {
	client lark.Client
}

func (srv *NotifyService) NotifyLark(ctx http.Context[*entity.LarkNotifyRequest, *entity.BaseResponse]) {
	logging.Info(logger.NewFields(ctx).WithBaseFields(baseField).WithMessage("processing lark notify request").WithData(ctx.Request()))

	var theme = lark.LarkMarkdownMessageThemeGrey
	switch ctx.Request().Level {
	case "info":
		theme = lark.LarkMarkdownMessageThemeBlue
	case "warn":
		theme = lark.LarkMarkdownMessageThemeOrange
	case "error":
		theme = lark.LarkMarkdownMessageThemeRed
	case "success":
		theme = lark.LarkMarkdownMessageThemeGreen
	}

	var errs []error
	title, content := ctx.Request().Title, ctx.Request().Content
	for _, receiver := range ctx.Request().UserReceivers {
		larkReceiver := lark.Receiver{Type: lark.LarkReceiverIdTypeUserID, Receiver: receiver}
		err := srv.client.SendMarkdownMessage(ctx, larkReceiver, title, content, theme)
		if err != nil {
			errs = append(errs, err)
		}
	}

	for _, receiver := range ctx.Request().GroupReceivers {
		larkReceiver := lark.Receiver{Type: lark.LarkReceiverIdTypeChatID, Receiver: receiver}
		err := srv.client.SendMarkdownMessage(ctx, larkReceiver, title, content, theme)
		if err != nil {
			errs = append(errs, err)
		}
	}

	response := &entity.LarkNotifyResponse{Errors: make([]string, 0)}
	switch len(errs) {
	case 0:
		response.Status = "success"
	case len(ctx.Request().UserReceivers) + len(ctx.Request().GroupReceivers):
		response.Status = "all failure"
		for _, err := range errs {
			response.Errors = append(response.Errors, err.Error())
		}
	default:
		response.Status = "partial failure"
		for _, err := range errs {
			response.Errors = append(response.Errors, err.Error())
		}
	}

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetResponse(&entity.BaseResponse{Data: response})
	logging.Info(logger.NewFields(ctx).WithBaseFields(baseField).WithMessage("lark notify request processed").WithData(response))
}

func NotifySrvCheckToken[req any](ctx http.Context[req, *entity.BaseResponse]) {
	token := strings.TrimPrefix(ctx.HeaderParams().GetString("Authorization"), "Bearer ")
	if token != entity.GlobalConfig.Token {
		ctx.Abort()
		ctx.SetStatusCode(http.StatusUnauthorized)
		ctx.SetResponse(&entity.BaseResponse{ErrorCode: 4011, ErrorMessage: "invalid authentication token"})
		return
	}
}