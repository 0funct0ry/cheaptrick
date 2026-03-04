package fixture

import "fmt"

func TemplateText() string {
	return `{
  "candidates": [
    {
      "content": {
        "role": "model",
        "parts": [
          {
            "text": "Your textual response here."
          }
        ]
      },
      "finishReason": "STOP"
    }
  ],
  "usageMetadata": {
    "promptTokenCount": 0,
    "candidatesTokenCount": 0,
    "totalTokenCount": 0
  }
}`
}

func TemplateFunctionCall(parsedBody map[string]interface{}) string {
	funcName := "your_function_name"
	args := "{}"

	if tools, ok := parsedBody["tools"].([]interface{}); ok && len(tools) > 0 {
		if tSet, ok := tools[0].(map[string]interface{}); ok {
			if funcs, ok := tSet["functionDeclarations"].([]interface{}); ok && len(funcs) > 0 {
				if fd, ok := funcs[0].(map[string]interface{}); ok {
					if name, ok := fd["name"].(string); ok {
						funcName = name
					}
				}
			}
		}
	}

	return fmt.Sprintf(`{
  "candidates": [
    {
      "content": {
        "role": "model",
        "parts": [
          {
            "functionCall": {
              "name": "%s",
              "args": %s
            }
          }
        ]
      },
      "finishReason": "STOP"
    }
  ]
}`, funcName, args)
}

func Template429() string {
	return `{
  "error": {
    "code": 429,
    "message": "Quota exceeded",
    "status": "RESOURCE_EXHAUSTED"
  }
}`
}

func Template500() string {
	return `{
  "error": {
    "code": 500,
    "message": "Internal error encountered.",
    "status": "INTERNAL"
  }
}`
}
