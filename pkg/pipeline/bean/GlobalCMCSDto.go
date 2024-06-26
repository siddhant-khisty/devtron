/*
 * Copyright (c) 2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bean

import (
	"encoding/json"
	"github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/internal/sql/repository"
)

type GlobalCMCSDto struct {
	Id         int    `json:"id"`
	ConfigType string `json:"configType" validate:"oneof=CONFIGMAP SECRET"`
	Name       string `json:"name"  validate:"required"`
	Type       string `json:"type" validate:"oneof=environment volume"`
	//map of key:value, example: '{ "a" : "b", "c" : "d"}'
	Data               map[string]string `json:"data"  validate:"required"`
	MountPath          string            `json:"mountPath"`
	Deleted            bool              `json:"deleted"`
	UserId             int32             `json:"-"`
	SecretIngestionFor string            `json:"SecretIngestionFor"` // value can be one of [ci, cd, ci/cd]
}

func (dto GlobalCMCSDto) ConvertToConfigSecretMap() (bean.ConfigSecretMap, error) {

	configSecretMap := bean.ConfigSecretMap{}
	configSecretMap.Name = dto.Name
	configSecretMap.Type = dto.Type
	configSecretMap.MountPath = dto.MountPath

	var jsonRawMsg []byte
	var err error
	// adding handling to get base64 encoded value in map value in case of secrets
	if dto.ConfigType == repository.CS_TYPE_CONFIG {
		jsonRawMsg, err = ConvertToEncodedForm(dto.Data)

	} else {
		jsonRawMsg, err = json.Marshal(dto.Data)
	}

	if err != nil {
		return configSecretMap, err
	}
	configSecretMap.Data = jsonRawMsg
	return configSecretMap, nil
}

// ConvertToEncodedForm Function to encode the values in the input map values
func ConvertToEncodedForm(data map[string]string) ([]byte, error) {
	var csDataMap = make(map[string][]byte)
	for key, value := range data {
		csDataMap[key] = []byte(value)
	}
	jsonRawMsg, err := json.Marshal(csDataMap)
	return jsonRawMsg, err
}
