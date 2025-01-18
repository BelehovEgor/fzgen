package main

// Edit if desired. Code generated by "fzgen net/http".

import (
	"bufio"
	"context"
	"io"
	http "net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/BelehovEgor/fzgen/fuzzer"
)

func Fuzz_Client_Do(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var req *http.Request
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c, &req)
		if err != nil || c == nil || req == nil {
			return
		}

		c.Do(req)
	})
}

func Fuzz_Client_Get(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url_0 string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c, &url_0)
		if err != nil || c == nil {
			return
		}

		c.Get(url_0)
	})
}

func Fuzz_Client_Head(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url_0 string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c, &url_0)
		if err != nil || c == nil {
			return
		}

		c.Head(url_0)
	})
}

func Fuzz_Client_Post(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url_0 string
		var contentType string
		var body io.Reader
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c, &url_0, &contentType, &body)
		if err != nil || c == nil {
			return
		}

		c.Post(url_0, contentType, body)
	})
}

func Fuzz_Client_PostForm(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.Client
		var url_0 string
		var data_0 url.Values
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c, &url_0, &data_0)
		if err != nil || c == nil {
			return
		}

		c.PostForm(url_0, data_0)
	})
}

func Fuzz_ProtocolError_Is(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pe *http.ProtocolError
		var err error
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err_1 := fz.Fill2(&pe, &err)
		if err_1 != nil || pe == nil {
			return
		}

		pe.Is(err)
	})
}

func Fuzz_Request_AddCookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var c *http.Cookie
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &c)
		if err != nil || r == nil || c == nil {
			return
		}

		r.AddCookie(c)
	})
}

func Fuzz_Request_Clone(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var ctx context.Context
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &ctx)
		if err != nil || r == nil {
			return
		}

		r.Clone(ctx)
	})
}

func Fuzz_Request_Cookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var name string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &name)
		if err != nil || r == nil {
			return
		}

		r.Cookie(name)
	})
}

func Fuzz_Request_CookiesNamed(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var name string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &name)
		if err != nil || r == nil {
			return
		}

		r.CookiesNamed(name)
	})
}

func Fuzz_Request_FormFile(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var key string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &key)
		if err != nil || r == nil {
			return
		}

		r.FormFile(key)
	})
}

func Fuzz_Request_FormValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var key string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &key)
		if err != nil || r == nil {
			return
		}

		r.FormValue(key)
	})
}

func Fuzz_Request_ParseMultipartForm(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var maxMemory int64
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &maxMemory)
		if err != nil || r == nil {
			return
		}

		r.ParseMultipartForm(maxMemory)
	})
}

func Fuzz_Request_PathValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var name string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &name)
		if err != nil || r == nil {
			return
		}

		r.PathValue(name)
	})
}

func Fuzz_Request_PostFormValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var key string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &key)
		if err != nil || r == nil {
			return
		}

		r.PostFormValue(key)
	})
}

func Fuzz_Request_ProtoAtLeast(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var major int
		var minor int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &major, &minor)
		if err != nil || r == nil {
			return
		}

		r.ProtoAtLeast(major, minor)
	})
}

func Fuzz_Request_SetBasicAuth(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var username string
		var password string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &username, &password)
		if err != nil || r == nil {
			return
		}

		r.SetBasicAuth(username, password)
	})
}

func Fuzz_Request_SetPathValue(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var name string
		var value string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &name, &value)
		if err != nil || r == nil {
			return
		}

		r.SetPathValue(name, value)
	})
}

func Fuzz_Request_WithContext(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var ctx context.Context
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &ctx)
		if err != nil || r == nil {
			return
		}

		r.WithContext(ctx)
	})
}

func Fuzz_Request_Write(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var w io.Writer
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &w)
		if err != nil || r == nil {
			return
		}

		r.Write(w)
	})
}

func Fuzz_Request_WriteProxy(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Request
		var w io.Writer
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &w)
		if err != nil || r == nil {
			return
		}

		r.WriteProxy(w)
	})
}

func Fuzz_Response_ProtoAtLeast(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Response
		var major int
		var minor int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &major, &minor)
		if err != nil || r == nil {
			return
		}

		r.ProtoAtLeast(major, minor)
	})
}

