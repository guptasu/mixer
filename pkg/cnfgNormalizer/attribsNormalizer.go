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
)
var (
	attributesDescriptor = map[string]dpb.ValueType{
		"response.http.code" : dpb.INT64,
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

func GetAttributesType() string {
	var knownAttributesFields bytes.Buffer
	var attributesTypeFields bytes.Buffer
	var constructorCode bytes.Buffer

	for attrName, attrType := range attributesDescriptor {
		attrUpperCamelCaseName := dotCaseToCamelCase(attrName)
		knownAttributesFields.WriteString(fmt.Sprintf(`%s: "%s" as "%s",`, attrUpperCamelCaseName, attrName, attrUpperCamelCaseName))

		attributesTypeFields.WriteString(fmt.Sprintf("%s: %s;\n", attrUpperCamelCaseName, valueTypeToJSType[attrType]))

		constructorCode.WriteString(fmt.Sprintf(`
		if (attribs.Get('%s')[1]) {
		  //this.attribsThatExists.add(KnownAttribute.%s);
		  this.%s = attribs.Get('%s')[0]
		}
		`, attrName, attrUpperCamelCaseName, attrUpperCamelCaseName, attrName))

	}
	knownAttributeClass := fmt.Sprintf(`
	    const KnownAttribute = {
	      %s
	    }`, knownAttributesFields.String())

	AttributesClass := fmt.Sprintf(`
	    class Attributes {
	      // All the well known attribute names.
	      %s
	      //attribsThatExists : Set<keyof typeof KnownAttribute> = new Set<keyof typeof KnownAttribute>();
	      constructor (attribs: any) {
	        // Fill the set of attribues that are part of the call (data is available inside the attribs).

	        %s
	      }
              //has (knownAttribute: keyof typeof KnownAttribute) {
              //  if (this.attribsThatExists.has(knownAttribute)) {
              //    return true;
              //  }
              //  return false;
              //}
	    }
            function ConstructAttributes(attr: any) : Attributes {
              return new Attributes(attr)
            }`, attributesTypeFields.String(), constructorCode.String())

	return knownAttributeClass + "\n" + AttributesClass
}