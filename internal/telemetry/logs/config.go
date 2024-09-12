package logs

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
	"github.com/zitadel/logging/otel"
	"github.com/zitadel/zitadel/internal/telemetry/logs/record"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/otel/log"
)

type Hook string

const (
	GCPLoggingOtelExporter Hook = "GCPLoggingOtelExporter"
)

type Config struct {
	Log   logging.Config `mapstructure:",squash"`
	Hooks map[string]map[string]interface{}
}

type GCPExporterConfig struct {
	AddedAttributes map[string]string
}

func (g *GCPExporterConfig) ToAttributes() (attributes []log.KeyValue) {
	for k, v := range g.AddedAttributes {
		attributes = append(attributes, log.KeyValue{Key: k, Value: log.StringValue(fmt.Sprintf("%v", v))})
	}
	return attributes
}

func (c *Config) SetLogger() (err error) {
	var hooks []logrus.Hook
	for name, rawCfg := range c.Hooks {
		switch name {
		case strings.ToLower(string(GCPLoggingOtelExporter)):
			var hook *otel.GcpLoggingExporterHook
			addedAttributes := &GCPExporterConfig{}
			if err = decodeRawConfig(rawCfg, addedAttributes); err != nil {
				return err
			}
			hook, err = otel.NewGCPLoggingExporterHook(
				otel.WithExporterConfig(func(cfg *googlecloudexporter.Config) {
					cfg.LogConfig.DefaultLogName = "zitadel"
					cfg.LogConfig.ServiceResourceLabels = false
					err = decodeRawConfig(rawCfg, cfg)
				}),
				otel.WithOtelSettings(func(cfg *exporter.Settings) {
					err = decodeRawConfig(rawCfg, cfg)
				}),
				otel.WithInclude(func(entry *logrus.Entry) bool {
					return entry.Data["stream"] == record.StreamActivity
				}),
				otel.WithLevels([]logrus.Level{logrus.InfoLevel}),
				otel.WithAttributes(addedAttributes.ToAttributes()),
				otel.WithMapBody(func(entry *logrus.Entry) (body string) {
					entryCopy := *entry
					lg := logrus.New()
					lg.Formatter = &logrus.TextFormatter{
						DisableColors:    true,
						DisableQuote:     true,
						DisableTimestamp: true,
						PadLevelText:     true,
						QuoteEmptyFields: true,
					}
					entryCopy.Logger = lg
					body, err = entryCopy.String()
					return body
				}),
			)
			if err != nil {
				return err
			}
			if err = hook.Start(); err != nil {
				return err
			}
			hooks = append(hooks, hook)
		default:
			return fmt.Errorf("unknown hook: %s", name)
		}
	}
	return c.Log.SetLogger(
		logging.AddHooks(hooks...),
	)
}

func decodeRawConfig(rawConfig map[string]interface{}, typedConfig any) (err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		MatchName: func(mapKey, fieldName string) bool {
			return strings.ToLower(mapKey) == strings.ToLower(fieldName)
		},
		WeaklyTypedInput: true,
		Result:           typedConfig,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(rawConfig)
}