func Fuzz_Response_Write(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *http.Response
		var w io.Writer
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &w)
		if err != nil || r == nil {
			return
		}

		r.Write(w)
	})
}

func Fuzz_ResponseController_SetReadDeadline(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.ResponseController
		var deadline time.Time
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c, &deadline)
		if err != nil || c == nil {
			return
		}

		c.SetReadDeadline(deadline)
	})
}

func Fuzz_ResponseController_SetWriteDeadline(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var c *http.ResponseController
		var deadline time.Time
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&c, &deadline)
		if err != nil || c == nil {
			return
		}

		c.SetWriteDeadline(deadline)
	})
}

func Fuzz_ServeMux_Handle(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var mux *http.ServeMux
		var pattern string
		var handler http.Handler
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&mux, &pattern, &handler)
		if err != nil || mux == nil {
			return
		}

		mux.Handle(pattern, handler)
	})
}

func Fuzz_ServeMux_HandleFunc(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var mux *http.ServeMux
		var pattern string
		var handler func(http.ResponseWriter, *http.Request)
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&mux, &pattern, &handler)
		if err != nil || mux == nil {
			return
		}

		mux.HandleFunc(pattern, handler)
	})
}

func Fuzz_ServeMux_Handler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var mux *http.ServeMux
		var r *http.Request
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&mux, &r)
		if err != nil || mux == nil || r == nil {
			return
		}

		mux.Handler(r)
	})
}

// skipping Fuzz_ServeMux_ServeHTTP because parameters include unsupported type: net/http.ResponseWriter

func Fuzz_Server_ListenAndServeTLS(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		var certFile string
		var keyFile string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&srv, &certFile, &keyFile)
		if err != nil || srv == nil {
			return
		}

		srv.ListenAndServeTLS(certFile, keyFile)
	})
}

// skipping Fuzz_Server_RegisterOnShutdown because parameters include unsupported type: func()

// skipping Fuzz_Server_Serve because parameters include unsupported type: net.Listener

// skipping Fuzz_Server_ServeTLS because parameters include unsupported type: net.Listener

func Fuzz_Server_SetKeepAlivesEnabled(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		var v bool
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&srv, &v)
		if err != nil || srv == nil {
			return
		}

		srv.SetKeepAlivesEnabled(v)
	})
}

func Fuzz_Server_Shutdown(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var srv *http.Server
		var ctx context.Context
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&srv, &ctx)
		if err != nil || srv == nil {
			return
		}

		srv.Shutdown(ctx)
	})
}

func Fuzz_Transport_CancelRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 *http.Transport
		var req *http.Request
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &req)
		if err != nil || t_0 == nil || req == nil {
			return
		}

		t_0.CancelRequest(req)
	})
}

func Fuzz_Transport_RegisterProtocol(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 *http.Transport
		var scheme string
		var rt http.RoundTripper
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &scheme, &rt)
		if err != nil || t_0 == nil {
			return
		}

		t_0.RegisterProtocol(scheme, rt)
	})
}

func Fuzz_Transport_RoundTrip(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var t_0 *http.Transport
		var req *http.Request
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&t_0, &req)
		if err != nil || t_0 == nil || req == nil {
			return
		}

		t_0.RoundTrip(req)
	})
}

func Fuzz_Dir_Open(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var d http.Dir
		var name string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&d, &name)
		if err != nil {
			return
		}

		d.Open(name)
	})
}

// skipping Fuzz_HandlerFunc_ServeHTTP because parameters include unsupported type: net/http.ResponseWriter

func Fuzz_Header_Add(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		var value string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &key, &value)
		if err != nil {
			return
		}

		h.Add(key, value)
	})
}

func Fuzz_Header_Del(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &key)
		if err != nil {
			return
		}

		h.Del(key)
	})
}

func Fuzz_Header_Get(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &key)
		if err != nil {
			return
		}

		h.Get(key)
	})
}

func Fuzz_Header_Set(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		var value string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &key, &value)
		if err != nil {
			return
		}

		h.Set(key, value)
	})
}

