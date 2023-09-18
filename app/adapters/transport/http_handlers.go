package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/utsavgupta/knowledge-hub/app/entities"
	"github.com/utsavgupta/knowledge-hub/app/logger"
	"github.com/utsavgupta/knowledge-hub/app/uc"
)

type apiError struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func NewSearchHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func NewListDomainsHandler(listDomainsUc uc.ListDomainsUc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		domains, err := listDomainsUc(r.Context())

		if err != nil {
			handleError(w, r, err)
			return
		}

		sendResponse(w, r, http.StatusOK, domains)
	}
}

func NewAddDomainHandler(addDomainUc uc.AddDomainUc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		domain := &entities.Domain{}

		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(domain); err != nil {
			handleClientError(w, r, fmt.Errorf("invalid message body. please check documentation."))
			return
		}

		domain, err := addDomainUc(r.Context(), *domain)

		if err != nil {
			handleError(w, r, err)
			return
		}

		sendResponse(w, r, http.StatusCreated, *domain)
	}
}

func NewDeleteDomainHandler(deleteDomainUc uc.DeleteDomainUc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		domainId, ok := vars["domain_id"]

		if !ok {
			handleClientError(w, r, fmt.Errorf("domain id not provided"))
			return
		}

		if err := deleteDomainUc(r.Context(), domainId); err != nil {
			handleServerError(w, r, fmt.Errorf("internal server error"))
			return
		}

		sendResponse(w, r, http.StatusOK, "")
	}
}

func NewListResourcesHandler(listResourcesUc uc.ListResourcesUc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		domainId, ok := vars["domain_id"]

		if !ok {
			handleClientError(w, r, fmt.Errorf("domain id not provided"))
			return
		}

		resources, err := listResourcesUc(r.Context(), domainId)

		if err != nil {
			handleError(w, r, err)
			return
		}

		sendResponse(w, r, http.StatusOK, resources)
	}
}

func NewAddResourceHandler(addResourceUc uc.AddResourceUc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		domainId, ok := vars["domain_id"]

		if !ok {
			handleClientError(w, r, fmt.Errorf("domain id not provided"))
			return
		}

		resource := &entities.Resource{}

		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
			logger.Instance().Debug(r.Context(), err.Error())
			handleClientError(w, r, fmt.Errorf("invalid message body. please check documentation."))
			return
		}

		resource.DomainId = domainId

		resource, err := addResourceUc(r.Context(), *resource)

		if err != nil {
			handleError(w, r, err)
			return
		}

		sendResponse(w, r, http.StatusCreated, *resource)
	}
}

func NewDeleteResourceHandler(deleteResourceUc uc.DeleteResourceUc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		domainId, ok := vars["domain_id"]

		if !ok {
			handleClientError(w, r, fmt.Errorf("domain id not provided"))
			return
		}

		resourceId, ok := vars["resource_id"]

		if !ok {
			handleClientError(w, r, fmt.Errorf("resource id not provided"))
			return
		}

		resourceIdInt, err := strconv.Atoi(resourceId)

		if err != nil {
			handleClientError(w, r, fmt.Errorf("resource id should be an integer"))
			return
		}

		if err := deleteResourceUc(r.Context(), domainId, resourceIdInt); err != nil {
			handleServerError(w, r, fmt.Errorf("internal server error"))
			return
		}

		sendResponse(w, r, http.StatusOK, "")
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {

	if errors.Is(err, uc.ValidationError) {
		handleClientError(w, r, err)
		return
	}

	handleServerError(w, r, err)
}

func sendResponse(w http.ResponseWriter, r *http.Request, status int, body any) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Instance().Error(r.Context(), err.Error())
	}
}

func handleClientError(w http.ResponseWriter, r *http.Request, err error) {

	body := apiError{http.StatusBadRequest, err.Error()}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Instance().Error(r.Context(), err.Error())
	}
}

func handleServerError(w http.ResponseWriter, r *http.Request, err error) {

	body := apiError{http.StatusInternalServerError, "Internal Server Error"}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.Instance().Error(r.Context(), err.Error())
	}
}
