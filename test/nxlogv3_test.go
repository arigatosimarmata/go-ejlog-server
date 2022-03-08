package test

import (
	"ejol/ejlog-server/controller"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func (c *ApplicationClient) TestV3MultilineWincor1_WithSingleRequest(t *testing.T) {
// 	fmt.Println("Here start")
// 	req, _ := http.NewRequest(http.MethodGet, "/", nil)
// 	// w := httptest.NewRecorder()

// 	resp, _ := c.HttpClient.Do(req)
// 	if resp.StatusCode == http.StatusUnauthorized {
// 		fmt.Println("Error ")
// 	}
// 	t.Log(resp.StatusCode)
// 	t.Log("Done")
// 	// V3MultilineWincor_1(w, req)
// 	// Index(w, req)
// 	// res := w.Result()
// 	// defer res.Body.Close()
// 	// data, err := ioutil.ReadAll(res.Body)
// 	// if err != nil {
// 	// 	t.Errorf("expected error to be nil got %v", err)
// 	// }

// 	// if string(data) != "ABC" {
// 	// 	t.Error("expected ABC got %v", string(data))
// 	// }
// }

func HealthCheckHandler_WithSingleRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.HealthCheckHandler)

	// handler.HandleFunc("/health-check", HealthCheckHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"alives": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
