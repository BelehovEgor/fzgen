package main

// Edit if desired. Code generated by "fzgen net/http".

import (
	"bufio"
	"context"
	"io"
	http "net/http"
	"net/url"
	"testing"
	"time"

	"github.com/thepudds/fzgen/fuzzer"
)

func fabric_interface_httpHandler() http.Handler {
	return &http.ServeMux{}
}

func fabric_func_43_1() func(http.ResponseWriter, *http.Request) {
	return http.NotFound
}

func fabric_interface_httpRoundTripper() http.RoundTripper {
	return &http.Transport{}
}

func fabric_func_79_1() func(http.ResponseWriter, *http.Request) {
	return http.NotFound
}

func Fuzz_Client_CloseIdleConnections(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c)
		if c == nil {
			return
		}

		c.CloseIdleConnections()
	})
}

func Fuzz_Client_Do(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var req *http.Request
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c, &req)
		if c == nil || req == nil {
			return
		}

		c.Do(req)
	})
}

func Fuzz_Client_Get(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c, &url)
		if c == nil {
			return
		}

		c.Get(url)
	})
}

func Fuzz_Client_Head(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c, &url)
		if c == nil {
			return
		}

		c.Head(url)
	})
}

func Fuzz_Client_Post(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url string
		var contentType string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c, &url, &contentType, &body)
		if c == nil {
			return
		}

		c.Post(url, contentType, body)
	})
}

func Fuzz_Client_PostForm(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url_ string
		var d3 url.Values
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c, &url_, &d3)
		if c == nil {
			return
		}

		c.PostForm(url_, d3)
	})
}

func Fuzz_Cookie_String(f *testing.F) {
	f.Fuzz(func(t *testing.T, line string) {
		c, err := http.ParseSetCookie(line)
		if err != nil {
			return
		}
		c.String()
	})
}

func Fuzz_Cookie_Valid(f *testing.F) {
	f.Fuzz(func(t *testing.T, line string) {
		c, err := http.ParseSetCookie(line)
		if err != nil {
			return
		}
		c.Valid()
	})
}

func Fuzz_MaxBytesError_Error(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var e *http.MaxBytesError
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&e)
		if e == nil {
			return
		}

		e.Error()
	})
}

func Fuzz_ProtocolError_Error(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pe *http.ProtocolError
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&pe)
		if pe == nil {
			return
		}

		pe.Error()
	})
}

// skipping Fuzz_ProtocolError_Is because parameters include func, chan, or unsupported interface: error

func Fuzz_Request_AddCookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var c *http.Cookie
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &c)
		if c == nil {
			return
		}

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.AddCookie(c)
	})
}

func Fuzz_Request_BasicAuth(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.BasicAuth()
	})
}

func Fuzz_Request_Clone(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var ctx context.Context
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &ctx)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.Clone(ctx)
	})
}

func Fuzz_Request_Context(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.Context()
	})
}

func Fuzz_Request_Cookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var name string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &name)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.Cookie(name)
	})
}

func Fuzz_Request_Cookies(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.Cookies()
	})
}

func Fuzz_Request_CookiesNamed(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var name string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &name)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.CookiesNamed(name)
	})
}

func Fuzz_Request_FormFile(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var key string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &key)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.FormFile(key)
	})
}

func Fuzz_Request_FormValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var key string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &key)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.FormValue(key)
	})
}

func Fuzz_Request_MultipartReader(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.MultipartReader()
	})
}

func Fuzz_Request_ParseForm(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.ParseForm()
	})
}

func Fuzz_Request_ParseMultipartForm(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var maxMemory int64
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &maxMemory)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.ParseMultipartForm(maxMemory)
	})
}

func Fuzz_Request_PathValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var name string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &name)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.PathValue(name)
	})
}

func Fuzz_Request_PostFormValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var key string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &key)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.PostFormValue(key)
	})
}

func Fuzz_Request_ProtoAtLeast(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var major int
		var minor int
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &major, &minor)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.ProtoAtLeast(major, minor)
	})
}

