package main

import "C"
import (
	"encoding/json"
	"fmt"
	"net/url"

	http "github.com/bogdanfinn/fhttp"
	tls_client_cffi_src "github.com/bogdanfinn/tls-client/cffi_src"
)

//export getCookiesFromSession
func getCookiesFromSession(getCookiesParams *C.char) *C.char {
	getCookiesParamsJson := C.GoString(getCookiesParams)

	cookiesInput := tls_client_cffi_src.GetCookiesFromSessionInput{}
	marshallError := json.Unmarshal([]byte(getCookiesParamsJson), &cookiesInput)

	if marshallError != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(marshallError)
		return handleErrorResponse("", false, clientErr)
	}

	tlsClient, err := tls_client_cffi_src.GetTlsClientFromSession(cookiesInput.SessionId)

	if err != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(err)
		return handleErrorResponse(cookiesInput.SessionId, true, clientErr)
	}

	u, parsErr := url.Parse(cookiesInput.Url)
	if parsErr != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(parsErr)
		return handleErrorResponse(cookiesInput.SessionId, true, clientErr)
	}

	cookies := tlsClient.GetCookies(u)

	jsonResponse, marshallError := json.Marshal(cookies)

	if marshallError != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(marshallError)
		return handleErrorResponse(cookiesInput.SessionId, true, clientErr)
	}

	return C.CString(string(jsonResponse))
}

//export request
func request(requestParams *C.char) *C.char {
	requestParamsJson := C.GoString(requestParams)

	requestInput := tls_client_cffi_src.RequestInput{}
	marshallError := json.Unmarshal([]byte(requestParamsJson), &requestInput)

	if marshallError != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(marshallError)
		return handleErrorResponse("", false, clientErr)
	}

	tlsClient, sessionId, withSession, err := tls_client_cffi_src.GetTlsClientFromInput(requestInput)

	if err != nil {
		return handleErrorResponse(sessionId, withSession, err)
	}

	req, err := tls_client_cffi_src.BuildRequest(requestInput)

	if err != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(err)
		return handleErrorResponse(sessionId, withSession, clientErr)
	}

	cookies := buildCookies(requestInput.RequestCookies)

	if len(cookies) > 0 {
		tlsClient.SetCookies(req.URL, cookies)
	}

	resp, reqErr := tlsClient.Do(req)

	if reqErr != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(fmt.Errorf("failed to do request: %w", reqErr))
		return handleErrorResponse(sessionId, withSession, clientErr)
	}

	sessionCookies := tlsClient.GetCookies(req.URL)

	response, err := tls_client_cffi_src.BuildResponse(sessionId, withSession, resp, sessionCookies, requestInput.IsByteResponse)
	if err != nil {
		return handleErrorResponse(sessionId, withSession, err)
	}

	jsonResponse, marshallError := json.Marshal(response)

	if marshallError != nil {
		clientErr := tls_client_cffi_src.NewTLSClientError(marshallError)
		return handleErrorResponse(sessionId, withSession, clientErr)
	}

	return C.CString(string(jsonResponse))
}

func handleErrorResponse(sessionId string, withSession bool, err *tls_client_cffi_src.TLSClientError) *C.char {
	response := tls_client_cffi_src.Response{
		Status:  0,
		Body:    err.Error(),
		Headers: nil,
		Cookies: nil,
	}

	if withSession {
		response.SessionId = sessionId
	}

	jsonResponse, marshallError := json.Marshal(response)

	if marshallError != nil {
		return C.CString(marshallError.Error())
	}

	return C.CString(string(jsonResponse))
}

func buildCookies(cookies []tls_client_cffi_src.CookieInput) []*http.Cookie {
	var ret []*http.Cookie

	for _, cookie := range cookies {
		ret = append(ret, &http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Path:    cookie.Path,
			Domain:  cookie.Domain,
			Expires: cookie.Expires.Time,
		})
	}

	return ret
}

func main() {

}
