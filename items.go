package jsonschema

import (
	"fmt"
	"strconv"
	"strings"
)

// EvaluateItems checks if the data's array items conform to the subschema or boolean condition specified in the 'items' attribute of the schema.
// According to the JSON Schema Draft 2020-12:
//   - The value of "items" MUST be either a valid JSON Schema or a boolean.
//   - If "items" is a Schema, each element of the instance array must conform to this subschema.
//   - If "items" is boolean and is true, any array elements are valid.
//   - If "items" is boolean and is false, no array elements are valid unless the array is empty.
//
// This method ensures that array elements conform to the constraints defined in the items attribute.
// If any array element does not conform, it returns a EvaluationError detailing the issue.
//
// Reference: https://json-schema.org/draft/2020-12/json-schema-core#name-items
func evaluateItems(schema *Schema, array []interface{}, evaluatedProps map[string]bool, evaluatedItems map[int]bool, dynamicScope *DynamicScope) ([]*EvaluationResult, *EvaluationError) {
	if schema.Items == nil {
		return nil, nil // // No 'items' constraints to validate against
	}

	invalid_indexs := []string{}
	results := []*EvaluationResult{}

	// Number of prefix items to skip before regular item validation
	startIndex := len(schema.PrefixItems)

	// Check if the general 'items' schema is available and proceed with validation if it's not explicitly false
	if schema.Items != nil {
		// Ensure that we only access indices within the range of existing array elements
		for i := startIndex; i < len(array); i++ {
			item := array[i]
			result, _, _ := schema.Items.evaluate(item, dynamicScope)
			if result != nil {
				anchor := fmt.Sprintf("/items/%d", i)
				instanceLocation := fmt.Sprintf("/%d", i)

				result.SetEvaluationPath(anchor).
					SetSchemaLocation(schema.GetSchemaLocation(anchor)).
					SetInstanceLocation(instanceLocation)

				result.Details = recursiveUpdateDetailsWithItemAnchor(schema, result.Details, anchor, instanceLocation)

				if result.IsValid() {
					evaluatedItems[i] = true // Mark the item as evaluated if it passes schema validation.
				} else {
					invalid_indexs = append(invalid_indexs, strconv.Itoa(i))
					results = append(results, result)
				}
			}
		}
	}

	if len(invalid_indexs) == 1 {
		return results, NewEvaluationError("items", "item_mismatch", "Item at index {index} does not match the schema", map[string]interface{}{
			"index": invalid_indexs[0],
		})
	} else if len(invalid_indexs) > 1 {
		return results, NewEvaluationError("items", "items_mismatch", "Items at index {indexs} do not match the schema", map[string]interface{}{
			"indexs": strings.Join(invalid_indexs, ", "),
		})
	}
	return results, nil
}

func recursiveUpdateDetailsWithItemAnchor(schema *Schema, details []*EvaluationResult, anchor string, instanceLocation string) []*EvaluationResult {
	for key, value := range details {
		originEvalPath := value.EvaluationPath
		value.SetEvaluationPath(fmt.Sprintf("%s%s", anchor, originEvalPath)).
			SetSchemaLocation(schema.GetSchemaLocation(fmt.Sprintf("%s%s", anchor, originEvalPath))).
			SetInstanceLocation(fmt.Sprintf("%s%s", instanceLocation, value.InstanceLocation))
		value.Details = recursiveUpdateDetailsWithItemAnchor(schema, value.Details, anchor, instanceLocation)
		details[key] = value
	}
	return details
}
