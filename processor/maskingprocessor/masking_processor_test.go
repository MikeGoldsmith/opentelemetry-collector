package maskingprocessor

import (
	"context"
	"testing"
	"time"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/translator/internaldata"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMask(t *testing.T) {
	result := Reverse("apple")
	assert.Equal(t, result, "elppa")
}

func TestMaskingProcessor(t *testing.T) {

	next := new(consumertest.MetricsSink)
	cfg := &Config{
		ProcessorSettings: config.NewProcessorSettings(typeStr),
	}
	factory := NewFactory()
	mp, err := factory.CreateMetricsProcessor(
		context.Background(),
		component.ProcessorCreateParams{
			Logger: zap.NewNop(),
		},
		cfg,
		next,
	)
	assert.NotNil(t, mp)
	assert.Nil(t, err)

	mds := []internaldata.MetricsData{
		{
			Metrics: metricsWithName([]string{"metric1"}),
		},
	}
	cErr := mp.ConsumeMetrics(context.Background(), internaldata.OCSliceToMetrics(mds))
	assert.Nil(t, cErr)

	got := next.AllMetrics()
	assert.Equal(t, 1, len(got))

	v := got[0]
	assert.Equal(t, v.ResourceMetrics().At(0).InstrumentationLibraryMetrics().At(0).Metrics().At(0).Name(), "1cirtem")
	assert.Equal(t, v.ResourceMetrics().At(0).InstrumentationLibraryMetrics().At(0).Metrics().At(0).Description(), "csed_1cirtem")
}

type metricWithResource struct {
	metricNames []string
	resource    *resourcepb.Resource
}

func metricsWithName(names []string) []*metricspb.Metric {
	ret := make([]*metricspb.Metric, len(names))
	now := time.Now()
	for i, name := range names {
		ret[i] = &metricspb.Metric{
			MetricDescriptor: &metricspb.MetricDescriptor{
				Name:        name,
				Description: name + "_desc",
				Type:        metricspb.MetricDescriptor_GAUGE_INT64,
			},
			Timeseries: []*metricspb.TimeSeries{
				{
					Points: []*metricspb.Point{
						{
							Timestamp: timestamppb.New(now.Add(10 * time.Second)),
							Value: &metricspb.Point_Int64Value{
								Int64Value: int64(123),
							},
						},
					},
				},
			},
		}
	}
	return ret
}
