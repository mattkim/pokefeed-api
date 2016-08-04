package libuuid

// NewUUID generates a random UUID according to RFC 4122
func getMap(j, key string) (interface{}, error) {
	jType := j.(type)

	if jType != map[string]interface{} {
		// TODO: throw exception
		return nil, nil
	}

	// for k, v := range m {
	// 	switch vv := v.(type) {
	// 	case string:
	// 		Info.Println(k, "is string", vv)
	// 	case int:
	// 		Info.Println(k, "is int", vv)
	// 	case []interface{}:
	// 		Info.Println(k, "is an array:")
	// 		for i, u := range vv {
	// 			Info.Println(i, u)
	// 		}
	// 	case map[string]interface{}:
	// 		Info.Println(k, "is a map:")
	// 		for i, u := range vv {
	// 			Info.Println(i, u)
	// 		}
	// 	default:
	// 		Info.Println(k, "is of a type I don't know how to handle")
	// 		Info.Println(v)
	// 	}
	// }
	return j[key], nil
}
