package models

type HealthcheckStatus struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type HealthcheckResponse struct {
	Database HealthcheckStatus `json:"database"`
	Storage  HealthcheckStatus `json:"storage"`
	Redis    HealthcheckStatus `json:"redis"`
	Discord  HealthcheckStatus `json:"discord"`
	AllOk    bool              `json:"all_ok"`
}

func HealthcheckStatusFromError(err error) HealthcheckStatus {
	s := HealthcheckStatus{}

	s.Ok = err == nil
	if err != nil {
		s.Message = err.Error()
	}

	return s
}
