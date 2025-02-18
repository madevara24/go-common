package mapper

import "context"

type IMapper interface {
	Insert(ctx context.Context, entity Entity, tableName string) (string, []interface{}, error)
	InsertMany(ctx context.Context, entities []Entity, tableName string) (string, []interface{}, error)

	Update(ctx context.Context, entity Entity, tableName string) (string, []interface{}, error)
	UpdateMany(ctx context.Context, entities []Entity, tableName string) (string, []interface{}, error)

	SoftDelete(tableName string, args SoftDeleteFilter) string
	HardDelete(tableName string, args SoftDeleteFilter) string
}
