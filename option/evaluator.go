package option

import "fmt"

type EvaluatorFunc func(optionName string) (string, error)

type AttributeEvaluationError struct {
	Attribute        string
	EvaluationOutput string
}

func (e *AttributeEvaluationError) Error() string {
	return fmt.Sprintf("failed to evaluate attribute %s", e.Attribute)
}
