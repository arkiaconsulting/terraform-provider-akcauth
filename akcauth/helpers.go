package akcauth

func expandString(input []interface{}) []string {
	output := make([]string, 0)

	for _, s := range input {
		output = append(output, s.(string))
	}

	return output
}
