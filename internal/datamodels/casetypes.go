package datamodels

// CaseEventElement типовой элемент описывающий события, 'case', приходящие из TheHive
type CaseEventElement struct {
	Operation  string              `json:"operation"`  //тип операции
	ObjectType string              `json:"objectType"` //тип объекта
	RootId     string              `json:"rootId"`     //основной идентификатор объекта
	Details    DetailsEventElement `json:"details"`    //частичные детали по объекту
	Object     ObjectEventElement  `json:"object"`     //частичная информация по объекту
}

// ObjectEventElement содержит информацию из поля 'object' приходящего из TheHive элемента
type ObjectEventElement struct {
	Tags      []string `json:"tags"`
	CreatedAt int64    `json:"createdAt"`
	CaseId    int      `json:"caseId"`
}

// DetailsEventElement содержит информацию из поля 'details'
type DetailsEventElement struct {
	Status string `json:"status"`
}

// BaseCaseEventElement приходит от TheHive на запрос по номеру кейса
type BaseCaseEventElement struct {
	CustomFields     []any    `json:"customFields"`
	ExtraData        any      `json:"extraData"`
	Tags             []string `json:"tags"`
	ID               string   `json:"_id"`
	Type             string   `json:"_type"`
	CreatedBy        string   `json:"_createdBy"`
	UpdatedBy        string   `json:"_updatedBy"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Status           string   `json:"status"`
	Summary          string   `json:"summary"`
	ImpactStatus     string   `json:"impactStatus"`
	ResolutionStatus string   `json:"resolutionStatus"`
	Assignee         string   `json:"assignee"`
	CreatedAt        int64    `json:"_createdAt"`
	UpdatedAt        int64    `json:"_updatedAt"`
	StartDate        int64    `json:"startDate"`
	EndDate          int64    `json:"endDate"`
	Number           int      `json:"number"`
	Severity         int      `json:"severity"`
	Tlp              int      `json:"tlp"`
	Pap              int      `json:"pap"`
	Flag             bool     `json:"flag"`
}
