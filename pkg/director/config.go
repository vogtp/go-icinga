package director

type Config struct {
	Command         command           `json:"Command,omitempty"`
	ServiceTemplate serviceTemplate   `json:"ServiceTemplate,omitempty"`
	Datafield       map[int]datafield `json:"Datafield,omitempty"`
}
