package webhookserver

// CreateEvenAlert генератор кейса содержащего в себе дополнительную информацию, такую как
// перечень значений observables и ttp
func CreateEvenAlert(uuidStorage, rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventAlert, error) {
	var (
		rmea ReadyMadeEventAlert = ReadyMadeEventAlert{}
	)

	return rmea, nil
}
