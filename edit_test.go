package mwclient

import (
	"fmt"
	"net/http"
	"testing"

	"cgt.name/pkg/go-mwclient/params"
)

func TestEdit(t *testing.T) {
	resp := `{"edit":{"result":"Success","pageid":42,"title":"PAGE",
	"contentmodel":"wikitext","oldrevid":7936766,"newrevid":7950155,
	"newtimestamp":"2015-02-12T17:13:01Z"}}`

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic("Bad HTTP form")
		}

		if r.Method != "POST" {
			t.Fatalf("edit requests must be posted. Method: %v", r.Method)
		}
		if v := r.Form.Get("action"); v != "edit" {
			t.Fatalf("action != edit: action=%s", v)
		}
		if v := r.Form.Get("token"); v != "VALIDTOKEN" {
			t.Fatalf("token != VALIDTOKEN: token=%s", v)
		}

		fmt.Fprint(w, resp)
	}

	server, client := setup(httpHandler)
	defer server.Close()

	client.Tokens[CSRFToken] = "VALIDTOKEN"
	err := client.Edit(params.Values{})
	if err != nil {
		t.Fatalf("edit request returned error: %v", err)
	}
}

func TestGetToken(t *testing.T) {
	resp := `{"batchcomplete":"","query":{"tokens":{"csrftoken":"+\\"}}}`
	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			panic("Bad HTTP form")
		}

		if v := r.Form.Get("action"); v != "query" {
			t.Fatalf("action != query: action=%s", v)
		}
		if v := r.Form.Get("meta"); v != "tokens" {
			t.Fatalf("meta != tokens: meta=%s", v)
		}
		if v := r.Form.Get("type"); v != CSRFToken {
			t.Fatalf("meta != %s: meta=%s", CSRFToken, v)
		}

		fmt.Fprint(w, resp)
	}

	server, client := setup(httpHandler)
	defer server.Close()

	token, err := client.GetToken(CSRFToken)
	if err != nil {
		t.Fatalf("token request failed: %v", err)
	}
	if token != "+\\" {
		t.Fatalf("received token does not match sent token")
	}
}

func TestGetCachedToken(t *testing.T) {
	client, err := New("http://example.com", "go-mwclient test")
	if err != nil {
		panic(err)
	}
	client.Tokens[CSRFToken] = "tokenvalue"
	gotToken, err := client.GetToken(CSRFToken)
	if err != nil {
		panic(err)
	}
	if gotToken != client.Tokens[CSRFToken] {
		t.Fatalf("got token does not match manually cached token: CSRFToken=%s",
			gotToken)
	}
}
