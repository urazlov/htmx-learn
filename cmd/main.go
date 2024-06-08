package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*html")),
	}
}

type Contact struct {
    Name string
    Email string
}

type Contacts = []Contact

type Data struct {
    Contacts Contacts

}

type FormData struct {
    Values map[string]string
    Errors map[string]string
}

type Page struct {
    Data Data
    Form FormData
}

func newPage() Page {
    return Page{
        Data: newData(),
        Form: newFormData(),
    }
}

func newFormData() FormData {
    return FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

func newContact(name, email string) Contact {
    return Contact{
        Name: name,
        Email: email,
    }
}

func newData() Data {
    return Data{
        Contacts: Contacts{
            newContact("John", "jd@gmail.com"),
            newContact("Clara", "cd@gmail.com"),
        },
    }
}

func (d *Data) hasEmail(email string) bool {
    for _, contact := range d.Contacts {
        if contact.Email == email {
            return true
        }
    }

    return false
}

func main() {
	e := echo.New()
    e.Use(middleware.Logger())

    page := newPage()
    e.Renderer = NewTemplates()

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", page)
    })

    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")

        if page.Data.hasEmail(email) {
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email
            formData.Errors["email"] = "Email already exist"
            return c.Render(422, "form", formData)
        }

        contact := newContact(name, email);
        page.Data.Contacts = append(page.Data.Contacts, newContact(name, email))
        c.Render(200, "form", newFormData())
        return c.Render(200, "oob-contact", contact)
    })

    e.Logger.Fatal(e.Start(":8080"))
}
