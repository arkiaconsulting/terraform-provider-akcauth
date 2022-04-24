package akcauth

func expandString(input []interface{}) []string {
	output := make([]string, len(input))

	for _, s := range input {
		output = append(output, s.(string))
	}

	return output
}
