package template

// RevisedTemplateTmpl defines the modified template proto with Type and InstanceParams
var RevisedTemplateTmpl = `// Copyright 2017 Istio Authors
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

// THIS FILE IS AUTOMATICALLY GENERATED.

syntax = "proto3";

package {{.PackageName}};

import "mixer/v1/config/descriptor/value_type.proto";
import "pkg/adapter/template/TemplateExtensions.proto";

option (istio.mixer.v1.config.template.template_variety) = {{.VarietyName}};
option (istio.mixer.v1.config.template.template_name) = "{{.Name}}";

{{.Comment}}

{{.TemplateMessage.Comment}}
message Type {
  {{range .TemplateMessage.Fields}}
  {{.Comment}}
  {{replacePrimitiveToValueType .ProtoType}} {{.ProtoName}} = {{.Number}};
  {{end}}
}

message InstanceParam {
  {{range .TemplateMessage.Fields}}
  {{replaceValueTypeToString .ProtoType}} {{.ProtoName}} = {{.Number}};
  {{end}}
}
`
