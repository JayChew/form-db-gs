package models

type FormModel struct {
	ID            int     `db:"id" json:"id" col:"Id"`
	CharCol       string  `db:"char_col" json:"char_col" col:"Char Column"`
	VarcharCol    string  `db:"varchar_col" json:"varchar_col" col:"Varchar Column"`
	TextCol       string  `db:"text_col" json:"text_col" col:"Text Column"`
	MediumtextCol string  `db:"mediumtext_col" json:"mediumtext_col" col:"Mediumtext Column"`
	LongtextCol   string  `db:"longtext_col" json:"longtext_col" col:"Longtext Column"`
	TinyintCol    int8    `db:"tinyint_col" json:"tinyint_col" col:"Tinyint Column"`
	SmallintCol   int16   `db:"smallint_col" json:"smallint_col" col:"Smallint Column"`
	MediumintCol  int32   `db:"mediumint_col" json:"mediumint_col" col:"Mediumint Column"`
	BigintCol     int64   `db:"bigint_col" json:"bigint_col" col:"Bigint Column"`
	DecimalCol    float64 `db:"decimal_col" json:"decimal_col" col:"Decimal Column"`
	FloatCol      float32 `db:"float_col" json:"float_col" col:"Float Column"`
	DoubleCol     float64 `db:"double_col" json:"double_col" col:"Double Column"`
	DateCol       string  `db:"date_col" json:"date_col" col:"Date Column"`                // or time.Time if using time package
	DatetimeCol   string  `db:"datetime_col" json:"datetime_col" col:"Datetime Column"`    // or time.Time if using time package
	TimestampCol  string  `db:"timestamp_col" json:"timestamp_col" col:"Timestamp Column"` // or time.Time if using time package
	TimeCol       string  `db:"time_col" json:"time_col" col:"Time Column"`                // or time.Time if using time package
	YearCol       int     `db:"year_col" json:"year_col" col:"Year Column"`
	EnumCol       string  `db:"enum_col" json:"enum_col" col:"Enum Column"`
	SetCol        string  `db:"set_col" json:"set_col" col:"Set Column"`
	JsonCol       string  `db:"json_col" json:"json_col" col:"Json Column"`
	Url           string  `db:"url" json:"url" col:"Url"`
}

type FormInterface interface {
	// Create(FormModel) (int64, error)
	GetAll() ([]FormModel, error)
}
