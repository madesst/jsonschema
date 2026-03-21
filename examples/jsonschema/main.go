package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"github.com/kaptinlin/jsonschema"
	"io"
	"log"
)

//go:embed *.json
var testFiles embed.FS

func main() {
	compiler := jsonschema.NewCompiler()

	compiler.RegisterLoader("json-ir", func(urlString string) (result io.ReadCloser, err error) {
		var schemaBytes []byte
		schemaBytes, err = testFiles.ReadFile("schema.json")
		return io.NopCloser(bytes.NewReader(schemaBytes)), nil
	})
	objectSchemaURI := "json-ir://00000000-0000-0000-0000-000000000000@madesst/profile/finance/CORPORATE_ACCOUNT?fc2b8ca2-b6b8-43f8-a187-87b3352d5e28"
	schema, err := compiler.GetSchema(objectSchemaURI)
	if err != nil {
		log.Fatalf("Failed to compile schema: %v", err)
	}

	objectMap := map[string]interface{}{}
	var objectBytes []byte
	objectBytes, err = testFiles.ReadFile("object.json")
	err = json.Unmarshal(objectBytes, &objectMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal object: %v", err)
	}

	validationResult := schema.Validate(objectMap)
	validationResultList := validationResult.ToList(true, true)
	log.Printf("Validation result: %v", validationResult.IsValid())
	validationResultListJson, _ := json.MarshalIndent(validationResultList, "", "  ")
	log.Printf(string(validationResultListJson))
}
