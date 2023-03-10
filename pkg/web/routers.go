/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package web

import (
	"context"
	"fmt"
	"github.com/traas-stack/chaosmetad/pkg/web/handler"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(ctx context.Context, isPprof bool) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	if isPprof {
		routes = append(routes, pprofRoutes...)
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(ctx, handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/",
		Index,
	},

	Route{
		"ExperimentInjectPost",
		strings.ToUpper("Post"),
		"/v1/experiment/inject",
		handler.ExperimentInjectPost,
	},

	Route{
		"ExperimentQueryPost",
		strings.ToUpper("Post"),
		"/v1/experiment/query",
		handler.ExperimentQueryPost,
	},

	Route{
		"ExperimentRecoverPost",
		strings.ToUpper("Post"),
		"/v1/experiment/recover",
		handler.ExperimentRecoverPost,
	},

	Route{
		"StatusGet",
		strings.ToUpper("Get"),
		"/v1/status",
		handler.StatusGet,
	},
}

var pprofRoutes = Routes{
	Route{
		"DebugIndex",
		strings.ToUpper("Get"),
		"/debug/pprof/",
		pprof.Index,
	},

	Route{
		"DebugProfile",
		strings.ToUpper("Get"),
		"/debug/pprof/profile",
		pprof.Profile,
	},

	Route{
		"DebugHeap",
		strings.ToUpper("Get"),
		"/debug/pprof/heap",
		pprof.Handler("heap").ServeHTTP,
	},

	Route{
		"DebugBlock",
		strings.ToUpper("Get"),
		"/debug/pprof/block",
		pprof.Handler("block").ServeHTTP,
	},

	Route{
		"DebugGoroutine",
		strings.ToUpper("Get"),
		"/debug/pprof/goroutine",
		pprof.Handler("goroutine").ServeHTTP,
	},

	Route{
		"DebugAllocs",
		strings.ToUpper("Get"),
		"/debug/pprof/allocs",
		pprof.Handler("allocs").ServeHTTP,
	},

	Route{
		"DebugCmdline",
		strings.ToUpper("Get"),
		"/debug/pprof/cmdline",
		pprof.Cmdline,
	},

	Route{
		"DebugThreadcreate",
		strings.ToUpper("Get"),
		"/debug/pprof/threadcreate",
		pprof.Handler("threadcreate").ServeHTTP,
	},

	Route{
		"DebugMutex",
		strings.ToUpper("Get"),
		"/debug/pprof/mutex",
		pprof.Handler("mutex").ServeHTTP,
	},

	Route{
		"DebugTrace",
		strings.ToUpper("Get"),
		"/debug/pprof/trace",
		pprof.Trace,
	},
}
