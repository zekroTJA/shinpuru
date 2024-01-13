// Package lokiwriter implements rogu.Writer to push
// logs to a Grafana Loki instance.
package lokiwriter

import (
	"fmt"

	"github.com/zekrotja/promtail"
	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/level"
)

// Writer implements rogu.Writer for pusing logs
// to a loki instance.
type Writer struct {
	client promtail.Client
}

var (
	_ rogu.Writer = (*Writer)(nil)
	_ rogu.Closer = (*Writer)(nil)
)

// NewWriter creates a new loki Writer with the
// given connection options.
func NewWriter(options Options) (t *Writer, err error) {
	t = &Writer{}

	t.client, err = promtail.NewJSONv1Client(options.Address, options.Labels,
		promtail.WithBasicAuth(options.Username, options.Password),
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t Writer) Write(
	lvl level.Level,
	fields []*rogu.Field,
	tag string,
	err error,
	errFormat string,
	callerFile string,
	callerLine int,
	msg string,
) error {
	labels := map[string]string{}

	for _, field := range fields {
		labels[stringify(field.Key)] = stringify(field.Val)
	}

	if tag != "" {
		labels["tag"] = tag
	}

	if err != nil {
		labels["error"] = err.Error()
	}

	if callerFile != "" {
		labels["caller"] = fmt.Sprintf("%s:%d", callerFile, callerLine)
	}

	t.client.LogfWithLabels(translateLevel(lvl), labels, msg)
	return nil
}

func (t Writer) Close() error {
	t.client.Close()
	return nil
}

func translateLevel(lvl level.Level) promtail.Level {
	switch lvl {
	case level.Trace, level.Debug:
		return promtail.Debug
	case level.Info:
		return promtail.Info
	case level.Warn:
		return promtail.Warn
	case level.Error:
		return promtail.Error
	case level.Fatal:
		return promtail.Fatal
	case level.Panic:
		return promtail.Panic
	default:
		return promtail.Info
	}
}

func stringify(v any) string {
	return fmt.Sprintf("%v", v)
}
