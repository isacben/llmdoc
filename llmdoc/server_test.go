package llmdoc

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	indexTemplate, err := template.New("index").Parse(`
	<html>
		<head>
			<title>LLMDoc</title>
		</head>
		<body>
			<h1>LLMDoc</h1>
		</body>
	</html>
	`)

	if err != nil {
		t.Fatal(err)
	}

	handler, err := IndexHandler(indexTemplate)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)
	status := recorder.Code
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `
	<html>
		<head>
			<title>LLMDoc</title>
		</head>
		<body>
			<h1>LLMDoc</h1>
		</body>
	</html>
	`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}
