package technomodule

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// APIData struct
type APIData struct {
	Prefix  string
	Version string
	API     string
}

// ResolveData struct
type ResolveData struct {
	Host    string `json:"host,omitempty"`
	Prefix  string `json:"prefix,omitempty"`
	Version string `json:"version,omitempty"`
}

// FuncData struct
type FuncData struct {
	Function string
	Writer   http.ResponseWriter
	Request  *http.Request
}

// ServiceProfile struct
type ServiceProfile struct {
	Profile string
	Action  string
}

// ServiceData struct
type ServiceData struct {
	Secret string
}

// variables
const EnvType = "TYPE"
const EnvDev = "DEV"
const Secret = "SECRET"
const ServiceProfileHeader = "ServiceProfile"
const HostPrefixVersion = "hostPrefixVersion"
const LoginEmp = "loginEmp"
const CheckSessionByToken = "checkSessionByToken"
const ServiceCheckSessionByToken = "serviceCheckSessionByToken"
const ServiceValidateAction = "serviceValidateAction"
const ServiceAPI = "services/"

var Api APIData
var TechnoIMGResolveData ResolveData
var Service ServiceData

// LogError writes log from function
func LogError(text string, function string, err error) {
	log.Printf("%s %s%s%s %s\n", text, Api.Prefix, Api.Version, function, err)
}

// GetResolveData function
func GetResolveData(data FuncData) (result ResolveData) {

	response, err := http.Get(TechnoIMGResolveData.Host + TechnoIMGResolveData.Prefix + TechnoIMGResolveData.Version + HostPrefixVersion + "/" + LoginEmp)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error getting host+prefix+version"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot get host+prefix+version'", data.Function, errWriter)
		}
		return
	}

	var decoder = json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error decoding host+prefix+version"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot decode host+prefix+version'", data.Function, errWriter)
		}
		return
	}

	return
}

// GetResolveDataService function
func GetResolveDataService(data FuncData, service string) (result ResolveData) {

	response, err := http.Get(TechnoIMGResolveData.Host + TechnoIMGResolveData.Prefix + TechnoIMGResolveData.Version + HostPrefixVersion + "/" + service)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error getting host+prefix+version"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot get host+prefix+version'", data.Function, errWriter)
		}
		return
	}

	var decoder = json.NewDecoder(response.Body)
	err = decoder.Decode(&result)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error decoding host+prefix+version"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot decode host+prefix+version'", data.Function, errWriter)
		}
		return
	}

	return
}

// GetToken get token from request
func GetToken(request *http.Request) (token string) {

	tokenRaw := request.Header.Get("Authorization")
	if tokenRaw == "" {
		return ""
	}

	tokenSplit := strings.Split(tokenRaw, " ")
	if len(tokenSplit) < 2 {
		return ""
	}

	return tokenSplit[1]
}

// SetToken set token to request
func SetToken(request *http.Request, token string) {

	request.Header.Add(ServiceProfileHeader, token)
}

// GetServiceProfile gets ServiceProfile from request
func GetServiceProfile(request *http.Request) (serviceProfile ServiceProfile) {

	profileRaw := request.Header.Get(ServiceProfileHeader)
	if profileRaw == "" {
		return
	}

	profileSplit := strings.Split(profileRaw, " ")
	if len(profileSplit) < 2 {
		return
	}

	serviceProfile.Profile = profileSplit[0]
	serviceProfile.Action = profileSplit[1]
	return
}

// SetServiceProfile set service profile to request
func SetServiceProfile(request *http.Request, serviceProfile ServiceProfile) {

	request.Header.Add(ServiceProfileHeader, serviceProfile.Profile+" "+serviceProfile.Action)
}

// ValidateToken function
func ValidateToken(data FuncData, resolve ResolveData, token string) (result bool) {

	validate, err := http.Get(resolve.Host + resolve.Prefix + resolve.Version + CheckSessionByToken + "/" + token)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error getting checkSessionByToken"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot get checkSessionByToken'", data.Function, errWriter)
		}
		return
	}

	var decoder = json.NewDecoder(validate.Body)
	err = decoder.Decode(&result)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error decoding checkSessionByToken"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot decode checkSessionByToken'", data.Function, errWriter)
		}
		return
	}

	return
}

// ValidateTokenService function
func ValidateTokenService(data FuncData, resolve ResolveData, mySecret string, token string, profile ServiceProfile) (result bool, err error) {

	// validate token
	validate, err := http.Get(resolve.Host + resolve.Prefix + resolve.Version + ServiceAPI + ServiceCheckSessionByToken + "/" + mySecret + "/" + token)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error getting checkSessionByToken"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot get checkSessionByToken'", data.Function, errWriter)
		}
		return
	}

	var decoder = json.NewDecoder(validate.Body)
	err = decoder.Decode(&result)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error decoding checkSessionByToken"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot decode checkSessionByToken'", data.Function, errWriter)
		}
		return
	}

	if !result {
		return
	}

	// validate profile
	validate, err = http.Get(resolve.Host + resolve.Prefix + resolve.Version + ServiceAPI + ServiceValidateAction + "/" + mySecret + "/" + profile.Profile + "/" + profile.Action)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error getting serviceValidateAction"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot get serviceValidateAction'", data.Function, errWriter)
		}
		return
	}

	decoder = json.NewDecoder(validate.Body)
	err = decoder.Decode(&result)
	if err != nil {
		data.Writer.WriteHeader(http.StatusInternalServerError)
		_, errWriter := data.Writer.Write([]byte("Error decoding serviceValidateAction"))
		if errWriter != nil {
			LogError("Error writing result 'Cannot decode serviceValidateAction'", data.Function, errWriter)
		}
		return
	}

	return
}
