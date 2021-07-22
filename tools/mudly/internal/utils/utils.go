package utils

func MergeMaps(maps ...map[string]string) map[string]string {
	output_map := map[string]string{}

	hasAny := false
	for _, obj := range maps {
		for key, value := range obj {
			output_map[key] = value
			hasAny = true
		}
	}

	if hasAny {
		return output_map
	} else {
		return nil
	}
}
