package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/distrib-jobs/internal/httpapi"
	"github.com/nats-io/nats.go"
)

type fakeJS struct{}

func (f *fakeJS) Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	return &nats.PubAck{Stream: "JOBS", Seq: 1}, nil
}

func (f *fakeJS) AccountInfo(opts ...nats.JSOpt) (*nats.AccountInfo, error) { return nil, nil }

func TestEnqueue(t *testing.T) {
	h := httpapi.NewHandler(&fakeAdapter{js: &fakeJS{}}, "JOBS")
	w := httptest.NewRecorder()
	body := map[string]any{"type":"email","payload":map[string]any{"to":"a@b.com"}}
	b,_ := json.Marshal(body)
	r := httptest.NewRequest("POST", "/jobs", bytes.NewReader(b))
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatal("unexpected status")
	}
}

type fakeAdapter struct{ js nats.JetStreamContext }

func (f *fakeAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
