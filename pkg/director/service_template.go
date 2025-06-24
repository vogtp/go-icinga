package director

type serviceTemplate map[string]service

type service struct {
	ActionURL             interface{}    `json:"action_url,omitempty"`
	ApplyFor              interface{}    `json:"apply_for,omitempty"`
	AssignFilter          interface{}    `json:"assign_filter,omitempty"`
	CheckCommand          string         `json:"check_command,omitempty"`
	CheckInterval         int            `json:"check_interval,omitempty"`
	CheckPeriod           interface{}    `json:"check_period,omitempty"`
	CheckTimeout          interface{}    `json:"check_timeout,omitempty"`
	CommandEndpoint       interface{}    `json:"command_endpoint,omitempty"`
	Disabled              bool           `json:"disabled,omitempty"`
	DisplayName           interface{}    `json:"display_name,omitempty"`
	EnableActiveChecks    interface{}    `json:"enable_active_checks,omitempty"`
	EnableEventHandler    interface{}    `json:"enable_event_handler,omitempty"`
	EnableFlapping        interface{}    `json:"enable_flapping,omitempty"`
	EnableNotifications   bool           `json:"enable_notifications,omitempty"`
	EnablePassiveChecks   interface{}    `json:"enable_passive_checks,omitempty"`
	EnablePerfdata        bool           `json:"enable_perfdata,omitempty"`
	EventCommand          interface{}    `json:"event_command,omitempty"`
	Fields                []interface{}  `json:"fields,omitempty"`
	FlappingThresholdHigh interface{}    `json:"flapping_threshold_high,omitempty"`
	FlappingThresholdLow  interface{}    `json:"flapping_threshold_low,omitempty"`
	Groups                []interface{}  `json:"groups,omitempty"`
	Host                  interface{}    `json:"host,omitempty"`
	IconImage             string         `json:"icon_image,omitempty"`
	IconImageAlt          interface{}    `json:"icon_image_alt,omitempty"`
	Imports               []string       `json:"imports,omitempty"`
	MaxCheckAttempts      int            `json:"max_check_attempts,omitempty"`
	Notes                 string         `json:"notes,omitempty"`
	NotesURL              string         `json:"notes_url,omitempty"`
	ObjectName            string         `json:"object_name,omitempty"`
	ObjectType            string         `json:"object_type,omitempty"`
	RetryInterval         int            `json:"retry_interval,omitempty"`
	ServiceSet            interface{}    `json:"service_set,omitempty"`
	TemplateChoice        interface{}    `json:"template_choice,omitempty"`
	UseAgent              bool           `json:"use_agent,omitempty"`
	UseVarOverrides       interface{}    `json:"use_var_overrides,omitempty"`
	Vars                  map[string]any `json:"vars,omitempty"`
	Volatile              interface{}    `json:"volatile,omitempty"`
	Zone                  interface{}    `json:"zone,omitempty"`
	UUID                  string         `json:"uuid,omitempty"`
}
