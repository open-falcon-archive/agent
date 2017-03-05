package g

import (
	"log"
	"regexp"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
)

var cachedRegexp map[string]*regexp.Regexp = make(map[string]*regexp.Regexp)

func cachedMatch(re string, tags string) bool {
	r, ok := cachedRegexp[re]
	if !ok {
		r = regexp.MustCompile(re)
		cachedRegexp[re] = r
	}
	return r.MatchString(tags)
}

func FilterMetrics(metrics *[]*model.MetricValue) *[]*model.MetricValue {
	addTags := Config().AddTags
	ignore := Config().Ignore
	debug := Config().Debug

	filtered := make([]*model.MetricValue, 0)

metricsLoop:
	for _, mv := range *metrics {
		for metricRe, tagsRe := range ignore {
			if cachedMatch(metricRe, mv.Metric) && cachedMatch(tagsRe, mv.Tags) {
				if debug {
					log.Println("=> Filtered metric", mv.Metric, "/", mv.Tags, "by rule", metricRe, "=>", tagsRe)
				}
				continue metricsLoop
			}
		}

		if addTags != nil {
			tags := utils.DictedTagstring(mv.Tags)
			for k, v := range addTags {
				if _, ok := tags[k]; !ok {
					tags[k] = v
				}
			}
			mv.Tags = utils.SortedTags(tags)
		}
		filtered = append(filtered, mv)
	}
	return &filtered
}
