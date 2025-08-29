package plugin

import (
	"context"

	"github.com/databricks/databricks-sdk-go/listing"
)

func fetchWithLimit[T any](ctx context.Context, it listing.Iterator[T], maxItems int) ([]T, error) {
	var result = []T{}

	for it.HasNext(ctx) && len(result) < maxItems {
		item, err := it.Next(ctx)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}
