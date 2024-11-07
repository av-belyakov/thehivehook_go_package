package webhookserver

// CreateEvenAlert генератор кейса содержащего в себе дополнительную информацию, такую как
// перечень значений observables и ttp. ЭТОТ МЕТОД ЕЩЕ НЕ ДОДЕЛАН
func CreateEvenAlert(rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventAlert, error) {
	var (
		rmea ReadyMadeEventAlert = ReadyMadeEventAlert{}
	)

	//попробовать запросить весь alert
	//http://192.168.9.38:9000/api/v1/alert/~74465718400

	return rmea, nil
}
