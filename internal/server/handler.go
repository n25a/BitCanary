package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/n25a/BitCanary/internal/config"
	"github.com/n25a/BitCanary/internal/log"
)

func CanaryTestingHandler(writer http.ResponseWriter, request *http.Request) {
	userID, err := extractUserID(request.Header)
	if err != nil {
		log.Logger.Error("error in extracting user id", zap.Error(err))
	}

}

func extractUserID(header http.Header) (uint64, error) {
	key := header.Get(config.C.UserIDHeaderKey)
	if key == "" {
		return 0, errors.New("user id header key is empty")
	}

	userIDString := key
	if config.C.UserNestedKey != "" {
		var data map[string]interface{}
		err := json.Unmarshal([]byte(key), &data)
		if err != nil {
			return 0, err
		}

		nestedKeys := strings.Split(config.C.UserNestedKey, ".")
		for _, nestedKey := range nestedKeys {
			if data[nestedKey] == nil {
				return 0, errors.New("nested key not found: " + nestedKey)
			}

			if nestedKey == nestedKeys[len(nestedKeys)-1] {
				userIDString = data[nestedKey].(string)
			} else {
				data = data[nestedKey].(map[string]interface{})
			}
		}
	}

	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
