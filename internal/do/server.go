package do

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/aybabtme/log"
	"github.com/julienschmidt/httprouter"
)

// Stub returns the URL to a stubbed out DigitalOcean API server.
func Stub() (u *url.URL, done func()) {
	api := &api{
		droplets: &dropletsAPI{db: newDatabase()},
	}

	router := httprouter.New()

	router.POST("/v2/droplets", api.droplets.create)
	router.DELETE("/v2/droplets/:id", api.droplets.delete)

	router.NotFound = http.HandlerFunc(api.NotFound)
	router.PanicHandler = api.panicHandler

	srv := httptest.NewServer(router)
	u, _ = url.Parse(srv.URL)
	return u, srv.Close
}

type api struct {
	droplets *dropletsAPI
}

type apiHandler func() (string, httprouter.Handle)

func (api *api) NotFound(w http.ResponseWriter, r *http.Request) {
	log.KV("req.body", rstring(r.Body)).
		KV("req.path", r.URL.Path).
		Info("unexpected request")
	w.WriteHeader(http.StatusNotFound)
}

func (api *api) panicHandler(w http.ResponseWriter, r *http.Request, pv interface{}) {
	switch v := pv.(type) {
	case *thrown4xx:
		w.WriteHeader(v.code)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": v.data,
		})
		log.KV("err", v.data).KV("http.code", v.code).Info("client error")
	case *thrown5xx:
		w.WriteHeader(v.code)
		log.KV("err", v.err).KV("http.code", v.code).Error("server error")
	default:
		defer func() { w.WriteHeader(http.StatusServiceUnavailable) }()
		log.KV("panic", v).Error("server panic")
		panic(v)
	}
}

func decodeReq(r *http.Request, v interface{}) {
	buf := bytes.NewBuffer(nil)
	_, err := io.CopyN(buf, r.Body, 100<<10)
	if err != io.EOF {
		throwBadRequest("can't read request, %v", err)
	}

	err = json.Unmarshal(buf.Bytes(), v)
	switch e := err.(type) {
	case nil:
	default:
		throwBadRequest("can't parse JSON request, %v: %s", e, buf.String())
	}
}

func respond(w http.ResponseWriter, code int, v interface{}) {
	w.WriteHeader(code)
	if data, err := json.Marshal(v); err != nil {
		throwInternalError("can't generate response: %v", err)
	} else {
		if _, err := w.Write(data); err != nil {
			log.Err(err).Info("couldn't answer client")
		}
	}
}

type thrown4xx struct {
	code int
	data string
}

func throwBadRequest(format string, args ...interface{}) {
	panic(&thrown4xx{
		code: http.StatusBadRequest,
		data: fmt.Sprintf(format, args...),
	})
}

type thrown5xx struct {
	code int
	err  error
}

func throwInternalError(format string, args ...interface{}) {
	panic(&thrown5xx{
		code: http.StatusBadRequest,
		err:  fmt.Errorf(format, args...),
	})
}

func rstring(r io.Reader) string {
	v, _ := ioutil.ReadAll(io.LimitReader(r, 10*(1<<10)))
	return string(v)
}
