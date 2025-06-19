package db

type TagRecord struct {
	ID           int64  `db:"id"`
	FullTag      string `db:"full_tag"`
	SystemCode   string `db:"system_code"`
	EquipmentID  string `db:"equipment_id"`
	InstrumentID string `db:"instrument_id"`
	FunctionCode string `db:"function_code"`
	UDCCode      string `db:"udc_code"`
	Description  string `db:"description"`
}
