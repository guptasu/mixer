// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cnfgNormalizer

import (
	"bytes"
	"fmt"
	dpb "istio.io/api/mixer/v1/config/descriptor"
	"strings"
	"istio.io/mixer/pkg/attribute"
)
var (
	attributesDescriptor = map[string]dpb.ValueType{
		"response.code" : dpb.INT64,
		"response.latency" : dpb.INT64,
		"api.method" : dpb.STRING,
		"target.name" : dpb.STRING,
		"api.name" : dpb.STRING,
		"source.name" : dpb.STRING,
	}
	valueTypeToJSType = map[dpb.ValueType]string{
		dpb.INT64:  "number",
		dpb.STRING: "string",
	}
)

func getAttributeFieldName (dotSeparatedAttribName string) string {
	return dotCaseToCamelCase(dotSeparatedAttribName)
}
func dotCaseToCamelCase(s string) string {
	subStrs := strings.Split(s, ".")
	for i, subStr := range subStrs {
		subStrs[i] = strings.Title(subStr)
	}
	return strings.Join(subStrs, "")
}

func GetTypeFromAttributes() string {
	var attributesTypeFields bytes.Buffer

	for attrName, attrType := range attributesDescriptor {
		attrUpperCamelCaseName := dotCaseToCamelCase(attrName)
		attributesTypeFields.WriteString(fmt.Sprintf("%s: %s;\n", attrUpperCamelCaseName, valueTypeToJSType[attrType]))
	}

	AttributesClass := fmt.Sprintf(`
	    class Attributes {
	      // All the well known attribute names.
	      %s
	    }`, attributesTypeFields.String())

	return AttributesClass
}

func constructAttributesForJS(requestBag *attribute.MutableBag) map[string]interface{} {
	attribs := make(map[string]interface{})
	for _, attribName := range requestBag.Names() {
          attribs[dotCaseToCamelCase(attribName)],_ = requestBag.Get(attribName)
	}
	return attribs
}