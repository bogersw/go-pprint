package pprint

import (
	"fmt"
	"reflect"
	"strconv"
)

type PPColor string

const (
	colorReset  PPColor = "\033[0m"
	colorRed    PPColor = "\033[31m"
	colorGreen  PPColor = "\033[32m"
	colorYellow PPColor = "\033[33m"
	colorBlue   PPColor = "\033[34m"
)

// ============================================================================
// PUBlIC FUNCTIONS
// ============================================================================

func PPrintRed(value interface{}) {
	fmt.Print(formatValue(value, colorRed, false))
}

func PPrintGreen(value interface{}) {
	fmt.Print(formatValue(value, colorGreen, false))
}

func PPrintYellow(value interface{}) {
	fmt.Print(formatValue(value, colorYellow, false))
}

func PPrintBlue(value interface{}) {
	fmt.Print(formatValue(value, colorBlue, false))
}

// ============================================================================
// PRIVATE FUNCTIONS
// ============================================================================

func formatValue(value interface{}, color PPColor, inline bool) string {

	// Check if value is an array, slice, map, or struct (special cases)
	isArray := reflect.ValueOf(value).Kind() == reflect.Array
	if isArray {
		return formatList(value, color, "Array", inline)
	}
	isSlice := reflect.ValueOf(value).Kind() == reflect.Slice
	if isSlice {
		return formatList(value, color, "Slice", inline)
	}
	isMap := reflect.ValueOf(value).Kind() == reflect.Map
	if isMap {
		return formatMap(value, color, inline)
	}
	isStruct := reflect.ValueOf(value).Kind() == reflect.Struct
	if isStruct {
		return formatStruct(value, color, inline)
	}
	// General case: just print tthe value
	if inline {
		return fmt.Sprintf("%#v", value)
	} else {
		valueType := reflect.ValueOf(value).Type().String()
		return fmt.Sprintf("%#v (%s%s%s)\n\n", value, color, valueType, colorReset)
	}
}

func formatList(value interface{}, color PPColor, listType string, inline bool) string {

	// A list can be an array or a slice. Determine the type of the items in
	// the list and the maximum length of the list index as a string.
	list := reflect.ValueOf(value)
	elType := list.Type().Elem()
	colWidth := len(strconv.Itoa(list.Len()))
	// Pretty print the list (index and value)
	output := ""
	if inline {
		output += "{"
	} else {
		// Add a header if the list is not printed inline
		output += fmt.Sprintf("%s with %s%d%s elements of type %s%s%s:\n\n",
			listType, color, list.Len(), colorReset, color, elType.String(), colorReset)
	}
	for i := 0; i < list.Len(); i++ {
		item := formatValue(list.Index(i).Interface(), color, true)
		if inline {
			// Inline items are printed on a single line, separated by commas
			output += fmt.Sprintf("%s[%d]%s %s", color, i, colorReset, item)
			if i < list.Len()-1 {
				output += ", "
			}
		} else {
			// Non-inline items are printed on multiple lines
			output += fmt.Sprintf("%s[%*d]%s %v\n", color, colWidth, i, colorReset, item)
		}
	}
	if inline {
		output += "}"
	} else {
		output += "\n"
	}
	return output
}

func formatMap(value interface{}, color PPColor, inline bool) string {

	// Determine types of keys and values
	mp := reflect.ValueOf(value)
	valueType := mp.Type().Elem().Name()
	keys := mp.MapKeys()
	keyType := keys[0].Type()
	// Pretty print the map
	output := ""
	if inline {
		output += "{"
	} else {
		// Add a header if the map is not printed inline
		output += fmt.Sprintf("Map: %s%d%s elements, key type %s%s%s, value type %s%s%s:\n\n",
			color, mp.Len(), colorReset, color, keyType, colorReset, color, valueType, colorReset)
	}
	for index, key := range keys {
		item := formatValue(mp.MapIndex(key).Interface(), color, true)
		if inline {
			// Inline items are printed on a single line, separated by commas
			output += fmt.Sprintf("%s[%#v]%s %s", color, key, colorReset, item)
			if index < len(keys)-1 {
				output += ", "
			}
		} else {
			output += fmt.Sprintf("%s[%#v]%s %v\n", color, key, colorReset, item)
		}
	}
	if inline {
		output += "}"
	} else {
		output += "\n"
	}
	return output
}

func formatStruct(value interface{}, color PPColor, inline bool) string {

	// Get value as struct
	s := reflect.ValueOf(value)
	// Get type information
	t := reflect.TypeOf(value)
	// Determine total number of fields
	numFields := s.NumField()
	// Determine length of longest field name for aligning
	maxNameLength := 0
	for i := range numFields {
		if len(t.Field(i).Name) > maxNameLength {
			maxNameLength = len(t.Field(i).Name)
		}
	}
	// Construct output string
	output := ""
	if !inline {
		// Add a header if the struct is not printed inline
		output += fmt.Sprintf("Struct with %s%d%s fields:\n\n",
			color, numFields, colorReset)
	}
	// Loop over fields and print name / type / value
	output += "{"
	if !inline {
		// Add fields on new lines if not printed inline
		output += "\n"
	}
	for i := range numFields {
		fieldName := t.Field(i).Name
		fieldValue := s.Field(i)
		// fieldType := s.Field(i).Type().String()
		if inline {
			output += fmt.Sprintf("%s%s%s: %#v", color, fieldName, colorReset, fieldValue)
			if i < numFields-1 {
				output += ", "
			}
		} else {
			output += fmt.Sprintf("  %s%*s%s: %#v\n",
				color, maxNameLength, fieldName, colorReset, fieldValue)
		}
	}
	output += "}"
	if !inline {
		// Add space below the struct if not printed inline
		output += "\n\n"
	}
	return output
}
