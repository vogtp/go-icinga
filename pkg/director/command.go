package director

type command map[string]CommandDefinition

type CommandDefinition struct {
	Arguments      map[string]argument `json:"arguments,omitempty"`
	Command        string              `json:"command,omitempty"`
	Disabled       bool                `json:"disabled,omitempty"`
	Fields         []cmdField          `json:"fields,omitempty"`
	Imports        []interface{}       `json:"imports,omitempty"`
	IsString       interface{}         `json:"is_string,omitempty"`
	MethodsExecute string              `json:"methods_execute,omitempty"`
	ObjectName     string              `json:"object_name,omitempty"`
	ObjectType     string              `json:"object_type,omitempty"`
	Timeout        int                 `json:"timeout,omitempty"`
	Vars           struct {
	} `json:"vars"`
	Zone interface{} `json:"zone,omitempty"`
	UUID string      `json:"uuid,omitempty"`
}

type argument struct {
	Value string `json:"value,omitempty"`
	SetIf    string `json:"set_if,omitempty"` // DO NOT USE: it breaks the parameter parsing
	Required bool `json:"required,omitempty"`
	SkipKey  bool `json:"skip_key,omitempty"`
	Order    int  `json:"order,omitempty"`
}

type cmdField struct {
	DatafieldID int         `json:"datafield_id,omitempty"`
	IsRequired  string      `json:"is_required,omitempty"`
	VarFilter   interface{} `json:"var_filter"`
}
