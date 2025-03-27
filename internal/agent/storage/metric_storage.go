package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/model"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"math/rand"
	"runtime"
	"sync"
)

type MetricStorage struct {
	mutex   sync.Mutex
	storage map[string]model.Metric
}

// NewMetricRepository создает и возвращает новый экземпляр IMetricStorage
func NewMetricRepository() IMetricStorage {
	return &MetricStorage{storage: generateMetrics(1)}
}

// GetMetrics возвращает собранные метрики
func (metrics *MetricStorage) GetMetrics() map[string]model.Metric {
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()
	return metrics.storage
}

// Update обновляет системные метрики
func (metrics *MetricStorage) Update() {
	metrics.mutex.Lock()
	metrics.storage = generateMetrics(metrics.storage["PollCount"].Value.(int64) + 1)
	metrics.mutex.Unlock()
}

// UpdateSys обновляет системные настройки памяти
func (metrics *MetricStorage) UpdateSys() {
	metrics.mutex.Lock()

	mt := metrics.GetMetrics()
	v, _ := mem.VirtualMemory()

	// Get CPU utilization
	cpuPercent, _ := cpu.Percent(0, false)
	cpuUtilization1 := cpuPercent[0]

	mt["TotalMemory"] = model.Metric{Value: v.Total, Type: model.Gauge}
	mt["FreeMemory"] = model.Metric{Value: v.Free, Type: model.Gauge}
	mt["CPUutilization1"] = model.Metric{Value: cpuUtilization1, Type: model.Gauge}

	metrics.mutex.Unlock()
}

// generateMetrics генерирует системный метрики
func generateMetrics(pollCount int64) map[string]model.Metric {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return map[string]model.Metric{
		"Alloc":         {Value: float64(memStats.Alloc), Type: model.Gauge},
		"BuckHashSys":   {Value: float64(memStats.BuckHashSys), Type: model.Gauge},
		"Frees":         {Value: float64(memStats.Frees), Type: model.Gauge},
		"GCCPUFraction": {Value: memStats.GCCPUFraction, Type: model.Gauge},
		"GCSys":         {Value: float64(memStats.GCSys), Type: model.Gauge},
		"HeapAlloc":     {Value: float64(memStats.HeapAlloc), Type: model.Gauge},
		"HeapIdle":      {Value: float64(memStats.HeapIdle), Type: model.Gauge},
		"HeapInuse":     {Value: float64(memStats.HeapInuse), Type: model.Gauge},
		"HeapObjects":   {Value: float64(memStats.HeapObjects), Type: model.Gauge},
		"HeapReleased":  {Value: float64(memStats.HeapReleased), Type: model.Gauge},
		"HeapSys":       {Value: float64(memStats.HeapSys), Type: model.Gauge},
		"LastGC":        {Value: float64(memStats.LastGC), Type: model.Gauge},
		"Lookups":       {Value: float64(memStats.Lookups), Type: model.Gauge},
		"MCacheInuse":   {Value: float64(memStats.MCacheInuse), Type: model.Gauge},
		"MCacheSys":     {Value: float64(memStats.MCacheSys), Type: model.Gauge},
		"MSpanInuse":    {Value: float64(memStats.MSpanInuse), Type: model.Gauge},
		"MSpanSys":      {Value: float64(memStats.MSpanSys), Type: model.Gauge},
		"Mallocs":       {Value: float64(memStats.Mallocs), Type: model.Gauge},
		"NextGC":        {Value: float64(memStats.NextGC), Type: model.Gauge},
		"NumForcedGC":   {Value: float64(memStats.NumForcedGC), Type: model.Gauge},
		"NumGC":         {Value: float64(memStats.NumGC), Type: model.Gauge},
		"OtherSys":      {Value: float64(memStats.OtherSys), Type: model.Gauge},
		"PauseTotalNs":  {Value: float64(memStats.PauseTotalNs), Type: model.Gauge},
		"StackInuse":    {Value: float64(memStats.StackInuse), Type: model.Gauge},
		"StackSys":      {Value: float64(memStats.StackSys), Type: model.Gauge},
		"Sys":           {Value: float64(memStats.Sys), Type: model.Gauge},
		"TotalAlloc":    {Value: float64(memStats.TotalAlloc), Type: model.Gauge},
		"PollCount":     {Value: pollCount, Type: model.Counter},
		"RandomValue":   {Value: rand.Float64(), Type: model.Gauge},
	}
}