func Fuzz_Header_Values(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var key string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &key)
		if err != nil {
			return
		}

		h.Values(key)
	})
}

func Fuzz_Header_Write(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var w io.Writer
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &w)
		if err != nil {
			return
		}

		h.Write(w)
	})
}

func Fuzz_Header_WriteSubset(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Header
		var w io.Writer
		var exclude map[string]bool
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &w, &exclude)
		if err != nil {
			return
		}

		h.WriteSubset(w, exclude)
	})
}

func Fuzz_AllowQuerySemicolons(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Handler
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h)
		if err != nil {
			return
		}

		http.AllowQuerySemicolons(h)
	})
}

func Fuzz_CanonicalHeaderKey(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var s string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&s)
		if err != nil {
			return
		}

		http.CanonicalHeaderKey(s)
	})
}

func Fuzz_DetectContentType(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var data_0 []byte
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&data_0)
		if err != nil {
			return
		}

		http.DetectContentType(data_0)
	})
}

// skipping Fuzz_Error because parameters include unsupported type: net/http.ResponseWriter

// skipping Fuzz_FS because parameters include unsupported type: io/fs.FS

// skipping Fuzz_FileServer because parameters include unsupported type: net/http.FileSystem

// skipping Fuzz_FileServerFS because parameters include unsupported type: io/fs.FS

func Fuzz_Get(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url_0 string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&url_0)
		if err != nil {
			return
		}

		http.Get(url_0)
	})
}

func Fuzz_Handle(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pattern string
		var handler http.Handler
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&pattern, &handler)
		if err != nil {
			return
		}

		http.Handle(pattern, handler)
	})
}

func Fuzz_HandleFunc(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var pattern string
		var handler func(http.ResponseWriter, *http.Request)
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&pattern, &handler)
		if err != nil {
			return
		}

		http.HandleFunc(pattern, handler)
	})
}

func Fuzz_Head(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url_0 string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&url_0)
		if err != nil {
			return
		}

		http.Head(url_0)
	})
}

func Fuzz_ListenAndServe(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var addr string
		var handler http.Handler
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&addr, &handler)
		if err != nil {
			return
		}

		http.ListenAndServe(addr, handler)
	})
}

func Fuzz_ListenAndServeTLS(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var addr string
		var certFile string
		var keyFile string
		var handler http.Handler
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&addr, &certFile, &keyFile, &handler)
		if err != nil {
			return
		}

		http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	})
}

func Fuzz_MaxBytesHandler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Handler
		var n int64
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &n)
		if err != nil {
			return
		}

		http.MaxBytesHandler(h, n)
	})
}

// skipping Fuzz_MaxBytesReader because parameters include unsupported type: net/http.ResponseWriter

// skipping Fuzz_NewFileTransport because parameters include unsupported type: net/http.FileSystem

// skipping Fuzz_NewFileTransportFS because parameters include unsupported type: io/fs.FS

func Fuzz_NewRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var method string
		var url_0 string
		var body io.Reader
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&method, &url_0, &body)
		if err != nil {
			return
		}

		http.NewRequest(method, url_0, body)
	})
}

func Fuzz_NewRequestWithContext(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var ctx context.Context
		var method string
		var url_0 string
		var body io.Reader
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&ctx, &method, &url_0, &body)
		if err != nil {
			return
		}

		http.NewRequestWithContext(ctx, method, url_0, body)
	})
}

// skipping Fuzz_NewResponseController because parameters include unsupported type: net/http.ResponseWriter

// skipping Fuzz_NotFound because parameters include unsupported type: net/http.ResponseWriter

func Fuzz_ParseCookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var line string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&line)
		if err != nil {
			return
		}

		http.ParseCookie(line)
	})
}

func Fuzz_ParseHTTPVersion(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var vers string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&vers)
		if err != nil {
			return
		}

		http.ParseHTTPVersion(vers)
	})
}

func Fuzz_ParseSetCookie(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var line string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&line)
		if err != nil {
			return
		}

		http.ParseSetCookie(line)
	})
}

func Fuzz_ParseTime(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var text string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&text)
		if err != nil {
			return
		}

		http.ParseTime(text)
	})
}

