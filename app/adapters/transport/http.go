package transport

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/utsavgupta/knowledge-hub/app/logger"
	"github.com/utsavgupta/knowledge-hub/app/runners"
	"github.com/utsavgupta/knowledge-hub/app/uc"
)

type httpRunner struct {
	port   int
	router http.Handler
}

type HttpRunnerDependencies struct {
	uc.SearchUc
	uc.ListDomainsUc
	uc.AddDomainUc
	uc.DeleteDomainUc
	uc.ListResourcesUc
	uc.AddResourceUc
	uc.DeleteResourceUc
}

func NewHttpRunner(port int, dependencies HttpRunnerDependencies) runners.Runner {

	router := mux.NewRouter()

	router.NewRoute().HandlerFunc(NewSearchHandler(dependencies.SearchUc)).Path("/search").Methods(http.MethodGet)
	router.NewRoute().HandlerFunc(NewListDomainsHandler(dependencies.ListDomainsUc)).Path("/domains").Methods(http.MethodGet)
	router.NewRoute().HandlerFunc(NewAddDomainHandler(dependencies.AddDomainUc)).Path("/domains").Methods(http.MethodPost)
	router.NewRoute().HandlerFunc(NewDeleteDomainHandler(dependencies.DeleteDomainUc)).Path("/domains/{domain_id}").Methods(http.MethodDelete)
	router.NewRoute().HandlerFunc(NewListResourcesHandler(dependencies.ListResourcesUc)).Path("/domains/{domain_id}/resources").Methods(http.MethodGet)
	router.NewRoute().HandlerFunc(NewAddResourceHandler(dependencies.AddResourceUc)).Path("/domains/{domain_id}/resources").Methods(http.MethodPost)
	router.NewRoute().HandlerFunc(NewDeleteResourceHandler(dependencies.DeleteResourceUc)).Path("/domains/{domain_id}/resources/{resource_id}").Methods(http.MethodDelete)

	return &httpRunner{port, router}
}

func (runner httpRunner) Run() error {

	logger.Instance().Info(context.Background(), fmt.Sprintf("Starting app server on port %d", runner.port))

	errChan := make(chan error)

	go func(c chan error) {

		if err := http.ListenAndServe(fmt.Sprintf(":%d", runner.port), runner.router); err != nil {
			c <- err
		}
	}(errChan)

	intChannel := make(chan os.Signal)
	signal.Notify(intChannel, os.Interrupt)

	select {
	case err := <-errChan:
		logger.Instance().Error(context.Background(), fmt.Sprintf("Exiting app server: %s", err.Error()))
		return err
	case <-intChannel:
		logger.Instance().Info(context.Background(), fmt.Sprintf("Stopping app server"))
	}

	return nil
}
