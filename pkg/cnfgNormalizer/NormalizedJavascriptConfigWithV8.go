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
	"istio.io/mixer/pkg/config"
)

/*


import (
	"istio.io/mixer/pkg/attribute"
	"github.com/augustoroman/v8"
	"fmt"
)

type NormalizedJavascriptConfigWithV8 struct {
	v8Context *v8.Context
	reportMtd *v8.Value
}

// invoked at runtime
func (n NormalizedJavascriptConfigWithV8) Evalaute(requestBag *attribute.MutableBag,
	callBack func(kind string, val interface{})) [][]interface {} {

	tmpCallBack := func(in v8.CallbackArgs) (*v8.Value, error) {
		callBack(in.Arg(0).String(), in.Arg(1))
		return nil, nil
	}
	n.v8Context.Global().Set(callbackMtdName, n.v8Context.Bind(callbackMtdName, tmpCallBack))


	attribsV8Value, err := n.v8Context.Create(constructAttributesForJS(requestBag))
	if err != nil {
		fmt.Println("ERROR constructing/binding attribs object", err)
	}
	val, errFromJS := n.reportMtd.Call(nil, attribsV8Value)
	if errFromJS != nil {
		fmt.Println("ERROR FROM JS with v8 engine", errFromJS)
	}
	//var vresult interface{}
	//vresult,_ = val.Get("result")

	vresult2,_ := val.Get("result")
	var returnVal [][]interface {}
	for i := 0 ; ; i++ {
		k, err := vresult2.GetIndex(i)
		if err == nil {
			var objToInsert []interface{}
			p,_ := k.GetIndex(0)
			objToInsert = append(objToInsert, p.String())
			m,_ := k.GetIndex(1)
			objToInsert = append(objToInsert, m)
			returnVal = append(returnVal, objToInsert)
			fmt.Printf("##valueresult=%v Typeresult= %T\n", returnVal, returnVal)
		} else {
			break;
		}

	}
	fmt.Printf("**value=%v, type=%T, valueresult=%v Typeresult= %T\n", val, val, returnVal, returnVal)


	return nil
}
*/

func createNormalizedJavascriptConfigWithV8(js string) config.NormalizedConfig {
	//ctx := v8.NewIsolate().NewContext()
	//
	//_, err := ctx.Eval(js, "")
	//if err != nil {
	//	fmt.Println("ERROR parsing JS", err)
	//	//return nil
	//}
	//reportMtd, err := ctx.Global().Get("report")
	//if err != nil {
	//	fmt.Println("ERROR finding report method", err)
	//	//return nil
	//}
	//return NormalizedJavascriptConfigWithV8{v8Context: ctx, reportMtd: reportMtd}
	return nil
}
