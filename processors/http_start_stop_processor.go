package processors

import (
	"strconv"
	"strings"

	"github.com/cloudcredo/graphite-nozzle/metrics"
	"github.com/cloudfoundry/noaa/events"
)

type HttpStartStopProcessor struct{}

func NewHttpStartStopProcessor() *HttpStartStopProcessor {
	return &HttpStartStopProcessor{}
}

func (p *HttpStartStopProcessor) Process(e *events.Envelope) []metrics.Metric {
	processedMetrics := make([]metrics.Metric, 2)
	httpStartStopEvent := e.GetHttpStartStop()

	processedMetrics[0] = metrics.Metric(p.ProcessHttpStartStopResponseTime(httpStartStopEvent))
	processedMetrics[1] = metrics.Metric(p.ProcessHttpStartStopStatusCodeCount(httpStartStopEvent))

	return processedMetrics
}

func (p *HttpStartStopProcessor) ProcessHttpStartStopResponseTime(event *events.HttpStartStop) *metrics.TimingMetric {
	statPrefix := "http.responsetimes."
	hostname := strings.Replace(strings.Split(event.GetUri(), "/")[0], ".", "_", -1)
	stat := statPrefix + hostname

	startTimestamp := event.GetStartTimestamp()
	stopTimestamp := event.GetStopTimestamp()
	durationNanos := stopTimestamp - startTimestamp
	durationMillis := durationNanos / 1000000 // NB: loss of precision here
	metric := metrics.NewTimingMetric(stat, durationMillis)

	return metric
}

func (p *HttpStartStopProcessor) ProcessHttpStartStopStatusCodeCount(event *events.HttpStartStop) *metrics.CounterMetric {
	var incrementValue int64

	statPrefix := "http.statuscodes."
	hostname := strings.Replace(strings.Split(event.GetUri(), "/")[0], ".", "_", -1)
	stat := statPrefix + hostname + "." + strconv.Itoa(int(event.GetStatusCode()))

	if event.GetPeerType() == events.PeerType_Client {
		incrementValue = 1
	} else {
		incrementValue = 0
	}

	metric := metrics.NewCounterMetric(stat, incrementValue)

	return metric
}
