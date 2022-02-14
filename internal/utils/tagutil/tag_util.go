package tagutil

import "github.com/lemoony/snipkit/internal/utils/stringutil"

func HasValidTag(validTagUUIDs stringutil.StringSet, snippetTagUUIDS []string) bool {
	if len(validTagUUIDs) == 0 {
		return true
	}

	for _, tagUUID := range snippetTagUUIDS {
		if validTagUUIDs.Contains(tagUUID) {
			return true
		}
	}
	return false
}