func Fuzz_Request_Referer(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.Referer()
	})
}

func Fuzz_Request_SetBasicAuth(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var username string
		var password string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &username, &password)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.SetBasicAuth(username, password)
	})
}

func Fuzz_Request_SetPathValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var name string
		var value string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &name, &value)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.SetPathValue(name, value)
	})
}

func Fuzz_Request_UserAgent(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.UserAgent()
	})
}

func Fuzz_Request_WithContext(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var ctx context.Context
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &ctx)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.WithContext(ctx)
	})
}

func Fuzz_Request_Write(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var w io.Writer
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &w)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.Write(w)
	})
}

func Fuzz_Request_WriteProxy(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		var w io.Writer
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body, &w)

		r, err := http.NewRequest(method, url, body)
		if err != nil {
			return
		}
		r.WriteProxy(w)
	})
}

func Fuzz_Response_Cookies(f *testing.F) {
	f.Fuzz(func(t *testing.T, url string) {
		r, err := http.Get(url)
		if err != nil {
			return
		}
		r.Cookies()
	})
}

func Fuzz_Response_Location(f *testing.F) {
	f.Fuzz(func(t *testing.T, url string) {
		r, err := http.Get(url)
		if err != nil {
			return
		}
		r.Location()
	})
}

func Fuzz_Response_ProtoAtLeast(f *testing.F) {
	f.Fuzz(func(t *testing.T, url string, major int, minor int) {
		r, err := http.Get(url)
		if err != nil {
			return
		}
		r.ProtoAtLeast(major, minor)
	})
}

func Fuzz_Response_Write(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url string
		var w io.Writer
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&url, &w)

		r, err := http.Get(url)
		if err != nil {
			return
		}
		r.Write(w)
	})
}

func Fuzz_ResponseController_EnableFullDuplex(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.ResponseController
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c)
		if c == nil {
			return
		}

		c.EnableFullDuplex()
	})
}

func Fuzz_ResponseController_Flush(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.ResponseController
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c)
		if c == nil {
			return
		}

		c.Flush()
	})
}

func Fuzz_ResponseController_Hijack(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.ResponseController
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c)
		if c == nil {
			return
		}

		c.Hijack()
	})
}

func Fuzz_ResponseController_SetReadDeadline(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.ResponseController
		var deadline time.Time
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c, &deadline)
		if c == nil {
			return
		}

		c.SetReadDeadline(deadline)
	})
}

func Fuzz_ResponseController_SetWriteDeadline(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.ResponseController
		var deadline time.Time
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c, &deadline)
		if c == nil {
			return
		}

		c.SetWriteDeadline(deadline)
	})
}

func Fuzz_ServeMux_Handle(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pattern string
		var handler http.Handler
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&pattern, &handler)

		mux := http.NewServeMux()
		mux.Handle(pattern, handler)
	})
}

func Fuzz_ServeMux_HandleFunc(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pattern string
		var handler func(http.ResponseWriter, *http.Request)
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&pattern, &handler)

		mux := http.NewServeMux()
		mux.HandleFunc(pattern, handler)
	})
}

func Fuzz_ServeMux_Handler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&r)
		if r == nil {
			return
		}

		mux := http.NewServeMux()
		mux.Handler(r)
	})
}

// skipping Fuzz_ServeMux_ServeHTTP because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

func Fuzz_Server_Close(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&srv)
		if srv == nil {
			return
		}

		srv.Close()
	})
}

func Fuzz_Server_ListenAndServe(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&srv)
		if srv == nil {
			return
		}

		srv.ListenAndServe()
	})
}

func Fuzz_Server_ListenAndServeTLS(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		var certFile string
		var keyFile string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&srv, &certFile, &keyFile)
		if srv == nil {
			return
		}

		srv.ListenAndServeTLS(certFile, keyFile)
	})
}

// skipping Fuzz_Server_RegisterOnShutdown because parameters include func, chan, or unsupported interface: func()

