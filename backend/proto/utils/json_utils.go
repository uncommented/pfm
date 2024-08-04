package utils

import (
	"log"
	"math"
	"strconv"
)

func UnmarshalToInt(json map[string]interface{}, key string) int64 {
	if value, ok := json[key]; ok {
		if value_str, ok := value.(string); ok {
			value_int, err := strconv.ParseInt(value_str, 10, 64)
			if err != nil {
				log.Printf("Failed to umarshal %v to int64 (key: %s)", value_str, key)
				return -math.MaxInt64
			}
			return value_int
		} else {
			log.Printf("Failed to umarshal %v to string (key: %s)", value, key)
			if value_int, ok := value.(int64); ok {
				return value_int
			} else {
				log.Printf("Failed to parse %v to int64", value)
				return -math.MaxInt64
			}
		}
	} else {
		log.Printf("KeyError: %s (%v)", key, json)
		return -math.MaxInt64
	}
}

func UnmarshalToFloat(json map[string]interface{}, key string) float64 {
	if value, ok := json[key]; ok {
		if value_str, ok := value.(string); ok {
			value_float, err := strconv.ParseFloat(value_str, 64)
			if err != nil {
				log.Printf("Failed to umarshal %v to float64 (key: %s)", value_str, key)
				return -math.MaxFloat64
			}
			return value_float
		} else {
			log.Printf("Failed to umarshal %v to string (key: %s)", value, key)
			if value_float, ok := value.(float64); ok {
				return value_float
			} else {
				log.Printf("Failed to parse %v to float64", value)
				return -math.MaxFloat64
			}
		}
	} else {
		log.Printf("KeyError: %s (%v)", key, json)
		return -math.MaxFloat64
	}
}

func UnmarshalToString(json map[string]interface{}, key string) string {
	if value, ok := json[key]; ok {
		if value_str, ok := value.(string); ok {
			return value_str
		} else {
			log.Printf("Failed to umarshal %v to string (key: %s)", value, key)
			return "<none>"
		}
	} else {
		log.Printf("KeyError: %s (%v)", key, json)
		return "<none>"
	}
}

func UnmarshalToMap(json map[string]interface{}, key string) map[string]interface{} {
	if value, ok := json[key]; ok {
		if value_map, ok := value.(map[string]interface{}); ok {
			return value_map
		} else {
			log.Printf("Failed to umarshal %v to map (key: %s)", value, key)
			return make(map[string]interface{})
		}
	} else {
		log.Printf("KeyError: %s (%v)", key, json)
		return make(map[string]interface{})
	}
}

func UnmarshalToList(json map[string]interface{}, key string) []interface{} {
	if value, ok := json[key]; ok {
		if value_list, ok := value.([]interface{}); ok {
			return value_list
		} else {
			log.Printf("Failed to umarshal %v to list (key: %s)", value, key)
			return []interface{}{}
		}
	} else {
		log.Printf("KeyError: %s (%v)", key, json)
		return []interface{}{}
	}
}
