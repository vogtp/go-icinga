package director

type datafield struct {
	Varname     string      `json:"varname,omitempty"`
	Caption     string      `json:"caption,omitempty"`
	Description string      `json:"description,omitempty"`
	Datatype    string      `json:"datatype,omitempty"`
	Format      interface{} `json:"format,omitempty"`
	Settings    struct {
	} `json:"settings,omitempty"`
	UUID string `json:"uuid,omitempty"`
}