// skipping Fuzz_Server_Serve because parameters include func, chan, or unsupported interface: net.Listener

// skipping Fuzz_Server_ServeTLS because parameters include func, chan, or unsupported interface: net.Listener

func Fuzz_Server_SetKeepAlivesEnabled(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		var v bool
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&srv, &v)
		if srv == nil {
			return
		}

		srv.SetKeepAlivesEnabled(v)
	})
}

func Fuzz_Server_Shutdown(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		var ctx context.Context
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&srv, &ctx)
		if srv == nil {
			return
		}

		srv.Shutdown(ctx)
	})
}

func Fuzz_Transport_CancelRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t1 *http.Transport
		var req *http.Request
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&t1, &req)
		if t1 == nil || req == nil {
			return
		}

		t1.CancelRequest(req)
	})
}

func Fuzz_Transport_Clone(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t1 *http.Transport
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&t1)
		if t1 == nil {
			return
		}

		t1.Clone()
	})
}

func Fuzz_Transport_CloseIdleConnections(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t1 *http.Transport
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&t1)
		if t1 == nil {
			return
		}

		t1.CloseIdleConnections()
	})
}

func Fuzz_Transport_RegisterProtocol(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t1 *http.Transport
		var scheme string
		var rt http.RoundTripper
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&t1, &scheme, &rt)
		if t1 == nil {
			return
		}

		t1.RegisterProtocol(scheme, rt)
	})
}

func Fuzz_Transport_RoundTrip(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t1 *http.Transport
		var req *http.Request
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&t1, &req)
		if t1 == nil || req == nil {
			return
		}

		t1.RoundTrip(req)
	})
}

func Fuzz_ConnState_String(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c http.ConnState
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&c)

		c.String()
	})
}

func Fuzz_Dir_Open(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var d http.Dir
		var name string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&d, &name)

		d.Open(name)
	})
}

// skipping Fuzz_HandlerFunc_ServeHTTP because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

func Fuzz_Header_Add(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		var value string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &key, &value)

		h.Add(key, value)
	})
}

func Fuzz_Header_Clone(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h)

		h.Clone()
	})
}

func Fuzz_Header_Del(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &key)

		h.Del(key)
	})
}

func Fuzz_Header_Get(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &key)

		h.Get(key)
	})
}

func Fuzz_Header_Set(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		var value string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &key, &value)

		h.Set(key, value)
	})
}

func Fuzz_Header_Values(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &key)

		h.Values(key)
	})
}

func Fuzz_Header_Write(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var w io.Writer
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &w)

		h.Write(w)
	})
}

func Fuzz_Header_WriteSubset(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var w io.Writer
		var exclude map[string]bool
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &w, &exclude)

		h.WriteSubset(w, exclude)
	})
}

func Fuzz_AllowQuerySemicolons(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Handler
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h)

		http.AllowQuerySemicolons(h)
	})
}

func Fuzz_CanonicalHeaderKey(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		http.CanonicalHeaderKey(s)
	})
}

func Fuzz_DetectContentType(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var d1 []byte
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&d1)

		http.DetectContentType(d1)
	})
}

// skipping Fuzz_Error because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

// skipping Fuzz_FS because parameters include func, chan, or unsupported interface: io/fs.FS

// skipping Fuzz_FileServer because parameters include func, chan, or unsupported interface: net/http.FileSystem

// skipping Fuzz_FileServerFS because parameters include func, chan, or unsupported interface: io/fs.FS

func Fuzz_Get(f *testing.F) {
	f.Fuzz(func(t *testing.T, url string) {
		http.Get(url)
	})
}

func Fuzz_Handle(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pattern string
		var handler http.Handler
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&pattern, &handler)

		http.Handle(pattern, handler)
	})
}

func Fuzz_HandleFunc(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pattern string
		var handler func(http.ResponseWriter, *http.Request)
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&pattern, &handler)

		http.HandleFunc(pattern, handler)
	})
}

func Fuzz_Head(f *testing.F) {
	f.Fuzz(func(t *testing.T, url string) {
		http.Head(url)
	})
}

