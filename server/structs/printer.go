package structs

type PrintSettings struct {
	TargetNozzleTemp int `json:"target_nozzle_temp,omitempty"`
	TargetHotbedTemp int `json:"target_hotbed_temp,omitempty"`
	FanSpeedPct      int `json:"fan_speed_pct,omitempty"`
	PrintSpeedMode   int `json:"print_speed_mode,omitempty"`
	Zcompensation    int `json:"z_comp,omitempty"`
}
