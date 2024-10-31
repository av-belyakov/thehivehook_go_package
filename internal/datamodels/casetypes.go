package datamodels

// CaseEventElement типовой элемент описывающий события, 'case', приходящие из TheHive
type CaseEventElement struct {
	Operation  string              `json:"operation"`  //тип операции
	ObjectType string              `json:"objectType"` //тип объекта
	RootId     string              `json:"rootId"`     //основной идентификатор объекта
	Object     ObjectEventElement  `json:"object"`     //частичная информация по объекту
	Details    DetailsEventElement `json:"details"`    //частичные детали по объекту
}

// ObjectEventElement содержит информацию из поля 'object' приходящего из TheHive элемента
type ObjectEventElement struct {
	CaseId    string   `json:"caseId"`
	CreatedAt int64    `json:"createdAt"`
	Tags      []string `json:"tags"`
}

// DetailsEventElement содержит информацию из поля 'details'
type DetailsEventElement struct {
	Status string `json:"status"`
}

// BaseCaseEventElement приходит от TheHive на запрос по номеру кейса
type BaseCaseEventElement struct {
	ID               string        `json:"_id"`
	Type             string        `json:"_type"`
	CreatedBy        string        `json:"_createdBy"`
	UpdatedBy        string        `json:"_updatedBy"`
	CreatedAt        int64         `json:"_createdAt"`
	UpdatedAt        int64         `json:"_updatedAt"`
	Number           int           `json:"number"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	Severity         int           `json:"severity"`
	StartDate        int64         `json:"startDate"`
	EndDate          int64         `json:"endDate"`
	Tags             []string      `json:"tags"`
	Flag             bool          `json:"flag"`
	Tlp              int           `json:"tlp"`
	Pap              int           `json:"pap"`
	Status           string        `json:"status"`
	Summary          string        `json:"summary"`
	ImpactStatus     string        `json:"impactStatus"`
	ResolutionStatus string        `json:"resolutionStatus"`
	Assignee         string        `json:"assignee"`
	CustomFields     []interface{} `json:"customFields"`
	ExtraData        interface{}   `json:"extraData"`
}
