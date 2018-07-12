package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	historyUtil "github.com/tokopedia/historykv/src/util"
)

type httpKey struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type httpListKey struct {
	List         []httpKey `json:"list"`
	MessageError string    `json:"message_error"`
}

type httpHistoryKey struct {
	ID   int64  `json:"id"`
	By   string `json:"by"`
	Time string `json:"time"`
}

type httpHistoryByID struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	MessageError string `json:"message_error"`
}

type httpValueKey struct {
	Key          string           `json:"key"`
	Value        string           `json:"value"`
	History      []httpHistoryKey `json:"history"`
	MessageError string           `json:"message_error"`
}

type httpAdminUser struct {
	ID    int64  `json:"id"`
	User  string `json:"user"`
	Token string `json:"token"`
}

type httpAdminGetList struct {
	List         []httpAdminUser `json:"list"`
	HasNext      bool            `json:"has_next"`
	MessageError string          `json:"message_error"`
}

type httpSuccess struct {
	Success      bool   `json:"success"`
	MessageError string `json:"message_error"`
}

type httpResponse struct {
	User string      `json:"user"`
	Data interface{} `json:"data"`
}

func (h *HTTP) CreateResponse(user string, data interface{}) string {
	dataResponse := httpResponse{
		User: user,
		Data: data,
	}

	dataByte, errData := json.Marshal(dataResponse)
	if errData == nil {
		return string(dataByte)
	}

	return "{ \"user\": \"\", \"data\": {} }"
}

