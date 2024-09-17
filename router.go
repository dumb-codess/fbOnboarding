package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func router(cfg *Config, controller *Controller) http.Handler {
	route := mux.NewRouter()
	route.Use(corsMiddleware)
	route.Use(panicRecovery)

	route.Handle("/v1/auth", AppHandler(controller.GetToken)).Methods(http.MethodPost)
	route.Handle("/v1/onboard/check-submission", Auth(AppHandler(controller.CheckSubmission))).Methods(http.MethodGet)
	route.Handle("/v1/onboard/upload", Auth(AppHandler(controller.Uploadfile))).Methods(http.MethodPost)
	route.Handle("/v1/interaction/check-approval-status", Auth(AppHandler(controller.GetApprovalStatus))).Methods(http.MethodGet)
	route.Handle("/v1/internal/approve-form", AppHandler(controller.ApproveApplicationStaus)).Methods(http.MethodPost)

	return route
}
