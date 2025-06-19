package director

type datafield struct {
	Varname     string      `json:"varname"`
	Caption     string      `json:"caption"`
	Description string      `json:"description"`
	Datatype    string      `json:"datatype"`
	Format      interface{} `json:"format"`
	Settings    struct {
	} `json:"settings"`
	UUID string `json:"uuid,omitempty"`
}
