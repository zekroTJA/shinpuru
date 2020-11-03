package jdoodle

type credentialsBody struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type execRequestBody struct {
	*credentialsBody

	Script   string `json:"script"`
	Language string `json:"language"`
}

type responseError struct {
	Error string `json:"error"`
}

type ExecResponse struct {
	Output  string `json:"output"`
	Memory  string `json:"memory"`
	CPUTime string `json:"cpuTime"`
}

type CreditsResponse struct {
	Used int `json:"used"`
}
