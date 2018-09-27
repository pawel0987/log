// Author: Paweł Konopko
// License: MIT

package log

import "github.com/pawel0987/utils/json_utils"

func formatJsonMessage (fieldsOrder []string, fields map[string]interface{}) string {
	result := "{\""
	for _, key := range fieldsOrder {
		result += key + "\":" + json_utils.Encode(fields[key]) + ",\""
	}

	return result[0:len(result)-2] + "}\n"
}
