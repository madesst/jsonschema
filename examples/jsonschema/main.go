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
  "definitions": {
    "inward_input": {
      "type": "object",
      "properties": {
        "payment_type": {
          "type": "string",
          "const": "inward"
        },
        "transaction_amount": {
          "type": "number"
        },
        "sender_jurisdiction": {
          "type": "string"
        },
        "sender_bank_jurisdiction": {
          "type": "string"
        },
        "payment_purpose": {
          "type": "string"
        }
      },
      "required": [
        "payment_type",
        "transaction_amount",
        "sender_jurisdiction",
        "sender_bank_jurisdiction",
        "payment_purpose"
      ]
    },
    "outward_input": {
      "type": "object",
      "properties": {
        "payment_type": {
          "type": "string",
          "const": "outward"
        },
        "transaction_amount": {
          "type": "number"
        },
        "recipient_jurisdiction": {
          "type": "string"
        },
        "recipient_bank_jurisdiction": {
          "type": "string"
        },
        "payment_purpose": {
          "type": "string"
        },
        "has_supporting_docs": {
          "type": "boolean"
        }
      },
      "required": [
        "payment_type",
        "transaction_amount",
        "recipient_jurisdiction",
        "recipient_bank_jurisdiction",
        "payment_purpose",
        "has_supporting_docs"
      ]
    },
    "indicator": {
      "type": "object",
      "properties": {
        "risk_level": {
          "type": "string",
          "enum": [
            "LOW",
            "MEDIUM",
            "HIGH"
          ]
        },
        "calculation": {
          "type": "object",
          "properties": {
            "base_score": {
              "type": "number"
            },
            "factor": {
              "type": "number"
            },
            "group_weight": {
              "type": "number"
            },
            "total_score": {
              "type": "number"
            }
          },
          "required": [
            "base_score",
            "factor",
            "group_weight",
            "total_score"
          ]
        }
      },
      "required": [
        "risk_level",
        "calculation"
      ]
    }
  },
  "allOf": [
    {

    },
    {
      "type": "object",
      "properties": {
        "@kind": {
          "const": "transaction_assessment_result"
        },
        "@data": {
          "allOf": [
            {
              
            },
            {
              "type": "object",
              "properties": {
                "initial_input": {
                  "allOf": [
                    {
                      "oneOf": [
                        {
                          "$ref": "#/definitions/inward_input"
                        },
                        {
                          "$ref": "#/definitions/outward_input"
                        }
                      ]
                    },
                    {
                      "type": "object",
                      "properties": {
                        "statistics": {
                          "type": "object",
                          "properties": {
                            "avg_volume": {
                              "type": "number"
                            },
                            "avg_frequency": {
                              "type": "number"
                            }
                          },
                          "required": [
                            "avg_volume",
                            "avg_frequency"
                          ]
                        },
                        "declared": {
                          "type": "object",
                          "properties": {
                            "avg_volume": {
                              "type": "string",
                              "enum": [
                                "UNDER_10K",
                                "10K_TO_100K",
                                "ABOVE_100K"
                              ]
                            },
                            "avg_frequency": {
                              "type": "string",
                              "enum": [
                                "UNDER_6",
                                "6_TO_20",
                                "ABOVE_20"
                              ]
                            }
                          },
                          "required": [
                            "avg_volume",
                            "avg_frequency"
                          ]
                        }
                      },
                      "required": [
                        "statistics",
                        "declared"
                      ]
                    }
                  ]
                },
                "calculation_log": {
                  "type": "object",
                  "properties": {
                    "risk_level": {
                      "type": "string",
                      "enum": [
                        "LOW",
                        "MEDIUM",
                        "HIGH"
                      ]
                    },
                    "risk_score": {
                      "type": "number"
                    },
                    "manual_review_needed": {
                      "type": "boolean"
                    },
                    "actions": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      }
                    },
                    "transaction_amount_and_frequency": {
                      "type": "object",
                      "properties": {
                        "transaction_amount": {
                          "$ref": "#/definitions/indicator"
                        },
                        "transaction_frequency": {
                          "$ref": "#/definitions/indicator"
                        },
                        "transaction_volume": {
                          "$ref": "#/definitions/indicator"
                        },
                        "group_score": {
                          "type": "number"
                        }
                      },
                      "required": [
                        "transaction_amount",
                        "transaction_frequency",
                        "transaction_volume",
                        "group_score"
                      ]
                    },
                    "actor_risk": {
                      "type": "object",
                      "properties": {
                        "history_with_csb": {
                          "$ref": "#/definitions/indicator"
                        },
                        "jurisdiction": {
                          "$ref": "#/definitions/indicator"
                        },
                        "bank_jurisdiction": {
                          "$ref": "#/definitions/indicator"
                        },
                        "group_score": {
                          "type": "number"
                        }
                      },
                      "required": [
                        "history_with_csb",
                        "jurisdiction",
                        "bank_jurisdiction",
                        "group_score"
                      ]
                    },
                    "client_risk_profile": {
                      "type": "object",
                      "properties": {
                        "client_risk_rating": {
                          "$ref": "#/definitions/indicator"
                        },
                        "amount_vs_declared": {
                          "$ref": "#/definitions/indicator"
                        },
                        "frequency_vs_declared": {
                          "$ref": "#/definitions/indicator"
                        },
                        "velocity": {
                          "$ref": "#/definitions/indicator"
                        },
                        "group_score": {
                          "type": "number"
                        }
                      },
                      "required": [
                        "client_risk_rating",
                        "amount_vs_declared",
                        "frequency_vs_declared",
                        "velocity",
                        "group_score"
                      ]
                    },
                    "purpose_and_documentation": {
                      "type": "object",
                      "properties": {
                        "has_transaction_purpose": {
                          "$ref": "#/definitions/indicator"
                        },
                        "has_supporting_documents": {
                          "$ref": "#/definitions/indicator"
                        },
                        "transaction_amount": {
                          "$ref": "#/definitions/indicator"
                        },
                        "group_score": {
                          "type": "number"
                        }
                      },
                      "required": [
                        "has_transaction_purpose",
                        "has_supporting_documents",
                        "transaction_amount",
                        "group_score"
                      ]
                    },
                    "actor_screening": {
                      "type": "object",
                      "properties": {
                        "actor_sanctioned": {
                          "$ref": "#/definitions/indicator"
                        },
                        "actor_pep": {
                          "$ref": "#/definitions/indicator"
                        },
                        "actor_adverse_media": {
                          "$ref": "#/definitions/indicator"
                        },
                        "actors_bank_sanctioned": {
                          "$ref": "#/definitions/indicator"
                        },
                        "group_score": {
                          "type": "number"
                        }
                      },
                      "required": [
                        "actor_sanctioned",
                        "actor_pep",
                        "actor_adverse_media",
                        "actors_bank_sanctioned",
                        "group_score"
                      ]
                    }
                  },
                  "required": [
                    "risk_level",
                    "risk_score",
                    "manual_review_needed",
                    "actions",
                    "transaction_amount_and_frequency",
                    "actor_risk",
                    "client_risk_profile",
                    "purpose_and_documentation",
                    "actor_screening"
                  ]
                }
              },
              "required": [
                "initial_input",
                "calculation_log"
              ]
            }
          ]
        }
      },
      "required": [
        "@kind",
        "@data"
      ]
    }
  ],
  "unevaluatedProperties": false
}

