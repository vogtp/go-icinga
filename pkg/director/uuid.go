package director

import "github.com/google/uuid"

var GenerateUUID = func() string {
	return uuid.NewString()
	//return ""
}
