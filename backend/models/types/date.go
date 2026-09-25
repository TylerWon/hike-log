package types

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// Date is a custom gorm type that stores a calendar day as a YYYY-MM-DD string. Only at the database level is it
// converted to/from a time.Time. Without this type, a date would have to use time.Time (a timestamp) which does not
// convert easily to a date.
type Date string

// Serializes a Date into a byte array. On success, the byte array will store a YYYY-MM-DD string (including start and
// end quotes).
func (date Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + string(date) + `"`), nil
}

// Deserializes a Date from a byte array. In order to succeed, the byte array must store a YYYY-MM-DD string.
func (date *Date) UnmarshalJSON(b []byte) error {
	// The string parsed from b will have quotes in it so the layout string needs to also include quotes for it to
	// match the string and parse correctly
	timestamp, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}

	dateString := timestamp.Format("2006-01-02")
	*date = Date(dateString)

	return nil
}

// Converts a value from the database to a Date.
func (date *Date) Scan(value interface{}) error {
	// Parse database value into Go type
	nullTime := &sql.NullTime{}
	err := nullTime.Scan(value)
	if err != nil {
		return err
	} else if !nullTime.Valid {
		// DB value NULL, set Date to empty string
		*date = ""
		return nil
	}

	dateString := nullTime.Time.Format("2006-01-02")
	*date = Date(dateString)

	return nil
}

// Converts a Date to a value that can be stored in the database.
func (date Date) Value() (driver.Value, error) {
	if date == "" {
		// Date is empty, return nil so NULL is stored in the DB
		return nil, nil
	}

	timestamp, err := time.Parse("2006-01-02", string(date))
	if err != nil {
		return nil, err
	}

	return timestamp, nil
}

// Tells GORM how to map the Date type to a SQL column type. Returns "date" so the column is day only, no time.
func (date Date) GormDataType() string {
	return "date"
}
