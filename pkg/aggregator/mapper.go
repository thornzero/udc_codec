package aggregator

func (db *AggregatedDatabase) LookupSystem(systemCode string) *AggregatedSystem {
	for _, sys := range db.Systems {
		if sys.SystemCode == systemCode {
			return &sys
		}
	}
	return nil
}
