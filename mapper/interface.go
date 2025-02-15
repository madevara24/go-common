package mapper

import "context"

type IMapper interface {
	InsertMany(ctx context.Context, entities []Entity, tableName string) (string, []interface{}, error)
}
