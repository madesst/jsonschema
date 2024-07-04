package jsonschema

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// resolveRef resolves a reference to another schema, either locally or globally, supporting both $ref and $dynamicRef.
func (s *Schema) resolveRef(ref string) (*Schema, error) {
	if ref == "#" {
		return s.getRootSchema(), nil
	}

	if strings.HasPrefix(ref, "#") {
		return s.resolveAnchor(ref[1:])
	}

	// Resolve the full URL if ref is a relative URL
	if !isAbsoluteURI(ref) && s.baseURI != "" {
		ref = resolveRelativeURI(s.baseURI, ref)
	}

	// Handle full URL references
	return s.resolveRefWithFullURL(ref)
}

func (s *Schema) resolveAnchor(anchorName string) (*Schema, error) {
	var schema *Schema
	var err error

	if strings.HasPrefix(anchorName, "/") {
		schema, err = s.resolveJSONPointer(anchorName)
	} else {
		if schema, ok := s.anchors[anchorName]; ok {
			return schema, nil
		}

		if schema, ok := s.dynamicAnchors[anchorName]; ok {
			return schema, nil
		}
	}

	if schema == nil && s.parent != nil {
		return s.parent.resolveAnchor(anchorName)
	}

	return schema, err
}

// resolveRefWithFullURL resolves a full URL reference to another schema.
func (s *Schema) resolveRefWithFullURL(ref string) (*Schema, error) {
	root := s.getRootSchema()
	if resolved, err := root.getSchema(ref); err == nil {
		return resolved, nil
	}

	// If not found in the current schema or its parents, look for the reference in the compiler
	if resolved, err := s.compiler.GetSchema(ref); err != nil {
		return nil, errors.New(fmt.Sprintf("%s: %s", ErrFailedToResolveGlobalReference, ref))
	} else {
		return resolved, nil
	}
}

// resolveJSONPointer resolves a JSON Pointer within the schema based on JSON Schema structure.
func (s *Schema) resolveJSONPointer(pointer string) (*Schema, error) {
	if pointer == "/" {
		return s, nil
	}

	segments := strings.Split(strings.TrimPrefix(pointer, "/"), "/")
	currentSchema := s
	previousSegment := ""

	for i, segment := range segments {
		decodedSegment, err := url.PathUnescape(strings.ReplaceAll(strings.ReplaceAll(segment, "~1", "/"), "~0", "~"))
		if err != nil {
			return nil, ErrFailedToDecodeSegmentWithJSONPointer
		}

		nextSchema, found := findSchemaInSegment(currentSchema, decodedSegment, previousSegment)
		if found {
			currentSchema = nextSchema
			previousSegment = decodedSegment // Update the context for the next iteration
			continue
		}

		if !found && i == len(segments)-1 {
			// If no schema is found and it's the last segment, throw error
			return nil, ErrSegmentNotFoundForJSONPointer
		}

		previousSegment = decodedSegment // Update the context for the next iteration
	}

	return currentSchema, nil
}

// Helper function to find a schema within a given segment
func findSchemaInSegment(currentSchema *Schema, segment string, previousSegment string) (*Schema, bool) {
	switch previousSegment {
	case "properties":
		if currentSchema.Properties != nil {
			if schema, exists := (*currentSchema.Properties)[segment]; exists {
				return schema, true
			}
		}
	case "prefixItems":
		index, err := strconv.Atoi(segment)

		if err == nil && currentSchema.PrefixItems != nil && index < len(currentSchema.PrefixItems) {
			return currentSchema.PrefixItems[index], true
		}
	case "$defs":
		if defSchema, exists := currentSchema.Defs[segment]; exists {
			return defSchema, true
		}
	case "items":
		if currentSchema.Items != nil {
			return currentSchema.Items, true
		}
	}
	return nil, false
}

func (s *Schema) resolveReferences() (err error) {
	// Resolve the root reference if this schema itself is a reference
	if s.Ref != "" {
		s.ResolvedRef, err = s.resolveRef(s.Ref) // Resolve against root schema
		if err != nil {
			return
		}
	}

	if s.DynamicRef != "" {
		s.ResolvedDynamicRef, err = s.resolveRef(s.DynamicRef) // Resolve dynamic references against root schema
		if err != nil {
			return
		}
	}

	// Recursively resolve references within definitions
	if s.Defs != nil {
		for _, defSchema := range s.Defs {
			err = defSchema.resolveReferences()
			if err != nil {
				return
			}
		}
	}

	// Recursively resolve references in properties
	if s.Properties != nil {
		for _, schema := range *s.Properties {
			if schema != nil {
				err = schema.resolveReferences()
				if err != nil {
					return
				}
			}
		}
	}

	// Additional fields that can have subschemas
	err = resolveSubschemaList(s.AllOf)
	if err != nil {
		return
	}
	err = resolveSubschemaList(s.AnyOf)
	if err != nil {
		return
	}
	err = resolveSubschemaList(s.OneOf)
	if err != nil {
		return
	}
	if s.Not != nil {
		err = s.Not.resolveReferences()
		if err != nil {
			return
		}
	}
	if s.Items != nil {
		err = s.Items.resolveReferences()
		if err != nil {
			return
		}
	}
	if s.PrefixItems != nil {
		for _, schema := range s.PrefixItems {
			err = schema.resolveReferences()
			if err != nil {
				return
			}
		}
	}

	if s.AdditionalProperties != nil {
		err = s.AdditionalProperties.resolveReferences()
		if err != nil {
			return
		}
	}
	if s.Contains != nil {
		err = s.Contains.resolveReferences()
		if err != nil {
			return
		}
	}
	if s.PatternProperties != nil {
		for _, schema := range *s.PatternProperties {
			err = schema.resolveReferences()
			if err != nil {
				return
			}
		}
	}
	return
}

// Helper function to resolve references in a list of schemas
func resolveSubschemaList(schemas []*Schema) (err error) {
	for _, schema := range schemas {
		if schema != nil {
			err = schema.resolveReferences()
			if err != nil {
				return
			}
		}
	}
	return
}
