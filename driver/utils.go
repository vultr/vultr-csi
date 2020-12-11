package driver

import (
	"context"
	"fmt"

	"github.com/vultr/govultr/v2"
)

func getUUID(ctx context.Context, client *govultr.Client, v1ID string) (string, error) {
	id := fmt.Sprintf("v1-%s", v1ID)
	instance, err := client.Instance.Get(ctx, id)
	if err != nil {
		return "", err
	}

	return instance.ID, nil
}
