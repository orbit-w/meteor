package subpub_redis

import (
	"fmt"
)

func ErrPublish(err error) error {
	return fmt.Errorf("subpub publish failed: %w", err)
}
