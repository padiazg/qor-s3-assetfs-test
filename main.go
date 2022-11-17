package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/padiazg/qor-render-s3/s3"
	"github.com/qor/render"
)

func main() {
	var (
		config = &s3.Config{
			AccessID:         "your-access-id",
			AccessKey:        "your-access-Key",
			Region:           "sa-east-1",
			Bucket:           "portal-test-pato",
			S3Endpoint:       "http://localhost:9000",
			S3ForcePathStyle: true,
		}

		// fs = assetfs.AssetFS()
		fs = s3.NewAssetFS(config)

		rnd = render.New(&render.Config{
			ViewPaths:       []string{"app/views"},
			DefaultLayout:   "application", // default value is application
			AssetFileSystem: fs,
			FuncMapMaker: func(render *render.Render, request *http.Request, writer http.ResponseWriter) template.FuncMap {
				// genereate FuncMap that could be used when render template based on request info
				funcMap := template.FuncMap{}
				return funcMap
			},
		})

		r = mux.NewRouter()
	)

	fs.RegisterPath("app/views")

	r.PathPrefix("/render").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rnd.Execute("test", nil, r, w)
	})

	fmt.Println("Listening on: 8000")
	http.ListenAndServe(":8000", r)

}
