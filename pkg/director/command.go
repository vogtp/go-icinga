package director

type command map[string]CommandDefinition

type CommandDefinition struct {
	Arguments      map[string]argument `json:"arguments"`
	Command        string              `json:"command"`
	Disabled       bool                `json:"disabled"`
	Fields         []cmdField          `json:"fields"`
	Imports        []interface{}       `json:"imports"`
	IsString       interface{}         `json:"is_string"`
	MethodsExecute string              `json:"methods_execute"`
	ObjectName     string              `json:"object_name"`
	ObjectType     string              `json:"object_type"`
	Timeout        int                 `json:"timeout"`
	Vars           struct {
	} `json:"vars"`
	Zone interface{} `json:"zone"`
	UUID string      `json:"uuid"`
}

type argument struct {
	Value    string `json:"value"`
	SetIf    string `json:"set_if,omitempty"`
	Required bool   `json:"required"`
	SkipKey  bool   `json:"skip_key"`
}

type cmdField struct {
	DatafieldID int         `json:"datafield_id"`
	IsRequired  string      `json:"is_required"`
	VarFilter   interface{} `json:"var_filter"`
}
