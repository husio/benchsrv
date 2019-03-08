package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func uploadHandler(store Store, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			httpFailf(w, http.StatusBadRequest, "cannot read body: %s", err)
			return
		}

		if sig := w.Header().Get("signature"); !signed(sig, content, secret) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(content) < 10 {
			// Ignore dummy content.
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		benchID, err := store.CreateBenchmark(r.Context(), string(content))
		if err != nil {
			httpFailf(w, http.StatusInternalServerError, "cannot upload: %s", err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, benchID)
	}
}

func signed(sig string, content []byte, secret string) bool {

	// TODO: check the signature of the content to make sure the signer
	// knows the secret.

	return true
}

func listHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		benchmarks, err := store.ListBenchmarks(r.Context(), time.Now(), 100)
		if err != nil {
			httpFailf(w, http.StatusInternalServerError, "cannot list benchmarks: %s", err)
			return
		}
		// TODO render template instead
		for _, b := range benchmarks {
			fmt.Fprintf(w, "%d\t%s\n", b.ID, b.Created)
		}
	}
}

func compareHandler(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		query := r.URL.Query()

		if query.Get("a") == "" || query.Get("b") == "" {
			httpFailf(w, http.StatusBadRequest, "Missing benchmarks IDs. Usage %s?a=<ID>&b=<ID>", r.URL.Path)
			return
		}

		aID, _ := strconv.ParseInt(query.Get("a"), 10, 64)
		a, err := store.FindBenchmark(ctx, aID)
		if err != nil {
			code := http.StatusInternalServerError
			if err == ErrNotFound {
				code = http.StatusNotFound
			}
			httpFailf(w, code, "cannot find benchmark %d: %s", aID, err)
			return
		}

		bID, _ := strconv.ParseInt(query.Get("b"), 10, 64)
		b, err := store.FindBenchmark(ctx, bID)
		if err != nil {
			code := http.StatusInternalServerError
			if err == ErrNotFound {
				code = http.StatusNotFound
			}
			httpFailf(w, code, "cannot find benchmark %d: %s", bID, err)
			return
		}

		cmp, err := Compare(a, b)
		if err != nil {
			httpFailf(w, http.StatusInternalServerError, "cannot compare: %s", err)
			return
		}

		w.Write(cmp)
	}
}

func httpFailf(w http.ResponseWriter, code int, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	w.Header().Set("content-type", "text/html;charset=utf-8")
	w.WriteHeader(code)
	tmpl.ExecuteTemplate(w, "error", msg)
}

var tmpl = template.Must(template.New("").Parse(`

{{define "header"}}
<!doctype html>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
	*            { box-sizing: border-box; }
	html         { position: relative; min-height: 100%; margin: 20px; }
	body         { margin: 40px auto 120px auto; max-width: 50em; line-height: 28px; }
</style>
{{ end}}

{{define "error"}}
{{template "header" .}}
	<div class="error">{{.}}</div>
{{end}}
`))
