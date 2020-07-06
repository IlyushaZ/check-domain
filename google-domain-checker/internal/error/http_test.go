package error

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIError(t *testing.T) {
	t.Run("bad request", func(t *testing.T) {
		message := "bad request"
		err := BadRequest(message)
		if err.StatusCode() != http.StatusBadRequest {
			t.Errorf(
				"expected StatusCode() to return %d, %d returned",
				http.StatusBadRequest, err.StatusCode(),
			)
		}

		if err.Error() != message {
			t.Errorf("expected Error() to return %s, %s returned", message, err.Error())
		}
	})

	t.Run("unprocessable entity", func(t *testing.T) {
		message := "unprocessable entity"
		err := UnprocessableEntity(message)

		if err.StatusCode() != http.StatusUnprocessableEntity {
			t.Errorf(
				"expected StatusCode() to return %d, %d returned",
				http.StatusUnprocessableEntity, err.StatusCode(),
			)
		}

		if err.Error() != message {
			t.Errorf("expected Error() to return %s, %s returned", message, err.Error())
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		err := Internal()

		if err.StatusCode() != http.StatusInternalServerError {
			t.Errorf(
				"expected StatusCode() to return %d, %d returned",
				http.StatusInternalServerError, err.StatusCode(),
			)
		}

		if err.Error() != internalErrMessage {
			t.Errorf("expected Error() to return %s, %s returned", internalErrMessage, err.Error())
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		tcs := []struct {
			got, allowed, message string
		}{
			{got: "GET", allowed: "POST", message: "method GET is not allowed. allowed method: POST"},
			{got: "POST", allowed: "PUT", message: "method POST is not allowed. allowed method: PUT"},
			{got: "PUT", allowed: "DELETE", message: "method PUT is not allowed. allowed method: DELETE"},
		}

		for _, tc := range tcs {
			t.Run("method not allowed", func(t *testing.T) {
				err := MethodNotAllowed(tc.got, tc.allowed)

				if err.StatusCode() != http.StatusMethodNotAllowed {
					t.Errorf(
						"expected StatusCode() to return %d, %d returned",
						http.StatusMethodNotAllowed, err.StatusCode(),
					)
				}

				if err.Error() != tc.message {
					t.Errorf("expected Error() to return %s, %s returned", tc.message, err.Error())
				}
			})
		}
	})
}

type testCase struct {
	handler         Handler
	expectedCode    int
	expectedMessage string
}

func TestHandler_RespondError(t *testing.T) {
	tcs := []testCase{
		{
			handler: func(w http.ResponseWriter, r *http.Request) HttpError {
				w.Write([]byte(`{"message":"OK"}`))
				w.WriteHeader(200)
				return nil
			},
			expectedCode:    200,
			expectedMessage: `{"message":"OK"}`,
		},
		{
			handler: func(w http.ResponseWriter, r *http.Request) HttpError {
				return Internal()
			},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: fmt.Sprintf(`{"message":"%s"}`, internalErrMessage),
		},
		{
			handler: func(w http.ResponseWriter, r *http.Request) HttpError {
				return BadRequest("bad request")
			},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: `{"message":"bad request"}`,
		},
	}

	for _, tc := range tcs {
		t.Run("respond error", func(t *testing.T) {
			handler := http.HandlerFunc(tc.handler.RespondError)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "", nil)
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.expectedCode {
				t.Errorf(
					"expected %d code to be returned, %d returned",
					tc.expectedCode, rec.Code,
				)
			}

			if rec.Body.String() != tc.expectedMessage {
				t.Errorf(
					"expected %s message to be returned, %s returned",
					tc.expectedMessage, rec.Body.String(),
				)
			}
		})
	}
}
