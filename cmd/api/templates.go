package main

type templateData struct {
	ErrorTitle   string
	ErrorContent string
}

func (app *application) newTemplateData() templateData {
	return templateData{}
}