func Fuzz_Post(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url_0 string
		var contentType string
		var body io.Reader
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&url_0, &contentType, &body)
		if err != nil {
			return
		}

		http.Post(url_0, contentType, body)
	})
}

func Fuzz_PostForm(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url_0 string
		var data_0 url.Values
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&url_0, &data_0)
		if err != nil {
			return
		}

		http.PostForm(url_0, data_0)
	})
}

func Fuzz_ProxyFromEnvironment(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var req *http.Request
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&req)
		if err != nil || req == nil {
			return
		}

		http.ProxyFromEnvironment(req)
	})
}

func Fuzz_ProxyURL(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var fixedURL *url.URL
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&fixedURL)
		if err != nil || fixedURL == nil {
			return
		}

		http.ProxyURL(fixedURL)
	})
}

func Fuzz_ReadRequest(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var b *bufio.Reader
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&b)
		if err != nil || b == nil {
			return
		}

		http.ReadRequest(b)
	})
}

func Fuzz_ReadResponse(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var r *bufio.Reader
		var req *http.Request
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&r, &req)
		if err != nil || r == nil || req == nil {
			return
		}

		http.ReadResponse(r, req)
	})
}

// skipping Fuzz_Redirect because parameters include unsupported type: net/http.ResponseWriter

func Fuzz_RedirectHandler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var url_0 string
		var code int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&url_0, &code)
		if err != nil {
			return
		}

		http.RedirectHandler(url_0, code)
	})
}

// skipping Fuzz_Serve because parameters include unsupported type: net.Listener

// skipping Fuzz_ServeContent because parameters include unsupported type: net/http.ResponseWriter

// skipping Fuzz_ServeFile because parameters include unsupported type: net/http.ResponseWriter

// skipping Fuzz_ServeFileFS because parameters include unsupported type: net/http.ResponseWriter

// skipping Fuzz_ServeTLS because parameters include unsupported type: net.Listener

// skipping Fuzz_SetCookie because parameters include unsupported type: net/http.ResponseWriter

func Fuzz_StatusText(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var code int
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&code)
		if err != nil {
			return
		}

		http.StatusText(code)
	})
}

func Fuzz_StripPrefix(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var prefix string
		var h http.Handler
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&prefix, &h)
		if err != nil {
			return
		}

		http.StripPrefix(prefix, h)
	})
}

func Fuzz_TimeoutHandler(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var h http.Handler
		var dt time.Duration
		var msg string
		fz := fuzzer.NewFuzzerV2(data, FabricFuncsForCustomTypes)
		err := fz.Fill2(&h, &dt, &msg)
		if err != nil {
			return
		}

		http.TimeoutHandler(h, dt, msg)
	})
}

func fabric_interface_httpRoundTripper(impl *http.Transport) http.RoundTripper {
	return impl
}

func fabric_interface_httpHandler(impl *http.ServeMux) http.Handler {
	return impl
}

func fabric_interface_empty(
	num int,
	transport http.Transport,
	server http.Server,
	servemux http.ServeMux,
	responsecontroller http.ResponseController,
	response http.Response,
	request http.Request,
	pushoptions http.PushOptions,
	protocolerror http.ProtocolError,
	maxbyteserror http.MaxBytesError,
	cookie http.Cookie,
	client http.Client,
) interface{} {
	switch num % 11 {
	case 0:
		return transport
	case 1:
		return server
	case 2:
		return servemux
	case 3:
		return responsecontroller
	case 4:
		return response
	case 5:
		return request
	case 6:
		return pushoptions
	case 7:
		return protocolerror
	case 8:
		return maxbyteserror
	case 9:
		return cookie
	case 10:
		return client
	default:
		panic("unreachable")
	}
}

func fabric_func_79_1() func(http.ResponseWriter, *http.Request) {
	return http.NotFound
}

var FabricFuncsForCustomTypes map[string][]reflect.Value

