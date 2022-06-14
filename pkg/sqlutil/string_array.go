package sqlutil

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type StringArray []string

func (a *StringArray) Contains(str string) bool {
	for _, item := range *a {
		if strings.EqualFold(item, str) {
			return true
		}
	}

	return false
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	var result []string
	err := json.Unmarshal(bytes, &result)
	*a = result

	return err
}

// Value return json value, implement driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if len(a) == 0 {
		return []byte("[]"), nil
	}

	return json.Marshal(a)
}

func (StringArray) GormDataType() string {
	return "json"
}

func (StringArray) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
