package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Entity interface {
	any
}

type PostgresMapper struct {
}

func NewPostgresMapper() IMapper {
	return &PostgresMapper{}
}

func (m *PostgresMapper) Insert(ctx context.Context, entity Entity, tableName string) (string, []interface{}, error) {
	var primaryKey string
	var insertQuery, insertValues string
	data := make([]interface{}, 0)
	val := reflect.Indirect(reflect.ValueOf(entity))

	// MAKE INSERT QUERY
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if field.Tag.Get("primarykey") == "true" {
			primaryKey = m.getFieldName(field.Tag, true)
		}
		if field.Tag.Get("insertable") != "false" {
			temp := m.formatData(val.FieldByName(field.Name).Interface())
			isZeroValue := val.Field(i).IsZero()
			isDBDefault := field.Tag.Get("hasdbdefault") == "true"
			if isDBDefault && isZeroValue {
				continue
			}
			fieldName := m.getFieldName(field.Tag, true)
			if insertQuery == "" {
				insertQuery = "INSERT INTO " + tableName + " (" + fieldName
				insertValues = ") VALUES (?"
			} else {
				insertQuery = insertQuery + "," + fieldName
				insertValues = insertValues + ",?"
			}

			data = append(data, temp)
		}
	}

	insertQuery = insertQuery + insertValues + ")"
	query := insertQuery

	if primaryKey != "" {
		query = insertQuery + " RETURNING " + tableName + "." + primaryKey
	}

	return query, data, nil
}

func (m *PostgresMapper) InsertMany(ctx context.Context, entities []Entity, tableName string) (string, []interface{}, error) {
	if len(entities) == 0 {
		return "", nil, fmt.Errorf("no entities to insert")
	}

	// MAKE BASE INSERT QUERY
	var primaryKey string
	fields := make([]string, 0)
	val := reflect.Indirect(reflect.ValueOf(entities[0]))
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if i == 0 && field.Tag.Get("primarykey") == "true" {
			primaryKey = m.getFieldName(field.Tag, true)
		}
		if field.Tag.Get("insertable") != "false" {
			fieldName := m.getFieldName(field.Tag, true)
			fields = append(fields, fieldName)
		}
	}

	values := make([]string, len(entities))
	bulkData := make([]interface{}, 0)

	// PROCESS FIELDS
	for idx, entity := range entities {
		val := reflect.Indirect(reflect.ValueOf(entity))
		singleValue := make([]string, 0, val.Type().NumField())
		paramStart := idx*val.Type().NumField() + 1

		for i := 0; i < val.Type().NumField(); i++ {
			field := val.Type().Field(i)
			if field.Tag.Get("insertable") != "false" {
				isZeroValue := val.Field(i).IsZero()
				isDBDefault := field.Tag.Get("hasdbdefault") == "true"

				if isDBDefault && isZeroValue {
					singleValue = append(singleValue, "DEFAULT")
				} else {
					singleValue = append(singleValue, fmt.Sprintf("$%d", paramStart))
					paramStart++
					bulkData = append(bulkData, m.formatData(val.Field(i).Interface()))
				}
			}
		}

		values[idx] = fmt.Sprintf("(%s)", strings.Join(singleValue, ","))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		tableName,
		strings.Join(fields, ","),
		strings.Join(values, ","))

	if primaryKey != "" {
		query += " RETURNING " + tableName + "." + primaryKey
	}

	return query, bulkData, nil
}

func (m *PostgresMapper) Update(ctx context.Context, entity Entity, tableName string) (string, []interface{}, error) {
	var key interface{}
	var updateQuery, match string
	data := make([]interface{}, 0)
	val := reflect.Indirect(reflect.ValueOf(entity))
	//Build the query
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if field.Tag.Get("updatable") == "true" { //Only if updatable true
			fieldName := m.getFieldName(field.Tag, true)
			if updateQuery == "" {
				updateQuery = "UPDATE " + tableName + " SET " + fieldName + "=?"
			} else {
				updateQuery = updateQuery + "," + fieldName + "=?"
			}

			temp := m.formatData(val.Field(i).Interface())
			data = append(data, temp)

		} else if field.Tag.Get("primarykey") == "true" {
			fieldName := m.getFieldName(field.Tag, true)
			match = " WHERE " + fieldName + "=?"
			key = m.formatData(val.Field(i).Interface())
		}
	}
	data = append(data, key)
	updateQuery = updateQuery + match

	return updateQuery, data, nil
}

func (m *PostgresMapper) getFieldName(tag reflect.StructTag, quoted bool) string {
	fieldName := strings.Split(tag.Get("json"), ",")[0]
	if tag.Get("db") != "" {
		fieldName = tag.Get("db")
	}
	if quoted {
		fieldName = strconv.Quote(fieldName)
	}

	return fieldName
}

func (m *PostgresMapper) formatData(val any) interface{} {
	if m.isDataType(val, "time") {
		t_byte, _ := json.Marshal(val)
		t := time.Time{}
		json.Unmarshal(t_byte, &t)
		if t.Unix() < 0 {
			val = time.Unix(0, 0).UTC()
		}
		if t.IsZero() {
			val = nil
		}
	}

	return val
}

func (m *PostgresMapper) isDataType(val any, pref string) bool {
	return strings.HasSuffix(strings.ToLower(reflect.TypeOf(val).String()), strings.ToLower(pref))
}
