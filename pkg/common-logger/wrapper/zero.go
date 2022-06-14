package wrapper

import (
	"github.com/rs/zerolog"
)

type Zerolog struct {
	log *zerolog.Logger
}

func (s Zerolog) Debug(msg string, fields map[string]interface{}) {
	s.log.Debug().Fields(fields).Msg(msg)
}

func (s Zerolog) Info(msg string, fields map[string]interface{}) {
	s.log.Info().Fields(fields).Msg(msg)
}

func (s Zerolog) Warn(msg string, fields map[string]interface{}) {
	s.log.Warn().Fields(fields).Msg(msg)
}

func (s Zerolog) Error(msg string, err error, fields map[string]interface{}) {
	if err != nil {
		if fields == nil {
			fields = make(map[string]interface{})
		}

		fields["error"] = err
	}

	s.log.Error().Fields(fields).Msg(msg)
}

func (s Zerolog) Fatal(msg string, err error, fields map[string]interface{}) {
	if err != nil {
		if fields == nil {
			fields = make(map[string]interface{})
		}

		fields["error"] = err
	}

	s.log.Fatal().Fields(fields).Msg(msg)
}

func WrapZerolog(l *zerolog.Logger) Zerolog {
	return Zerolog{log: l}
}
