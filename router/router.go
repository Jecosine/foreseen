package router

import (
	"github.com/alioth-center/foreseen/app/api"
	"github.com/alioth-center/foreseen/app/entity"
	"github.com/alioth-center/infrastructure/network/http"
)

func init() {
	engine := http.NewEngine("/foreseen")

	engine.AddEndPoints(
		http.NewEndPointBuilder[*entity.LarkNotifyRequest, *entity.BaseResponse]().
			SetRouter(http.NewRouter("/v1/notify/lark")).
			SetAllowMethods(http.POST).
			SetNecessaryHeaders("Content-Type", "Authorization").
			SetHandlerChain(api.NotifyApi.NotifyLark()).
			Build(),
	)

	engine.ServeAsync("0.0.0.0:8881", make(chan struct{}))
}
