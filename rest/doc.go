// Tideland GoREST - REST
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package rest of Tideland GoREST provides types for the
// implementation of servers with a RESTful API. The business has to
// be implemented in types fullfilling the ResourceHandler interface.
// This basic interface only allows the initialization of the handler.
// More interesting are the other interfaces like GetResourceHandler
// which defines the Get() method for the HTTP request method GET.
// Others are for PUT, POST, HEAD, PATCH, DELETE, and OPTIONS. Their
// according methods get a Job as argument. It provides convenient
// helpers for the processing of the job.
//
//     type myHandler struct {
//         id string
//     }
//
//     func NewMyHandler(id string) rest.ResourceHandler {
//         return &myHandler{id}
//     }
//
//     func (h *myHandler) ID() string {
//         return h.id
//     }
//
//     func (h *myHandler) Init(env rest.Environment, domain, resource string) error {
//         // Nothing to do in this example.
//         return nil
//     }
//
//     // Get handles reading of resources, here simplified w/o
//     // error handling.
//     func (h *myHandler) Get(job rest.Job) (bool, error) {
//         id := job.ResourceID()
//         if id == "" {
//            all := model.GetAllData()
//            job.JSON(true).Write(all)
//            return true, nil
//         }
//         one := model.GetOneData(id)
//         job.JSON(true).Write(one)
//         return true, nil
//     }
//
// The processing methods return two values: a boolean and an error.
// The latter is pretty clear, it signals a job processing error. The
// boolean is more interesting. Registering a handler is based on a
// domain and a resource. The URL
//
//     /<DOMAIN>/<RESOURCE>
//
// leads to a handler, or even better, to a list of handlers. All
// are used as long as the returned boolean value is true. E.g. the
// first handler can check the authentication, the second one the
// authorization, and the third one does the business. Additionally
// the URL
//
//     /<DOMAIN>/<RESOURCE>/<ID>
//
// provides the resource identifier via Job.ResourceID().
//
// The handlers then are deployed to the Multiplexer which implements
// the Handler interface of the net/http package. So the typical order
// is
//
//     mux := rest.NewMultiplexer(ctx, cfg)
//
// to start the multiplexer with a given context and the configuration
// for the multiplexer. The configuration is using the Tideland Go
// Library etc.Etc, parameters can be found at the NewMultiplexer
// documentation. After creating the multiplexer call
//
//     mux.Register("domain", "resource-type-a", NewTypeAHandler("foo"))
//     mux.Register("domain", "resource-type-b", NewTypeBHandler("bar"))
//     mux.Register("admin", "user", NewUserManagementHandler())
//
// to register the handlers per domain and resource. The server then can
// be started by the standard
//
//     http.ListenAndServe(":8000", mux)
//
// Additionally further handlers can be registered or running ones
// removed during the runtime.
package rest

// EOF
