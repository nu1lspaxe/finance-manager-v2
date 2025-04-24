package server

import (
	"net/http"
	"path"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
)

func openAPIServer(dir string) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
			grpclog.Errorf("Not Found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		grpclog.Infof("Serving %s", r.URL.Path)
		p := strings.TrimPrefix(r.URL.Path, "/openapiv2/")
		p = path.Join(dir, p)
		http.ServeFile(w, r, p)
	}
}
