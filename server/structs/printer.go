package structs

type PrintSettings struct {
	TargetNozzleTemp int `json:"target_nozzle_temp,omitempty"`
	TargetHotbedTemp int `json:"target_hotbed_temp,omitempty"`
	FanSpeedPct      int `json:"fan_speed_pct,omitempty"`
	PrintSpeedMode   int `json:"print_speed_mode,omitempty"`
	Zcompensation    int `json:"z_comp,omitempty"`
}

type PrintState struct {
	// Print state
	Taskid        string `json:"taskid,omitempty"`
	Filename      string `json:"filename,omitempty"`
	PrintTime     int    `json:"print_time,omitempty"`
	SuppliesUsage int    `json:"supplies_usage,omitempty"`
	TotalLayers   int    `json:"total_layers,omitempty"`
	CurrLayer     int    `json:"curr_layer,omitempty"`
	RemainTime    int    `json:"remain_time,omitempty"`
}
