package option

type OptionLoader func() (NixosOptionSource, error)

type Scope struct {
	Name        string
	Description string
	Loader      OptionLoader
	Evaluator   EvaluatorFunc
}
