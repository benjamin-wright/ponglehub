package deployments

import (
	"fmt"
	"strings"
)

func mapHas(data map[string][]byte, keys []string) error {
	missing := []string{}
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing keys: %s", strings.Join(missing, ", "))
	}

	return nil
}