func (h *HTTP) CreateUpdateKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName == "" {
		returnData.MessageError = "NOT_LOGIN"
	} else {
		postKey := r.PostFormValue("key")
		postValue := r.PostFormValue("value")

		if postKey == "" || postKey == "/" {
			returnData.MessageError = fmt.Sprintf("Invalid Key: %s", postKey)
		} else {
			isFolder := postKey[len(postKey)-1] == '/'
			if isFolder {
				postValue = ""
			}

			isCreated := h.API.CreateUpdateKey(postKey, postValue, h.GetTokenFromRequest(r))
			if isCreated {
				if !isFolder {
					h.DB.AddHistory(postKey, postValue, userName, time.Now().Unix())
				}
				returnData.Success = true
			} else {
				returnData.MessageError = "Cannot save key. Check your token?"
			}
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) DeleteKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName == "" {
		returnData.MessageError = "NOT_LOGIN"
	} else {
		postKey := r.PostFormValue("key")

		if postKey == "" || postKey == "/" {
			returnData.MessageError = fmt.Sprintf("Invalid Key: %s", postKey)
		} else {
			isDeleted := h.API.DeleteKey(postKey, h.GetTokenFromRequest(r))
			if isDeleted {
				if postKey[len(postKey)-1] != '/' {
					h.DB.AddHistory(postKey, "[DELETED]", userName, time.Now().Unix())
				}
				returnData.Success = true
			} else {
				returnData.MessageError = "Cannot delete key. Check your token?"
			}
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) GetValueKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	getKey := r.FormValue("key")

	returnData := httpValueKey{
		Key:          getKey,
		Value:        "",
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName == "" {
		returnData.MessageError = "NOT_LOGIN"
	} else {
		returnData.Value = h.API.GetKeyValue(getKey, h.GetTokenFromRequest(r))

		getHistoryKey := h.DB.GetHistoryList(getKey)
		for _, dataHistory := range getHistoryKey {
			var histData httpHistoryKey
			histData.ID = dataHistory["id"].(int64)
			histData.By = dataHistory["user"].(string)

			tempTime := dataHistory["time"].(int64)
			timeDeduct := time.Now().Unix() - tempTime

			timeSecondMaxSecond := int64(60)
			timeSecondMaxMinute := int64(timeSecondMaxSecond * 60)
			timeSecondMaxHour := int64(timeSecondMaxMinute * 24)

			if timeDeduct <= 1 {
				histData.Time = fmt.Sprint("A second ago.")
			} else if timeDeduct < timeSecondMaxSecond {
				histData.Time = fmt.Sprintf("%d seconds ago.", timeDeduct)
			} else if timeDeduct == timeSecondMaxSecond {
				histData.Time = fmt.Sprint("A minute ago.")
			} else if timeDeduct < timeSecondMaxMinute {
				histData.Time = fmt.Sprintf("%d minutes ago.", int(timeDeduct/timeSecondMaxSecond))
			} else if timeDeduct == timeSecondMaxMinute {
				histData.Time = fmt.Sprint("An hour ago.")
			} else if timeDeduct < timeSecondMaxHour {
				histData.Time = fmt.Sprintf("%d hours ago.", int(timeDeduct/timeSecondMaxMinute))
			} else if timeDeduct == timeSecondMaxHour {
				histData.Time = fmt.Sprint("A day ago.")
			} else {
				histData.Time = fmt.Sprintf("%d days ago.", int(timeDeduct/timeSecondMaxHour))
			}

			returnData.History = append(returnData.History, histData)
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpListKey{
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName == "" {
		returnData.MessageError = "NOT_LOGIN"
	} else {
		getKey := r.FormValue("key")

		if getKey == "/" {
			getKey = ""
		}

		if getKey != "" {
			if getKey[len(getKey)-1] == '/' {
				getKey = fmt.Sprintf("%s/", getKey)
			}
		}

		getList := h.API.GetListKey(getKey, h.GetTokenFromRequest(r))

		for _, dataList := range getList {
			var dataKey httpKey
			dataKey.Key = dataList["key"]
			dataKey.Type = dataList["type"]

			returnData.List = append(returnData.List, dataKey)
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) GetHistoryID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpHistoryByID{
		Key:          "",
		Value:        "",
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName == "" {
		returnData.MessageError = "NOT_LOGIN"
	} else {
		getID := r.FormValue("id")

		intID, intIDErr := strconv.ParseInt(getID, 10, 64)

		if intIDErr == nil && intID > 0 {
			key, value := h.DB.GetHistoryID(intID)

			if key == "" {
				returnData.MessageError = "History not found."
			} else {
				returnData.Key = key
				returnData.Value = value
			}
		} else {
			returnData.MessageError = "Invalid History ID."
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	if h.IsDisableLogin {
		returnData.MessageError = "Login is disabled"
		fmt.Fprintf(w, "%s", h.CreateResponse("anonymous", returnData))
		return
	}

	postUser := r.PostFormValue("user")
	postPass := r.PostFormValue("pass")

	passwordMD5 := historyUtil.CreateHashPassword(postUser, postPass)

	getUser := h.DB.GetUser(postUser, passwordMD5)
	if getUser == "" {
		returnData.MessageError = "Invalid User or Password!"
	} else {

		newCookie := http.Cookie{
			Name:    "SID_HKV",
			Value:   h.CreateSession(getUser),
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		}
		http.SetCookie(w, &newCookie)

		tokenCookie := http.Cookie{
			Name:    "ACL_TOKEN",
			Value:   h.DB.GetToken(getUser),
			Expires: time.Now().Add(10 * 365 * 24 * time.Hour),
			Path:    "/",
		}

		http.SetCookie(w, &tokenCookie)

		returnData.Success = true
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(getUser, returnData))
}

func (h *HTTP) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	if h.IsDisableLogin {
		returnData.MessageError = "Login is disabled"
		fmt.Fprintf(w, "%s", h.CreateResponse("anonymous", returnData))
		return
	}

	cookieSID, errCookieSID := r.Cookie("SID_HKV")
	if errCookieSID == nil {
		sid := cookieSID.Value
		h.Session.Delete(sid)

		newCookie := http.Cookie{
			Name:    "SID_HKV",
			Value:   "",
			Expires: time.Now().Add(-24 * time.Hour),
			Path:    "/",
		}
		http.SetCookie(w, &newCookie)
		returnData.Success = true
	} else {
		returnData.MessageError = "NOT_LOGIN"
	}

	fmt.Fprintf(w, "%s", h.CreateResponse("", returnData))
}

func (h *HTTP) ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName == "" {
		returnData.MessageError = "NOT_LOGIN"
	} else {
		postPass := r.PostFormValue("pass")
		if postPass != "" {
			passwordMD5 := historyUtil.CreateHashPassword(userName, postPass)
			if h.DB.UpdateUser(userName, passwordMD5) {
				returnData.Success = true
			} else {
				returnData.MessageError = "Failed to change password. Please try again later."
			}
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) ChangeToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName == "" {
		returnData.MessageError = "NOT_LOGIN"
	} else {
		postToken := r.PostFormValue("token")
		if h.DB.UpdateToken(userName, postToken) {
			tokenCookie := http.Cookie{
				Name:    "ACL_TOKEN",
				Value:   postToken,
				Expires: time.Now().Add(10 * 365 * 24 * time.Hour),
				Path:    "/",
			}

			http.SetCookie(w, &tokenCookie)

			returnData.Success = true
		} else {
			returnData.MessageError = "Failed to change token. Please try again later."
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) AdminGetUserList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpAdminGetList{
		HasNext:      false,
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName != "admin" {
		returnData.MessageError = "NOT_ADMIN"
	} else {
		getUser := r.FormValue("user")
		getLastID := r.FormValue("lastid")

		intLastID, errIntLastID := strconv.ParseInt(getLastID, 10, 64)

		if errIntLastID == nil && intLastID >= 0 {
			listUser, listHasNext := h.DB.GetUserList(getUser, intLastID)
			returnData.HasNext = listHasNext

			for _, dataList := range listUser {
				if dataList["user"].(string) == userName {
					continue
				}
				var tempUser httpAdminUser
				tempUser.ID = dataList["id"].(int64)
				tempUser.User = dataList["user"].(string)
				tempUser.Token = dataList["token"].(string)
				returnData.List = append(returnData.List, tempUser)
			}
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) AdminCreateChangeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName != "admin" {
		returnData.MessageError = "NOT_ADMIN"
	} else {
		postUser := r.PostFormValue("user")
		if postUser != "" && postUser != userName {
			postToken := r.PostFormValue("token")
			postPass := r.PostFormValue("pass")
			changePassword := false
			if postPass != "" {
				postPass = historyUtil.CreateHashPassword(postUser, postPass)
				changePassword = true
			}

			if h.DB.IsUserExists(postUser) {
				if changePassword {
					if h.DB.UpdateUser(postUser, postPass) {
						h.DB.UpdateToken(postUser, postToken)
						returnData.Success = true
					} else {
						returnData.MessageError = "Failed to update user password and token."
					}
				} else {
					if h.DB.UpdateToken(postUser, postToken) {
						returnData.Success = true
					} else {
						returnData.MessageError = "Failed to update user token."
					}
				}
			} else {
				if changePassword {
					if h.DB.AddUser(postUser, postPass, postToken) {
						returnData.Success = true
					} else {
						returnData.MessageError = "Failed to create new user."
					}
				} else {
					returnData.MessageError = "User not exists, please specify password for new user."
				}
			}
		} else {
			returnData.MessageError = "User is empty."
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}

func (h *HTTP) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	returnData := httpSuccess{
		Success:      false,
		MessageError: "",
	}

	userName := h.GetUserFromRequest(r)
	if userName != "admin" {
		returnData.MessageError = "NOT_ADMIN"
	} else {
		postUser := r.PostFormValue("user")

		if postUser != "" && postUser != userName {
			if h.DB.DeleteUser(postUser) {
				returnData.Success = true
			} else {
				returnData.MessageError = "Failed to delete user."
			}
		} else {
			returnData.MessageError = "User is empty."
		}
	}

	fmt.Fprintf(w, "%s", h.CreateResponse(userName, returnData))
}
