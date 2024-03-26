package main

import (
	"fmt"
	"net/http"

	"github.com/WanderingAura/quotable/internal/data"
	"github.com/WanderingAura/quotable/internal/validator"
)

func (app *application) createQuoteHandler(w http.ResponseWriter, r *http.Request) {

	// an anonymous struct to hold the information that we expect to be in the request body
	// note any key value pairs which do not match one of the struct fields will be silently
	// ignored
	var input struct {
		Content string      `json:"content"`
		Author  string      `json:"author,omitempty"`
		Source  data.Source `json:"source"`
		Tags    []string    `json:"tags,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	// copying input vals into quote struct prevents user from
	// inputing unwanted quote fields like Version and ID
	quote := data.Quote{
		UserID:  user.ID,
		Content: input.Content,
		Author:  input.Author,
		Source:  input.Source,
		Tags:    input.Tags,
	}
	// initialising validator inside of the handlers gives us
	// flexibility when we have to have multiple validation checks
	v := validator.New()
	data.ValidateQuote(v, &quote)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Quotes.Insert(&quote)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", quote.ID))

	err = app.writeJSON(w, envelope{"quote": quote}, http.StatusOK, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listQuotesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Content string
		Tags    []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Content = app.readString(qs, "content", "")
	input.Tags = app.readCSV(qs, "tags", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"content", "modified", "created", "user_id", "-content", "-modified", "-created", "-user_id"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

}
