package thehiveapi

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// Start инициализирует новый модуль взаимодействия с API TheHive
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *apiTheHiveModule) Start(ctx context.Context) (chan<- commoninterfaces.ChannelRequester, error) {
	//обработка кэша
	api.cache.StartAutomaticExecution(ctx)

	//инициализация автоматической очистки хранилища
	api.storageCache.Start(ctx)

	//обработка маршрутов
	go api.router(ctx)

	return api.receivingChannel, nil
}
