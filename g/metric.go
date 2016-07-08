package g

import (
	"log"

	//"github.com/k0kubun/pp"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
)

func MetricFilter(metrics *[]*model.MetricValue) (err error) {
	defaultTagsMap := Config().DefaultTags
	ignoreMetrics := Config().IgnoreMetrics
	//ignoreMetricsWithTags := Config().IgnoreMetricsWithTags
	debug := Config().Debug

	mvs := *metrics
	for index, mv := range mvs {
		if b, ok := ignoreMetrics[mv.Metric]; ok && b {
			// delete mvs[index]
			mvs = append(mvs[:index], mvs[index+1:]...)
		} else {
			// merge defaultTags
			if defaultTagsMap != nil {
				sourceTagsMap := utils.DictedTagstring(mv.Tags)
				tagsMap := utils.MergeTagsMap(sourceTagsMap, defaultTagsMap)
				mv.Tags = utils.SortedTags(tagsMap)
			}
			if debug {
				log.Println("=> Metric: ", mv.Metric, " Tags: ", mv.Tags, "\n")
			}
		}
	}
	return nil
}
