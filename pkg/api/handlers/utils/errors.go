package utils

import (
	"net/http"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/errorhandling"
	"github.com/containers/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	ErrLinkNotSupport = errors.New("Link is not supported")
)

// TODO: document the exported functions in this file and make them more
// generic (e.g., not tied to one ctr/pod).

// Error formats an API response to an error
//
// apiMessage and code must match the container API, and are sent to client
// err is logged on the system running the podman service
func Error(w http.ResponseWriter, code int, err error) {
	// Log detailed message of what happened to machine running podman service
	log.Infof("Request Failed(%s): %s", http.StatusText(code), err.Error())
	em := errorhandling.ErrorModel{
		Because:      (errors.Cause(err)).Error(),
		Message:      err.Error(),
		ResponseCode: code,
	}
	WriteJSON(w, code, em)
}

func VolumeNotFound(w http.ResponseWriter, name string, err error) {
	if errors.Cause(err) != define.ErrNoSuchVolume {
		InternalServerError(w, err)
	}
	Error(w, http.StatusNotFound, err)
}

func ContainerNotFound(w http.ResponseWriter, name string, err error) {
	if errors.Cause(err) != define.ErrNoSuchCtr {
		InternalServerError(w, err)
	}
	Error(w, http.StatusNotFound, err)
}

func ImageNotFound(w http.ResponseWriter, name string, err error) {
	if errors.Cause(err) != storage.ErrImageUnknown {
		InternalServerError(w, err)
	}
	Error(w, http.StatusNotFound, err)
}

func NetworkNotFound(w http.ResponseWriter, name string, err error) {
	if errors.Cause(err) != define.ErrNoSuchNetwork {
		InternalServerError(w, err)
	}
	Error(w, http.StatusNotFound, err)
}

func PodNotFound(w http.ResponseWriter, name string, err error) {
	if errors.Cause(err) != define.ErrNoSuchPod {
		InternalServerError(w, err)
	}
	Error(w, http.StatusNotFound, err)
}

func SessionNotFound(w http.ResponseWriter, name string, err error) {
	if errors.Cause(err) != define.ErrNoSuchExecSession {
		InternalServerError(w, err)
	}
	Error(w, http.StatusNotFound, err)
}

func SecretNotFound(w http.ResponseWriter, nameOrID string, err error) {
	if errors.Cause(err).Error() != "no such secret" {
		InternalServerError(w, err)
	}
	Error(w, http.StatusNotFound, err)
}

func ContainerNotRunning(w http.ResponseWriter, containerID string, err error) {
	Error(w, http.StatusConflict, err)
}

func InternalServerError(w http.ResponseWriter, err error) {
	Error(w, http.StatusInternalServerError, err)
}

func BadRequest(w http.ResponseWriter, key string, value string, err error) {
	e := errors.Wrapf(err, "failed to parse query parameter '%s': %q", key, value)
	Error(w, http.StatusBadRequest, e)
}

// UnsupportedParameter logs a given param by its string name as not supported.
func UnSupportedParameter(param string) {
	log.Infof("API parameter %q: not supported", param)
}
