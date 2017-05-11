package YamlExprEvaluator

import (
	"istio.io/mixer/pkg/expr"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/adapter"
	multierror "github.com/hashicorp/go-multierror"
	"fmt"
)


// Metric aspect specific stuff
type MetricValue struct {
	// DATA that has to be computed using attributes parsing.
	// Fields that need attribute processing is identified via annotations
	// on the descriptor.
	Value interface{}
	Labels map[string]interface{}

	// Other fields that might be specified by the user which do not require attribute processing
	// ...
}

//////////////////////////////



type AspectReqObjectConversionInfo struct {
	Kind string
        dataMappingInfo DataMappingInfo
}

type DataMappingInfo struct {

}

type YamlExprToAdapterReqEvalutor struct {
	mapper expr.Evaluator
}

func (ev *YamlExprToAdapterReqEvalutor) Evaluate(attrs attribute.Bag) (interface{}, error) {
	result := &multierror.Error{}
	var values []adapter.Value

	//for name, md := range w.metadata {
	//	metricValue, err := ev.mapper.Eval(md.value, attrs)
	//	if err != nil {
	//		result = multierror.Append(result, fmt.Errorf("failed to eval metric value for metric '%s': %v", name, err))
	//		continue
	//	}
	//	labels, err := evalAll(md.labels, attrs, ev.mapper)
	//	if err != nil {
	//		result = multierror.Append(result, fmt.Errorf("failed to eval labels for metric '%s': %v", name, err))
	//		continue
	//	}
//
	//	// TODO: investigate either pooling these, or keeping a set around that has only its field's values updated.
	//	// we could keep a map[metric name]value, iterate over the it updating only the fields in each value
	//	values = append(values, adapter.Value{
	//		Definition: md.definition,
	//		Labels:     labels,
	//		// TODO: extract standard timestamp attributes for start/end once we det'm what they are
	//		//StartTime:   time.Now(),
	//		//EndTime:     time.Now(),
	//		MetricValue: metricValue,
	//	})
	//}


	return values, result.ErrorOrNil()
}

func evalAll(expressions map[string]string, attrs attribute.Bag, eval expr.Evaluator) (map[string]interface{}, error) {
	result := &multierror.Error{}
	labels := make(map[string]interface{}, len(expressions))
	for label, texpr := range expressions {
		val, err := eval.Eval(texpr, attrs)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("failed to construct value for label '%s': %v", label, err))
			continue
		}
		labels[label] = val
	}
	return labels, result.ErrorOrNil()
}
