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

// variables
const envType = "TYPE"
const envDev = "DEV"
const hostPrefixVersion = "hostPrefixVersion"
const loginEmp = "loginEmp"
const checkSessionByToken = "checkSessionByToken"

var api APIData
var technoIMGResolveData ResolveData

// LogError writes log from function
func LogError(text string, function string, err error) {
	log.Printf("%s %s%s%s %s\n", text, api.Prefix, api.Version, function, err)
}

// GetResolveData function
func GetResolveData(data FuncData) (result ResolveData) {

	response, err := http.Get(technoIMGResolveData.Host + technoIMGResolveData.Prefix + technoIMGResolveData.Version + hostPrefixVersion + "/" + loginEmp)
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

// ValidateToken function
func ValidateToken(data FuncData, resolve ResolveData, token string) (result bool) {

	validate, err := http.Get(resolve.Host + resolve.Prefix + resolve.Version + checkSessionByToken + "/" + token)
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