func TestMain(m *testing.M) {
	FabricFuncsForCustomTypes = make(map[string][]reflect.Value)
	FabricFuncsForCustomTypes["http.RoundTripper"] = append(FabricFuncsForCustomTypes["http.RoundTripper"], reflect.ValueOf(fabric_interface_httpRoundTripper))
	FabricFuncsForCustomTypes["interface {}"] = append(FabricFuncsForCustomTypes["interface {}"], reflect.ValueOf(fabric_interface_empty))
	FabricFuncsForCustomTypes["http.Handler"] = append(FabricFuncsForCustomTypes["http.Handler"], reflect.ValueOf(fabric_interface_httpHandler))
	FabricFuncsForCustomTypes["func(http.ResponseWriter, *http.Request)"] = append(FabricFuncsForCustomTypes["func(http.ResponseWriter, *http.Request)"], reflect.ValueOf(fabric_func_79_1))
	FabricFuncsForCustomTypes["http.Handler"] = append(FabricFuncsForCustomTypes["http.Handler"], reflect.ValueOf(http.TimeoutHandler))
	FabricFuncsForCustomTypes["http.Handler"] = append(FabricFuncsForCustomTypes["http.Handler"], reflect.ValueOf(http.StripPrefix))
	FabricFuncsForCustomTypes["http.Handler"] = append(FabricFuncsForCustomTypes["http.Handler"], reflect.ValueOf(http.RedirectHandler))
	FabricFuncsForCustomTypes["*http.Response"] = append(FabricFuncsForCustomTypes["*http.Response"], reflect.ValueOf(http.ReadResponse))
	FabricFuncsForCustomTypes["*http.Request"] = append(FabricFuncsForCustomTypes["*http.Request"], reflect.ValueOf(http.ReadRequest))
	FabricFuncsForCustomTypes["*url.URL"] = append(FabricFuncsForCustomTypes["*url.URL"], reflect.ValueOf(http.ProxyFromEnvironment))
	FabricFuncsForCustomTypes["*http.Response"] = append(FabricFuncsForCustomTypes["*http.Response"], reflect.ValueOf(http.PostForm))
	FabricFuncsForCustomTypes["*http.Response"] = append(FabricFuncsForCustomTypes["*http.Response"], reflect.ValueOf(http.Post))
	FabricFuncsForCustomTypes["time.Time"] = append(FabricFuncsForCustomTypes["time.Time"], reflect.ValueOf(http.ParseTime))
	FabricFuncsForCustomTypes["*http.Cookie"] = append(FabricFuncsForCustomTypes["*http.Cookie"], reflect.ValueOf(http.ParseSetCookie))
	FabricFuncsForCustomTypes["http.Handler"] = append(FabricFuncsForCustomTypes["http.Handler"], reflect.ValueOf(http.NotFoundHandler))
	FabricFuncsForCustomTypes["*http.ServeMux"] = append(FabricFuncsForCustomTypes["*http.ServeMux"], reflect.ValueOf(http.NewServeMux))
	FabricFuncsForCustomTypes["*http.Request"] = append(FabricFuncsForCustomTypes["*http.Request"], reflect.ValueOf(http.NewRequestWithContext))
	FabricFuncsForCustomTypes["*http.Request"] = append(FabricFuncsForCustomTypes["*http.Request"], reflect.ValueOf(http.NewRequest))
	FabricFuncsForCustomTypes["http.Handler"] = append(FabricFuncsForCustomTypes["http.Handler"], reflect.ValueOf(http.MaxBytesHandler))
	FabricFuncsForCustomTypes["error"] = append(FabricFuncsForCustomTypes["error"], reflect.ValueOf(http.ListenAndServeTLS))
	FabricFuncsForCustomTypes["error"] = append(FabricFuncsForCustomTypes["error"], reflect.ValueOf(http.ListenAndServe))
	FabricFuncsForCustomTypes["*http.Response"] = append(FabricFuncsForCustomTypes["*http.Response"], reflect.ValueOf(http.Head))
	FabricFuncsForCustomTypes["*http.Response"] = append(FabricFuncsForCustomTypes["*http.Response"], reflect.ValueOf(http.Get))
	FabricFuncsForCustomTypes["http.Handler"] = append(FabricFuncsForCustomTypes["http.Handler"], reflect.ValueOf(http.AllowQuerySemicolons))
	m.Run()
}
