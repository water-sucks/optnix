package option

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	glamourStyles "github.com/charmbracelet/glamour/styles"
	"github.com/fatih/color"
)

type ValuePrinterInput struct {
	Value string
	Err   error
}

func (o *NixosOption) PrettyPrint(value *ValuePrinterInput) string {
	var sb strings.Builder

	var (
		titleStyle  = color.New(color.Bold)
		italicStyle = color.New(color.Italic)
	)

	desc := strings.TrimSpace(stripInlineCodeAnnotations(o.Description))
	if desc == "" {
		desc = italicStyle.Sprint("(none)")
	} else {
		d, err := renderer.Render(desc)
		if err != nil {
			desc = italicStyle.Sprintf("warning: failed to render description: %v\n", err) + desc
		} else {
			desc = strings.TrimSpace(d)
		}
	}

	valueText := ""
	if value != nil {
		if value.Err != nil {
			valueText = "failed to evaluate value"

			if e, ok := value.Err.(*AttributeEvaluationError); ok {
				valueText = fmt.Sprintf("%v: %v", valueText, e.EvaluationOutput)
			}

			valueText = color.RedString(valueText)
		} else {
			valueText = color.WhiteString(strings.TrimSpace(value.Value))
		}
	}

	var defaultText string
	if o.Default != nil {
		defaultText = color.WhiteString(strings.TrimSpace(o.Default.Text))
	} else {
		defaultText = italicStyle.Sprint("(none)")
	}

	exampleText := ""
	if o.Example != nil {
		exampleText = color.WhiteString(strings.TrimSpace(o.Example.Text))
	}

	fmt.Fprintf(&sb, "%v\n%v\n\n", titleStyle.Sprint("Name"), o.Name)
	fmt.Fprintf(&sb, "%v\n%v\n\n", titleStyle.Sprint("Description"), desc)
	fmt.Fprintf(&sb, "%v\n%v\n\n", titleStyle.Sprint("Type"), italicStyle.Sprint(o.Type))

	if valueText != "" {
		fmt.Fprintf(&sb, "%v\n%v\n\n", titleStyle.Sprint("Value"), valueText)
	}
	fmt.Fprintf(&sb, "%v\n%v\n\n", titleStyle.Sprint("Default"), defaultText)
	if exampleText != "" {
		fmt.Fprintf(&sb, "%v\n%v\n\n", titleStyle.Sprint("Example"), exampleText)
	}

	if len(o.Declarations) > 0 {
		fmt.Fprintf(&sb, "%v\n", titleStyle.Sprint("Declared In"))
		for _, v := range o.Declarations {
			fmt.Fprintf(&sb, "  - %v\n", italicStyle.Sprint(v))
		}
	}
	if o.ReadOnly {
		fmt.Fprintf(&sb, "\n%v\n", color.YellowString("This option is read-only."))
	}

	return sb.String()
}

var (
	markdownRenderIndentWidth uint = 0
	renderer                       = NewMarkdownRenderer()
)

func NewMarkdownRenderer() *glamour.TermRenderer {
	glamourStyles.DarkStyleConfig.Document.Margin = &markdownRenderIndentWidth

	r, _ := glamour.NewTermRenderer(
		glamour.WithStyles(glamourStyles.DarkStyleConfig),
		glamour.WithWordWrap(80),
	)

	return r
}

var annotationsToRemove = []string{
	"{option}`",
	"{var}`",
	"{file}`",
	"{env}`",
	"{command}`",
	"{manpage}`",
}

func stripInlineCodeAnnotations(slice string) string {
	result := slice

	for _, input := range annotationsToRemove {
		result = strings.ReplaceAll(result, input, "`")
	}

	return result
}
