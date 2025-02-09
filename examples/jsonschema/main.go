package main

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/kaptinlin/jsonschema"
	"io"
	"log"
)

func main() {
	objectSchema := `
{
  "properties": {
    "@id": {
      "type": "string"
    },
    "@schema": {
      "type": "string"
    },
    "@archetype": {
      "const": "unit"
    },
    "@kind": {
      "const": "INDIVIDUAL_ENTITY"
    },
    "@meta": {
  "type": "object",
  "properties": {
    "@createdAt": {
      "type": "string",
      "format": "date-time"
    },
    "@updatedAt": {
      "type": "string",
      "format": "date-time"
    },
    "@etag": {
      "type": "string"
    },
    "@with": {
      "type": "object",
      "additionalProperties": {
        "type": "array",
        "items": {
          "type": "string",
		  "format": "uuid"
        }
      }
    }
  },
  "required": [
    "@createdAt",
    "@updatedAt",
    "@etag"
  ],
  "unevaluatedProperties": false
},
    "@labels": {
      "properties": {
        "name": {
          "type": "string"
        }
      },
      "type": "object",
      "required": [
        "name"
      ]
    },
    "@forms": {
      "properties": {
        "INDIVIDUAL_ENTITY__SELF": {
          "properties": {
            "@object": {
              "type": "string",
              "x-tf-object-ref": {
                "json-ir://local@madesst/form/finance/INDIVIDUAL_ENTITY__SELF": {}
              }
            }
          },
          "type": "object",
          "required": [
            "@object"
          ]
        }
      },
      "type": "object",
      "unevaluatedProperties": false
    },
    "@context-forms": {},
    "@connections": {},
    "@context-connections": {}
  },
  "type": "object",
  "required": [
    "@id",
    "@schema",
    "@archetype",
    "@kind",
    "@meta",
    "@labels",
    "@forms",
    "@context-forms",
    "@connections",
    "@context-connections"
  ],
  "unevaluatedProperties": false
}
`

	compiler := jsonschema.NewCompiler()

	compiler.RegisterLoader("new", func(urlString string) (result io.ReadCloser, err error) {
		refObject := `
    {
      "@id": "new://INDIVIDUAL_ENTITY__SELF",
      "@schema": "json-ir://local@madesst/form/finance/INDIVIDUAL_ENTITY__SELF?12321312123123",
      "@archetype": "form",
      "@fields": {
        "name": "T"
      }
    }
`
		refSchemaBytes := []byte(refObject)
		return io.NopCloser(bytes.NewReader(refSchemaBytes)), nil
	})

	compiler.RegisterLoader("json-ir", func(urlString string) (result io.ReadCloser, err error) {
		refSchema := `
{
  "properties": {
    "@id": {
      "type": "string"
    },
    "@schema": {
      "type": "string"
    },
    "@archetype": {
      "const": "form"
    },
    "@kind": {
      "const": "INDIVIDUAL_ENTITY__SELF"
    },
    "@fields": {
      "properties": {
        "name": {
          "type": "string"
        }
      },
      "type": "object",
      "required": [
        "name"
      ],
      "unevaluatedProperties": false
    }
  },
  "type": "object",
  "required": [
    "@id",
    "@schema",
    "@archetype",
    "@kind",
    "@fields"
  ],
  "unevaluatedProperties": false
}
`
		refSchemaBytes := []byte(refSchema)
		return io.NopCloser(bytes.NewReader(refSchemaBytes)), nil
	})

	schema, err := compiler.Compile([]byte(objectSchema))
	if err != nil {
		log.Fatalf("Failed to compile schema: %v", err)
	}

	objectMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(`
{
  "@connections": {},
  "@context-connections": {},
  "@context-forms": {},
  "@forms": {
    "INDIVIDUAL_ENTITY__SELF": {
      "@object": "new://INDIVIDUAL_ENTITY__SELF"
    }
  },
  "@id": "tf://LOCAL@madesst/unit/c35ac8c6-ea79-4398-8ed7-521eea14a472",
  "@archetype": "unit",
  "@kind": "INDIVIDUAL_ENTITY",
  "@meta": {
    "@createdAt": "2021-08-31T07:00:00Z",
    "@updatedAt": "2021-08-31T07:00:00Z",
	"@etag": "c35ac8c6-ea79-4398-8ed7-521eea14a472",
    "@with": {
      "DIRECTOR": ["tf://LOCAL@madesst/unit/9767f7b0-1c25-4010-9138-40f26eee33bc"]
    }
  },
  "@labels": {
    "name": "Test Individual Entity"
  },
  "@schema": "json-ir://local@madesst/unit/finance/INDIVIDUAL_ENTITY"
}
`), &objectMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal object: %v", err)
	}

	validationResult := schema.Validate(objectMap)
	validationResultList := validationResult.ToList()
	log.Printf("Validation result: %v", validationResult.IsValid())
	log.Printf("Validation result details: %v", validationResultList)
}
