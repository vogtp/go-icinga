package director

type serviceTemplate map[string]service

type service struct {
	ActionURL             interface{}    `json:"action_url"`
	ApplyFor              interface{}    `json:"apply_for"`
	AssignFilter          interface{}    `json:"assign_filter"`
	CheckCommand          string         `json:"check_command"`
	CheckInterval         int            `json:"check_interval"`
	CheckPeriod           interface{}    `json:"check_period"`
	CheckTimeout          interface{}    `json:"check_timeout"`
	CommandEndpoint       interface{}    `json:"command_endpoint"`
	Disabled              bool           `json:"disabled"`
	DisplayName           interface{}    `json:"display_name"`
	EnableActiveChecks    interface{}    `json:"enable_active_checks"`
	EnableEventHandler    interface{}    `json:"enable_event_handler"`
	EnableFlapping        interface{}    `json:"enable_flapping"`
	EnableNotifications   bool           `json:"enable_notifications"`
	EnablePassiveChecks   interface{}    `json:"enable_passive_checks"`
	EnablePerfdata        bool           `json:"enable_perfdata"`
	EventCommand          interface{}    `json:"event_command"`
	Fields                []interface{}  `json:"fields"`
	FlappingThresholdHigh interface{}    `json:"flapping_threshold_high"`
	FlappingThresholdLow  interface{}    `json:"flapping_threshold_low"`
	Groups                []interface{}  `json:"groups"`
	Host                  interface{}    `json:"host"`
	IconImage             string         `json:"icon_image"`
	IconImageAlt          interface{}    `json:"icon_image_alt"`
	Imports               []string       `json:"imports"`
	MaxCheckAttempts      int            `json:"max_check_attempts"`
	Notes                 string         `json:"notes"`
	NotesURL              string         `json:"notes_url"`
	ObjectName            string         `json:"object_name"`
	ObjectType            string         `json:"object_type"`
	RetryInterval         int            `json:"retry_interval"`
	ServiceSet            interface{}    `json:"service_set"`
	TemplateChoice        interface{}    `json:"template_choice"`
	UseAgent              bool           `json:"use_agent"`
	UseVarOverrides       interface{}    `json:"use_var_overrides"`
	Vars                  map[string]any `json:"vars"`
	Volatile              interface{}    `json:"volatile"`
	Zone                  interface{}    `json:"zone"`
	UUID                  string         `json:"uuid,omitempty"`
}
