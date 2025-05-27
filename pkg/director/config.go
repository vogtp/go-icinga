package director

type Config struct {
	Command         command           `json:"Command"`
	ServiceTemplate serviceTemplate   `json:"ServiceTemplate"`
	Datafield       map[int]datafield `json:"Datafield"`
}
