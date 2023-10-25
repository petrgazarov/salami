package metrics

import "time"

const ObjectsAdded = "objectsAdded"
const ObjectsRemoved = "objectsRemoved"
const ObjectsChanged = "objectsChanged"
const SourceFilesProcessed = "sourceFilesProcessed"

type Metrics struct {
	ObjectsChanged       int
	ObjectsAdded         int
	ObjectsRemoved       int
	SourceFilesProcessed int
	startTime            time.Time
}

var metrics *Metrics

func InitializeMetrics() {
	metrics = &Metrics{
		startTime: time.Now(),
	}
}

func GetMetrics() *Metrics {
	return metrics
}

func GetDuration() time.Duration {
	return time.Since(metrics.startTime)
}

func SetMetric(metricName string, value int) {
	switch metricName {
	case ObjectsChanged:
		metrics.ObjectsChanged = value
	case ObjectsAdded:
		metrics.ObjectsAdded = value
	case ObjectsRemoved:
		metrics.ObjectsRemoved = value
	case SourceFilesProcessed:
		metrics.SourceFilesProcessed = value
	default:
		panic("Unknown metric name: " + metricName)
	}
}
