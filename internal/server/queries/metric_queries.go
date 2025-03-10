package queries

const (
	// CreateMetricsTable создает таблицу для работы с метриками
	CreateMetricsTable = "CREATE TABLE IF NOT EXISTS Metrics(ID	SERIAL PRIMARY KEY, MiD varchar(255), MType varchar(255), Delta bigint, Value double precision);"

	// CreateIndex создает индекс для mid, mtype
	CreateIndex = "CREATE UNIQUE INDEX unq_name_type_idx ON metrics (mid, mtype);"

	// GetMetric получает метрику из базы данных по параметрам(mid, mtype)
	GetMetric = "SELECT Metrics.mid, Metrics.mtype, Metrics.delta, Metrics.value FROM Metrics WHERE Metrics.mid = $1 and mtype = $2"

	// GetMetrics получает все метрики из базы данных
	GetMetrics = "SELECT Metrics.mid, Metrics.mtype, Metrics.delta, Metrics.value FROM Metrics"

	// InsertMetric добавляет метрику в базу данных.
	// В случае конфликта изменяет значение метрики на новое.
	InsertMetric = "INSERT INTO metrics (mid, mtype, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (mid, mtype) DO UPDATE SET delta = COALESCE(metrics.delta, 0) + EXCLUDED.delta, value = EXCLUDED.value;"
)
