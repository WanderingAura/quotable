package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/WanderingAura/quotable/internal/data"
	"github.com/WanderingAura/quotable/internal/validator"
)

var quoteSortSafeList = []string{
	"id",
	"content",
	"modified_at",
	"created_at",
	"user_id",
	"-id",
	"-content",
	"-modified_at",
	"-created_at",
	"-user_id",
}

func (app *application) createQuoteHandler(w http.ResponseWriter, r *http.Request) {

	// an anonymous struct to hold the information that we expect to be in the request body
	// note any key value pairs which do not match one of the struct fields will be silently
	// ignored
	var input struct {
		Content string      `json:"content"`
		Author  string      `json:"author,omitempty"`
		Source  data.Source `json:"source,omitempty"`
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

func (app *application) getQuoteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamByName(r, "quote_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	quote, err := app.models.Quotes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	likeCount, err := app.models.Like.GetLikeDislikeNumForQuote(quote.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"quote": quote, "like_count": likeCount}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

type quoteSearchFields struct {
	Content string
	Tags    []string
	data.Filters
}

func (app *application) readQuoteSearch(r *http.Request, input *quoteSearchFields, v *validator.Validator) {
	qs := r.URL.Query()
	input.Content = app.readString(qs, "content", "")
	input.Tags = app.readCSV(qs, "tags", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = quoteSortSafeList
}

func (app *application) listQuotesHandler(w http.ResponseWriter, r *http.Request) {

	var input quoteSearchFields
	v := validator.New()
	app.readQuoteSearch(r, &input, v)

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	quotes, metadata, err := app.models.Quotes.GetAll(input.Content, input.Tags, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"quotes": quotes, "metadata": metadata}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listUserQuotesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := app.readParamByName(r, "user_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input quoteSearchFields
	v := validator.New()
	app.readQuoteSearch(r, &input, v)
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	quotes, metadata, err := app.models.Quotes.GetAllForUser(userID, input.Content, input.Tags, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"quotes": quotes, "metadata": metadata}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteQuotesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParamByName(r, "quote_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	quote, err := app.models.Quotes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.ID != quote.UserID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Quotes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, envelope{"message": "quote successfully deleted"}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) LikeQuoteHandler(w http.ResponseWriter, r *http.Request) {
	quoteID, err := app.readParamByName(r, "quote_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	var input struct {
		LikeType string `json:"like_type"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.LikeType == "like" || input.LikeType == "dislike", "like_type", "invalid like type")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	val := data.LikeType(data.LikeValue)

	if input.LikeType == "dislike" {
		val = data.DislikeValue
	}

	like := data.Like{
		QuoteID: quoteID,
		UserID:  user.ID,
		Val:     val,
	}

	err = app.models.Like.LikeOrDislikeQuote(like)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"message": "successful", "like": like}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
