package director

type Config struct {
	Command         Command           `json:"Command"`
	ServiceTemplate ServiceTemplate   `json:"ServiceTemplate"`
	Datafield       map[int]Datafield `json:"Datafield"`
}
