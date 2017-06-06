package nim

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestStackServeHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	result := ""

	ns := New()
	ns.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "_1bef"
		next(w, r)
		result += "_1aft"
	})
	ns.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "_2bef"
		next(w, r)
		result += "_2aft"
	})
	ns.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "_3here"
		w.WriteHeader(http.StatusBadRequest)
	})

	ns.ServeHTTP(rec, (*http.Request)(nil))

	expect(t, result, "_1bef_2bef_3here_2aft_1aft")
	expect(t, rec.Code, http.StatusBadRequest)
}

// Ensures that the middleware chain
// can correctly return all of its handlers.
func TestStackWithHandlerFunc(t *testing.T) {
	rec := httptest.NewRecorder()
	ns := New()
	handles := ns.handlers
	expect(t, 0, len(handles))

	ns.WithHandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.WriteHeader(http.StatusOK)
	})

	// Expects the length of handlers to be exactly 1
	// after adding exactly one handler to the middleware chain
	handles = ns.handlers
	expect(t, 1, len(handles))

	// Ensures that the first handler that is in sequence behaves
	// exactly the same as the one that was registered earlier
	handles[0](rec, (*http.Request)(nil), nil)
	expect(t, rec.Code, http.StatusOK)
}

func TesStackWithNil(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Expected Stack.With(nil) to panic, but it did not")
		}
	}()

	ns := New()
	ns.With(nil)
}
