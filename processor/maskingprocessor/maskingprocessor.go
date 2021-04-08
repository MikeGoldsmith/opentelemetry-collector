package maskingprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer/pdata"
	"go.uber.org/zap"
)

type maskingProcessor struct {
}

func newMaskingProcessor(logger *zap.Logger, cfg *Config) (*maskingProcessor, error) {
	processor := &maskingProcessor{}
	return processor, nil
}

func (mp *maskingProcessor) ProcessMetrics(_ context.Context, pdm pdata.Metrics) (pdata.Metrics, error) {

	rms := pdm.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)

		ilms := rm.InstrumentationLibraryMetrics()
		for j := 0; j < ilms.Len(); j++ {
			ms := ilms.At(j).Metrics()
			for k := 0; k < ms.Len(); k++ {
				m := ms.At(k)
				m.SetName(Reverse(m.Name()))
				m.SetDescription(Reverse(m.Description()))
			}
		}
	}

	return pdm, nil
}

// mock mask function - this would call out to some masking service, ideally in a batch
func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}
