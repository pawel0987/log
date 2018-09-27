// Author: Paweł Konopko
// License: MIT

package log

func joinData(data []Data) Data {
	result := Data{}

	// append data slice if exists
	if data != nil && len(data) != 0 {
		for _, someMap := range data {
			for key, value := range someMap {
				result[key] = value
			}
		}
	}

	return result
}
