package libuuid

// NewUUID generates a random UUID according to RFC 4122
func getMap(j, key string) (interface{}, error) {
	jType := j.(type)

	if jType != map[string]interface{} {
		// TODO: throw exception
		return nil, nil
	}
	// TODO: use typs.JSONText
	// TODO: saving this here becuase it was really hard to figure out
	// for _, feed := range feeds {
	// 	user, err := u.GetByUUID(nil, feed.CreatedByUserUUID)
	//
	// 	if err != nil {
	// 		libhttp.HandleBadRequest(w, err)
	// 		return
	// 	}
	// 	result := &GetLatestFeedsStruct{}
	//
	// 	// Find the right geocode here
	// 	var goodg map[string]interface{}
	// 	// var f []interface{}
	// 	var gs []interface{}
	// 	// json.Unmarshal(feed.Geocodes, &f)
	// 	err3 := json.Unmarshal(feed.Geocodes, &gs)
	// 	if gs == nil {
	// 		// TODO: this is super weird, but sometimes we cannot unmarshal geocodes even though it exists and is valid.
	// 		Info.Println("**** skipping")
	// 		Info.Println(fmt.Sprintf("%+v\n", feed))
	// 		Info.Println(gs)
	// 		Info.Println(err3)
	// 	} else {
	// 		// gs := f.([]interface{})
	// 		goodg = gs[0].(map[string]interface{}) // default to the first one.
	// 		for _, g := range gs {
	// 			gn := g.(map[string]interface{})
	// 			gnTypes := gn["types"].([]interface{})
	// 			for _, t := range gnTypes {
	// 				if t == feed.DisplayType {
	// 					goodg = gn
	// 				}
	// 			}
	// 		}
	//
	// 		// Info.Println(goodg)
	// 		// Fetch the formatted address and lat long here.
	// 		formattedAddress := goodg["formatted_address"].(string)
	// 		lat := goodg["geometry"].(map[string]interface{})["location"].(map[string]interface{})["lat"].(float64)
	// 		long := goodg["geometry"].(map[string]interface{})["location"].(map[string]interface{})["lng"].(float64)
	//
	// 		result.Username = user.Username
	// 		result.Message = feed.Message
	// 		result.Pokemon = feed.Pokemon
	// 		// TODO: fetch url from map.
	// 		result.CreatedAt = feed.CreatedAt.Time
	// 		result.Lat = lat
	// 		result.Long = long
	// 		result.FormattedAddress = formattedAddress
	// 		results = append(results, result)
	// 	}
	// }


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
