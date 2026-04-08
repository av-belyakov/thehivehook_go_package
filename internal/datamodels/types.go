package datamodels

// VerifiedObjectEventCase готовый, сформированный объект содержащий информацию по кейсу
type VerifiedObjectEventCase struct {
	Case        map[string]any `json:"event"`
	Observables []any          `json:"observables"`
	TTPs        []any          `json:"ttp"`
	Source      string         `json:"source"`
}

// VerifiedObjectEventAlert готовый, сформированный объект содержащий информацию по алерту
type VerifiedObjectEventAlert struct {
	Event  map[string]any `json:"event"`
	Alert  map[string]any `json:"alert"`
	Source string         `json:"source"`
}

type VerifiedResponseAcceptedCommand struct {
	Data       any    `json:"data"`
	Id         string `json:"id"`
	Source     string `json:"source"`
	Command    string `json:"command"`
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}
