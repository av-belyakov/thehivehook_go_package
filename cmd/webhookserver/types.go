package webhookserver

import (
	"context"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	temporarystorage "github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver/temporarystorage"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

// WebHookServer непосредственно сам сервер
type WebHookServer struct {
	ttl       int
	port      int
	host      string
	name      string //gcm, rcmmsk и т.д.
	version   string
	ctx       context.Context
	server    *http.Server
	storage   *temporarystorage.WebHookTemporaryStorage
	logger    *logginghandler.LoggingChan
	chanInput chan<- ChanFormWebHookServer
}

// webHookServerOptions функциональные параметры
type webHookServerOptions func(*WebHookServer)

// WebHookServerOptions основные опции
type WebHookServerOptions struct {
	TTL     int
	Port    int
	Host    string
	Name    string
	Version string
}

// ChanFormWebHookServer структура канала для взаимодействия сторонних сервисов с webhookserver
type ChanFormWebHookServer struct {
	ForSomebody string
	Data        commoninterfaces.ChannelRequester
}

type EventElement struct {
	Operation  string `json:"operation"`
	ObjectType string `json:"objectType"`
	RootId     string `json:"rootId"`
}

type ReadyMadeEventCase struct {
	Source string                 `json:"source"`
	Case   map[string]interface{} `json:"event"`
	//Case        CaseEvent     `json:"event"`
	Observables []interface{} `json:"observables"`
	TTPs        []interface{} `json:"ttp"`
}

// ************** case ***************
type CaseEvent struct {
	Operation      string  `json:"operation"`
	Details        Details `json:"details"`
	ObjectType     string  `json:"objectType"`
	ObjectID       string  `json:"objectId"`
	Base           bool    `json:"base"`
	StartDate      int64   `json:"startDate"`
	RootID         string  `json:"rootId"`
	RequestID      string  `json:"requestId"`
	Object         Object  `json:"object"`
	OrganisationID string  `json:"organisationId"`
	Organisation   string  `json:"organisation"` //nolint
}

type Details struct {
	EndDate          int64                  `json:"endDate,omitempty"`
	CustomFields     map[string]interface{} `json:"customFields,omitempty"`
	ResolutionStatus string                 `json:"resolutionStatus,omitempty"`
	Summary          string                 `json:"summary,omitempty"`
	Status           string                 `json:"status,omitempty"`
	ImpactStatus     string                 `json:"impactStatus,omitempty"`
}

type (
	Stats  struct{}
	Object struct {
		ID_              string                 `json:"_id"` //nolint
		ID               string                 `json:"id"`
		CreatedBy        string                 `json:"createdBy"`
		UpdatedBy        string                 `json:"updatedBy"`
		CreatedAt        int64                  `json:"createdAt"`
		UpdatedAt        int64                  `json:"updatedAt"`
		Type             string                 `json:"_type"`
		CaseID           int                    `json:"caseId"`
		Title            string                 `json:"title"`
		Description      string                 `json:"description"`
		Severity         int                    `json:"severity"`
		StartDate        int64                  `json:"startDate"`
		EndDate          int64                  `json:"endDate"`
		ImpactStatus     string                 `json:"impactStatus"`
		ResolutionStatus string                 `json:"resolutionStatus"`
		Tags             []string               `json:"tags"`
		Flag             bool                   `json:"flag"`
		Tlp              int                    `json:"tlp"`
		Pap              int                    `json:"pap"`
		Status           string                 `json:"status"`
		Summary          string                 `json:"summary"`
		Owner            string                 `json:"owner"`
		CustomFields     map[string]interface{} `json:"customFields"`
		Stats            Stats                  `json:"stats"`
		Permissions      []interface{}          `json:"permissions"`
	}
)