`

	compiler := jsonschema.NewCompiler()

	//	compiler.RegisterLoader("new", func(urlString string) (result io.ReadCloser, err error) {
	//		refObject := `
	//    {
	//      "@id": "new://INDIVIDUAL_ENTITY__SELF",
	//      "@schema": "json-ir://local@madesst/form/finance/INDIVIDUAL_ENTITY__SELF?12321312123123",
	//      "@archetype": "form",
	//      "@fields": {
	//        "name": "T"
	//      }
	//    }
	//`
	//		refSchemaBytes := []byte(refObject)
	//		return io.NopCloser(bytes.NewReader(refSchemaBytes)), nil
	//	})

	compiler.RegisterLoader("json-ir", func(urlString string) (result io.ReadCloser, err error) {
		refSchema := `
{
  "type": "object",
  "properties": {
    "@id": {
      "type": "string"
    },
    "@schema": {
      "type": "string"
    },
    "@archetype": {
      "allOf": [
        { "type": "string" },
        { "const": "form" }
      ]
    },
    "@meta": {
      "type": "object",
      "additionalProperties": true
    }
  },
  "required": [
    "@id",
    "@schema",
    "@archetype",
    "@kind",
    "@meta"
  ]
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
  "@archetype": "record",
  "@data": {
    "calculation_log": {
      "actions": [
        "Medium-Risk Transactions require approval by the Account Management.",
        "High or Very high risk country. Risk level overridden to Medium and payment route to be confirmed before processing"
      ],
      "actor_risk": {
        "bank_jurisdiction": {
          "calculation": {
            "base_score": 10,
            "factor": 0.35,
            "group_weight": 0.3,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.35 (Factor) * 0.30 (Group Weight) = 1.05 (Final Indicator Score)",
            "total_score": 1.0499999999999998
          },
          "description": "Bank Jurisdiction",
          "risk_level": "LOW",
          "text": "Recipient's Bank Jurisdiction: United Kingdom (code: GB, band: LOW)"
        },
        "group_score": 4.5,
        "history_with_csb": {
          "calculation": {
            "base_score": 30,
            "factor": 0.25,
            "group_weight": 0.3,
            "text": "Calculation: 30.00 (Base Indicator Score) * 0.25 (Factor) * 0.30 (Group Weight) = 2.25 (Final Indicator Score)",
            "total_score": 2.25
          },
          "description": "History with CSB",
          "risk_level": "HIGH",
          "text": "No history with CSB"
        },
        "jurisdiction": {
          "calculation": {
            "base_score": 10,
            "factor": 0.4,
            "group_weight": 0.3,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.40 (Factor) * 0.30 (Group Weight) = 1.20 (Final Indicator Score)",
            "total_score": 1.2000000000000002
          },
          "description": "Jurisdiction",
          "risk_level": "LOW",
          "text": "Recipient's Jurisdiction: Netherlands (code: NL, band: LOW)"
        }
      },
      "actor_screening": {
        "actor_adverse_media": {
          "calculation": {
            "base_score": 10,
            "factor": 0.25,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.25 (Factor) * 0.20 (Group Weight) = 0.50 (Final Indicator Score)",
            "total_score": 0.5
          },
          "description": "Actor Adverse Media",
          "risk_level": "LOW",
          "text": "Actor Adverse Media: No Adverse Media Hits detected for the Recipient"
        },
        "actor_pep": {
          "calculation": {
            "base_score": 10,
            "factor": 0.25,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.25 (Factor) * 0.20 (Group Weight) = 0.50 (Final Indicator Score)",
            "total_score": 0.5
          },
          "description": "Actor PEP",
          "risk_level": "LOW",
          "text": "Actor PEP: No PEP Hits detected for the Recipient"
        },
        "actor_sanctioned": {
          "calculation": {
            "base_score": 10,
            "factor": 0.25,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.25 (Factor) * 0.20 (Group Weight) = 0.50 (Final Indicator Score)",
            "total_score": 0.5
          },
          "description": "Actor Sanctioned",
          "risk_level": "LOW",
          "text": "Actor Sanctioned: No Sanction Hits detected for the Recipient"
        },
        "actors_bank_sanctioned": {
          "calculation": {
            "base_score": 10,
            "factor": 0.25,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.25 (Factor) * 0.20 (Group Weight) = 0.50 (Final Indicator Score)",
            "total_score": 0.5
          },
          "description": "Actors Bank Sanctioned",
          "risk_level": "LOW",
          "text": "Actors Bank Sanctioned: No Sanction Hits detected for the Recipient"
        },
        "group_score": 2
      },
      "client_risk_profile": {
        "amount_vs_declared": {
          "calculation": {
            "base_score": 10,
            "factor": 0.25,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.25 (Factor) * 0.20 (Group Weight) = 0.50 (Final Indicator Score)",
            "total_score": 0.5
          },
          "description": "Amount vs Declared",
          "risk_level": "LOW",
          "text": "Amount vs Declared: Amount is less than 1.5 times the maximum declared amount"
        },
        "client_risk_rating": {
          "calculation": {
            "base_score": 30,
            "factor": 0.4,
            "group_weight": 0.2,
            "text": "Calculation: 30.00 (Base Indicator Score) * 0.40 (Factor) * 0.20 (Group Weight) = 2.40 (Final Indicator Score)",
            "total_score": 2.4000000000000004
          },
          "description": "Client Risk Rating",
          "risk_level": "HIGH",
          "text": "Client Risk Rating: undefined"
        },
        "frequency_vs_declared": {
          "calculation": {
            "base_score": 10,
            "factor": 0.2,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.20 (Factor) * 0.20 (Group Weight) = 0.40 (Final Indicator Score)",
            "total_score": 0.4
          },
          "description": "Frequency vs Declared",
          "risk_level": "LOW",
          "text": "Frequency vs Declared: Frequency is less than 1.5 times the maximum declared frequency"
        },
        "group_score": 4.2,
        "velocity": {
          "calculation": {
            "base_score": 30,
            "factor": 0.15,
            "group_weight": 0.2,
            "text": "Calculation: 30.00 (Base Indicator Score) * 0.15 (Factor) * 0.20 (Group Weight) = 0.90 (Final Indicator Score)",
            "total_score": 0.8999999999999999
          },
          "description": "Velocity",
          "risk_level": "HIGH",
          "text": "Velocity: Sends out >50% of incoming third-party funds within 24h"
        }
      },
      "manual_review_needed": true,
      "purpose_and_documentation": {
        "group_score": 2,
        "has_supporting_documents": {
          "calculation": {
            "base_score": 10,
            "factor": 0.5,
            "group_weight": 0.1,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.50 (Factor) * 0.10 (Group Weight) = 0.50 (Final Indicator Score)",
            "total_score": 0.5
          },
          "description": "Supporting Documents",
          "risk_level": "LOW",
          "text": "Supporting Documents: Supporting documents provided"
        },
        "has_transaction_purpose": {
          "calculation": {
            "base_score": 30,
            "factor": 0.5,
            "group_weight": 0.1,
            "text": "Calculation: 30.00 (Base Indicator Score) * 0.50 (Factor) * 0.10 (Group Weight) = 1.50 (Final Indicator Score)",
            "total_score": 1.5
          },
          "description": "Transaction Purpose",
          "risk_level": "HIGH",
          "text": "Transaction Purpose: No transaction purpose provided"
        },
        "transaction_amount": {
          "calculation": {
            "base_score": 10,
            "factor": 0,
            "group_weight": 0.1,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.00 (Factor) * 0.10 (Group Weight) = 0.00 (Final Indicator Score)",
            "total_score": 0
          },
          "description": "Transaction Amount",
          "risk_level": "LOW",
          "text": "Transaction Amount: Not applicable for outward transactions"
        }
      },
      "risk_level": "MEDIUM",
      "risk_score": 14.7,
      "transaction_amount_and_frequency": {
        "group_score": 2,
        "transaction_amount": {
          "calculation": {
            "base_score": 10,
            "factor": 0.4,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.40 (Factor) * 0.20 (Group Weight) = 0.80 (Final Indicator Score)",
            "total_score": 0.8
          },
          "description": "Transaction amount",
          "risk_level": "LOW",
          "text": "Transaction amount: is less than $10,000"
        },
        "transaction_frequency": {
          "calculation": {
            "base_score": 10,
            "factor": 0.35,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.35 (Factor) * 0.20 (Group Weight) = 0.70 (Final Indicator Score)",
            "total_score": 0.7
          },
          "description": "Transaction frequency",
          "risk_level": "LOW",
          "text": "Transaction frequency: less than 5 payments for the last 30 days"
        },
        "transaction_volume": {
          "calculation": {
            "base_score": 10,
            "factor": 0.25,
            "group_weight": 0.2,
            "text": "Calculation: 10.00 (Base Indicator Score) * 0.25 (Factor) * 0.20 (Group Weight) = 0.50 (Final Indicator Score)",
            "total_score": 0.5
          },
          "description": "Transaction volume",
          "risk_level": "LOW",
          "text": "Transaction volume: is less than $50,000 for the last 30 days"
        }
      }
    },
    "initial_input": {
      "declared": {
        "avg_frequency": "UNDER_6",
        "avg_volume": "10K_TO_100K"
      },
      "has_supporting_docs": true,
      "payment_purpose": "",
      "payment_type": "outward",
      "recipient_bank_jurisdiction": "GB",
      "recipient_jurisdiction": "NL",
      "statistics": {
        "avg_frequency": 0,
        "avg_volume": 0
      },
      "transaction_amount": 9999.99
    }
  },
  "@id": "tf://8a72d881-a83a-46ee-9e13-75baf05c2966@madesst/record/transaction_assessment_result/d4c1e8b0-e4fb-4a33-b6cf-a60f1d28fc4d",
  "@kind": "transaction_assessment_result",
  "@meta": {
    "@author": {
      "@id": "adasadad",
      "@kind": "SYSTEM"
    },
    "@createdAt": "2025-05-23T06:17:20Z",
    "@etag": ":@etag_s3:",
    "@isDataValid": true,
    "@updatedAt": "2025-05-23T06:17:20Z",
    "@version": ":@version_s3:"
  },
  "@schema": "json-ir://local@madesst/types/record/transaction_assessment_result?79cc49e0-fa44-431c-9594-7e34aa13eee0?79cc49e0-fa44-431c-9594-7e34aa13eee0"
}
`), &objectMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal object: %v", err)
	}

	validationResult := schema.Validate(objectMap)
	validationResultList := validationResult.ToList(true)
	log.Printf("Validation result: %v", validationResult.IsValid())
	log.Printf("Validation result details: %v", validationResultList)
}
