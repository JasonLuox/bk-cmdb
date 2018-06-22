/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package validator

import (
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	api "configcenter/src/source_controller/api/object"
	"context"
	"encoding/json"
	"net/http"
)

type ValRule struct {
	ctx context.Context
	*backbone.Engine
	pheader        http.Header
	IsRequireArr   []string
	IsOnlyArr      []string
	AllFiledArr    []string
	NoEnumFiledArr []string
	PropertyKv     map[string]metadata.Attribute
	FieldRule      map[string]map[string]interface{}
	ownerID        string
	objID          string
	AllFieldAttDes []api.ObjAttDes
}

type MetaRst struct {
	Result  bool                     `json:result`
	Code    int                      `json:code`
	Message interface{}              `json:message`
	Data    []map[string]interface{} `json:data`
}

// NewValRule returns a new ValRule
func NewValRule(ownerID, objCtrl string) *ValRule {
	return &ValRule{ownerID: ownerID, objCtrl: objCtrl}
}

// GetObjAttrByID fetch object attributes
func (valid *ValRule) GetObjAttrByID(forward *api.ForwardParam, objID string) error {
	fieldRule := make(map[string]map[string]interface{})
	data := make(map[string]interface{})
	valid.PropertyKv = make(map[string]string)
	data[common.BKOwnerIDField] = valid.ownerID
	data[common.BKObjIDField] = objID
	info, _ := json.Marshal(data)
	client := api.NewClient(valid.objCtrl)

	re, _ := valid.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(valid.ctx, valid.pheader, dat)
	re.Data
	result, _ := client.SearchMetaObjectAttExceptInnerFiled(forward, []byte(info))
	valid.AllFieldAttDes = result
	blog.Infof("valid result:%+v selector:%s", result, info)
	for _, j := range result {
		cell := make(map[string]interface{})
		propertyID := j.PropertyID
		propertyType := j.PropertyType
		propertyName := j.PropertyName
		cell[common.BKPropertyTypeField] = j.PropertyType
		cell[common.BKOptionField] = j.Option
		if j.IsReadOnly {
			continue
		}
		fieldRule[propertyID] = cell
		isRequired := j.IsRequired
		isOnly := j.IsOnly
		if isRequired {
			valid.IsRequireArr = append(valid.IsRequireArr, propertyID)
		}
		if isOnly {
			valid.IsOnlyArr = append(valid.IsOnlyArr, propertyID)
		}
		if propertyType != common.FieldTypeEnum {
			valid.NoEnumFiledArr = append(valid.NoEnumFiledArr, propertyID)
		}
		valid.AllFiledArr = append(valid.AllFiledArr, propertyID)
		valid.PropertyKv[propertyID] = propertyName
	}
	valid.FieldRule = fieldRule
	return nil

}
