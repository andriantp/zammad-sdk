package article

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/andriantp/zammad-sdk/internal/transport"
)

func TestService_OpenAttachment(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/ticket_attachment/58/79/72" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Disposition", `inline; filename="image-1.png"; filename*=UTF-8''image-1.png`)
		w.Header().Set("Content-Length", "9")
		_, _ = w.Write([]byte("png-bytes"))
	}))

	file, err := service.OpenAttachment(context.Background(), 58, 79, 72)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer file.Body.Close()

	data, err := io.ReadAll(file.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if string(data) != "png-bytes" {
		t.Errorf("unexpected body %q", data)
	}

	if file.Filename != "image-1.png" || file.ContentType != "image/png" || file.Size != 9 {
		t.Errorf("unexpected metadata: %+v", file)
	}
}

func TestService_OpenAttachment_APIError(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"The file does not belong to the specified article."}`))
	}))

	file, err := service.OpenAttachment(context.Background(), 58, 79, 999)
	if file != nil {
		t.Error("expected nil file on error")
	}

	var apiErr *transport.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 APIError, got %v", err)
	}
}

func TestService_OpenAttachment_InvalidIDs(t *testing.T) {
	service := newTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no HTTP request")
	}))

	if _, err := service.OpenAttachment(context.Background(), 1, 0, 1); err == nil {
		t.Error("expected error for invalid identifiers")
	}
}

// rawOnlyTransport implements Transport but not StreamTransport, like a
// caller's own mock written against v1.
type rawOnlyTransport struct{ data []byte }

func (rawOnlyTransport) Do(context.Context, string, string, any, any) error { return nil }

func (t rawOnlyTransport) DoRaw(context.Context, string, string) ([]byte, error) {
	return t.data, nil
}

func TestService_OpenAttachment_FallbackWithoutStreaming(t *testing.T) {
	service := New(rawOnlyTransport{data: []byte("abc")})

	file, err := service.OpenAttachment(context.Background(), 1, 2, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer file.Body.Close()

	data, _ := io.ReadAll(file.Body)
	if string(data) != "abc" || file.Size != 3 || file.Filename != "" {
		t.Errorf("unexpected fallback result: %q %+v", data, file)
	}
}

func TestFilenameFromHeader(t *testing.T) {
	for _, tc := range []struct{ header, want string }{
		{`inline; filename="image-1.png"`, "image-1.png"},
		{`attachment; filename*=UTF-8''r%C3%A9sum%C3%A9.pdf`, "résumé.pdf"},
		{`attachment; filename="../../.bashrc"`, ".bashrc"},
		{`attachment; filename="/etc/passwd"`, "passwd"},
		{`attachment; filename="..\\..\\evil.txt"`, "evil.txt"},
		{`attachment; filename=".."`, ""},
		{`attachment; filename="..."`, ""},
		{`attachment; filename=""`, ""},
		{`attachment; filename="a` + "\x01" + `b.txt"`, ""},
		{`attachment`, ""},
		{``, ""},
		{`not a valid header;;;`, ""},
	} {
		if got := filenameFromHeader(tc.header); got != tc.want {
			t.Errorf("filenameFromHeader(%q) = %q, want %q", tc.header, got, tc.want)
		}
	}
}
