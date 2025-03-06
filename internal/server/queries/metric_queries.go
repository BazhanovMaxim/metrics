package queries

const (
	CreateMetricsTable = "CREATE TABLE IF NOT EXISTS Metrics\n(" +
		"ID	SERIAL PRIMARY KEY,\n" +
		"MiD   varchar(255),\n" +
		"MType varchar(255),\n" +
		"Delta int,\n" +
		"Value double precision\n" +
		");\n" +
		"CREATE UNIQUE INDEX unique_name_type_idx ON metrics (mid, mtype);"
	GetMetric    = "SELECT Metrics.mid, Metrics.mtype, Metrics.delta, Metrics.value FROM Metrics WHERE Metrics.mid = $1 and mtype = $1"
	GetMetrics   = "SELECT * FROM Metrics"
	InsertMetric = "INSERT INTO metrics (mid, mtype, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (mid, mtype) DO UPDATE SET delta = EXCLUDED.delta, value = EXCLUDED.value;"
)
