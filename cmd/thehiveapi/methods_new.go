package thehiveapi

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// Start_New инициализирует новый модуль взаимодействия с API TheHive
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *TheHiveApi_New) Start_New(ctx context.Context) (chan<- interfaces.ChannelRequester, error) {
	//обработка маршрутов
	go api.router_new(ctx)

	return api.receivingChannel, nil
}
