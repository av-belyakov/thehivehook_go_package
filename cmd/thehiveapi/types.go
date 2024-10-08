package thehiveapi

type apiTheHive struct {
	port   int
	host   string
	apiKey string
}

type RootQuery struct {
	Query []Query `json:"query"`
}

type Query struct {
	Name      string   `json:"_name,omitempty"`
	IDOrName  string   `json:"idOrName,omitempty"`
	From      int64    `json:"from"`
	To        int      `json:"to,omitempty"`
	ExtraData []string `json:"extraData,omitempty"`
}

type ErrorAnswer struct {
	Err     string `json:"type"`
	Message string `json:"message"`
}
