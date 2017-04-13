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
	"istio.io/mixer/pkg/expr"
	"fmt"
)

func getJSForExpression(expression string) string{
	ex, err := expr.Parse(expression)
	var out string
	if err != nil {
		out = "false"
	} else {
		condition, _ := EvalJSExpession(ex, expr.FuncMap(), "attributes.")
		out = condition
	}
	return out
}

// Eval evaluates the expression given an attribute bag and a function map.
func EvalJSExpession(e *expr.Expression, fMap map[string]expr.FuncBase, getPropMtdName string) (string, error) {
	if e.Const != nil {
		return e.Const.StrValue, nil
	}
	if e.Var != nil {
		return fmt.Sprintf(getPropMtdName + "%s", getAttributeFieldName(e.Var.Name)), nil
	}

	fn := fMap[e.Fn.Name]
	if fn == nil {
		return "", fmt.Errorf("unknown function: %s", e.Fn.Name)
	}
	// may panic
	if e.Fn.Name == "EQ" {
		leftStr, _ := EvalJSExpession(e.Fn.Args[0], fMap, getPropMtdName)
		rightStr, _ := EvalJSExpession(e.Fn.Args[1], fMap, getPropMtdName)
		return fmt.Sprintf("%s == %s", leftStr, rightStr), nil
	}
	if e.Fn.Name == "OR" {
		allArgs := e.Fn.Args
		if len(allArgs) > 0 {
			var chkIfExists string
			leftexp,_ := EvalJSExpession(e.Fn.Args[0], fMap, getPropMtdName)
			chkIfExists = fmt.Sprintf("%s !== undefined", leftexp)
			rightexp, _ := EvalJSExpession(e.Fn.Args[1], fMap, getPropMtdName)
			return fmt.Sprintf("%s ? %s : %s", chkIfExists, leftexp, rightexp), nil
		}
	}
	return "", nil
}