package config

type ValidationError interface {
	error
	Key() string
}

type ve struct {
	key, message string
}

func (e ve) Error() string {
	return e.message
}

func (e ve) Key() string {
	return e.key
}

func Validate(p Provider) error {
	cfg := p.Config()

	if cfg.Privacy.NoticeURL == "" {
		return ve{"privacy.noticeurl", "A privacy notice URL must be provided."}
	}

	if len(cfg.Privacy.Contact) == 0 || cfg.Privacy.Contact[0].Title == "" || cfg.Privacy.Contact[0].Value == "" {
		return ve{"privacy.contact", "At least one valid privacy contact must be provided."}
	}

	return nil
}
