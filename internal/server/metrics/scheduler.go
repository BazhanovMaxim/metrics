package metrics

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"time"
)

func StartPeriodicSave(config configs.Config, service service.MetricService) {
	ticker := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
	defer ticker.Stop()

	// Канал для остановки горутины
	stopChan := make(chan struct{})

	// Горутина для периодического сохранения данных
	for {
		select {
		case <-ticker.C:
			service.SaveMetricsToStorage()
		case <-stopChan:
			return
		}
	}
}