func Fuzz_ListenAndServe(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var addr string
		var handler http.Handler
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&addr, &handler)

		http.ListenAndServe(addr, handler)
	})
}

func Fuzz_ListenAndServeTLS(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var addr string
		var certFile string
		var keyFile string
		var handler http.Handler
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&addr, &certFile, &keyFile, &handler)

		http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	})
}

func Fuzz_MaxBytesHandler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Handler
		var n int64
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &n)

		http.MaxBytesHandler(h, n)
	})
}

// skipping Fuzz_MaxBytesReader because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

// skipping Fuzz_NewFileTransport because parameters include func, chan, or unsupported interface: net/http.FileSystem

// skipping Fuzz_NewFileTransportFS because parameters include func, chan, or unsupported interface: io/fs.FS

func Fuzz_NewRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&method, &url, &body)

		http.NewRequest(method, url, body)
	})
}

func Fuzz_NewRequestWithContext(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var ctx context.Context
		var method string
		var url string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&ctx, &method, &url, &body)

		http.NewRequestWithContext(ctx, method, url, body)
	})
}

// skipping Fuzz_NewResponseController because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

// skipping Fuzz_NotFound because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

func Fuzz_ParseCookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, line string) {
		http.ParseCookie(line)
	})
}

func Fuzz_ParseHTTPVersion(f *testing.F) {
	f.Fuzz(func(t *testing.T, vers string) {
		http.ParseHTTPVersion(vers)
	})
}

func Fuzz_ParseSetCookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, line string) {
		http.ParseSetCookie(line)
	})
}

func Fuzz_ParseTime(f *testing.F) {
	f.Fuzz(func(t *testing.T, text string) {
		http.ParseTime(text)
	})
}

func Fuzz_Post(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url string
		var contentType string
		var body io.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&url, &contentType, &body)

		http.Post(url, contentType, body)
	})
}

func Fuzz_PostForm(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url_ string
		var d2 url.Values
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&url_, &d2)

		http.PostForm(url_, d2)
	})
}

func Fuzz_ProxyFromEnvironment(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var req *http.Request
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&req)
		if req == nil {
			return
		}

		http.ProxyFromEnvironment(req)
	})
}

func Fuzz_ProxyURL(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var fixedURL *url.URL
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&fixedURL)
		if fixedURL == nil {
			return
		}

		http.ProxyURL(fixedURL)
	})
}

func Fuzz_ReadRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var b *bufio.Reader
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&b)
		if b == nil {
			return
		}

		http.ReadRequest(b)
	})
}

func Fuzz_ReadResponse(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *bufio.Reader
		var req *http.Request
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&r, &req)
		if r == nil || req == nil {
			return
		}

		http.ReadResponse(r, req)
	})
}

// skipping Fuzz_Redirect because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

func Fuzz_RedirectHandler(f *testing.F) {
	f.Fuzz(func(t *testing.T, url string, code int) {
		http.RedirectHandler(url, code)
	})
}

// skipping Fuzz_Serve because parameters include func, chan, or unsupported interface: net.Listener

// skipping Fuzz_ServeContent because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

// skipping Fuzz_ServeFile because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

// skipping Fuzz_ServeFileFS because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

// skipping Fuzz_ServeTLS because parameters include func, chan, or unsupported interface: net.Listener

// skipping Fuzz_SetCookie because parameters include func, chan, or unsupported interface: net/http.ResponseWriter

func Fuzz_StatusText(f *testing.F) {
	f.Fuzz(func(t *testing.T, code int) {
		http.StatusText(code)
	})
}

func Fuzz_StripPrefix(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var prefix string
		var h http.Handler
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&prefix, &h)

		http.StripPrefix(prefix, h)
	})
}

func Fuzz_TimeoutHandler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Handler
		var dt time.Duration
		var msg string
		fz := fuzzer.NewFuzzer(data)
		fz.Fill(&h, &dt, &msg)

		http.TimeoutHandler(h, dt, msg)
	})
}