package funcs

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
)

// TODO: BTW, I'd like to move these functions to "github.com/open-falcon/common/model/metric.go"
func NewMetricValue(metric string, val interface{}, dataType string, tags ...string) *model.MetricValue {
	mv := model.MetricValue{
		Metric: metric,
		Value:  val,
		Type:   dataType,
	}

	size := len(tags)

	if size > 0 {
		if size > 0 {
			if 1 == size {
				mv.Tags = tags[0]
			} else {
				// tags need to be sorted
				err, m := splitTagsString(tags...)

				if nil != err {
					mv.Tags = utils.SortedTags(m)
				}

				// TODO: err msg is missed here
			}
		}
	}

	return &mv
}

func GaugeValue(metric string, val interface{}, tags ...string) *model.MetricValue {
	return NewMetricValue(metric, val, "GAUGE", tags...)
}

func CounterValue(metric string, val interface{}, tags ...string) *model.MetricValue {
	return NewMetricValue(metric, val, "COUNTER", tags...)
}

func splitTagsString(strSlice ...string) (err error, tags map[string]string) {
	err = nil
	tags = make(map[string]string)

	for _, s := range strSlice {
		err, tagsT := utils.SplitTagsString(s)
		if nil != err {
			return err, tags        // TODO: continue or return is a question.
		} else {
			for k, v := range tagsT {
				tagsT[k] = v
			}
		}
	}

	return
}