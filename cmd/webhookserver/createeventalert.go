package webhookserver

// CreateEvenAlert генератор кейса содержащего в себе дополнительную информацию, такую как
// перечень значений observables и ttp. ЭТОТ МЕТОД ЕЩЕ НЕ ДОДЕЛАН
func CreateEvenAlert(rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventAlert, error) {
	var (
		rmea ReadyMadeEventAlert = ReadyMadeEventAlert{}
	)

	//попробовать запросить весь alert
	//http://192.168.9.38:9000/api/v1/alert/~74465718400

	/*
	   			query
	   :
	   [{_name: "getAlert", idOrName: "~76432666760"},…]
	   0
	   :
	   {_name: "getAlert", idOrName: "~76432666760"}
	   1
	   :
	   {_name: "similarCases", caseFilter: {_and: [{_field: "status", _value: "Open"},…]}}
	   caseFilter
	   :
	   {_and: [{_field: "status", _value: "Open"},…]}
	*/

	return rmea, nil
}
