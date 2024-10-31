package webhookserver

// CreateEvenAlert генератор кейса содержащего в себе дополнительную информацию, такую как
// перечень значений observables и ttp. ЭТОТ МЕТОД ЕЩЕ НЕ ДОДЕЛАН
func CreateEvenAlert(rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventAlert, error) {
	var (
		rmea ReadyMadeEventAlert = ReadyMadeEventAlert{}
	)

	return rmea, nil
}
