package templates

import (
	"html/template"
	"io"
    "strconv"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.New("").Funcs(template.FuncMap{
            "WithCommas": WithCommas,
        }).ParseGlob("views/*.html")),
	}
}

func WithCommas(number int) string {
    numberString := strconv.Itoa(number)
    var result string
    for i, digit := range numberString {
        if i != 0 && (len(numberString) - i) % 3 == 0 {
            result += ","
        }
        result += string(digit)
    }
    return result
}
