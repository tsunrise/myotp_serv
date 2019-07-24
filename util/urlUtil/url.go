package urlUtil

import "strings"

func MatchExact(urlPath string, target string) bool {
	return cutFirstEndSlash(urlPath) == cutFirstEndSlash(target)
}

func Match(path string, target string) bool {
	path = cutFirstEndSlash(path)
	target = cutFirstEndSlash(target)

	pathArr := strings.Split(path, "/")
	targetArr := strings.Split(target, "/")

	if len(pathArr) < len(targetArr) {
		return false
	}

	for i, v := range targetArr {
		if pathArr[i] != v && v != "" {
			return false
		}
	}

	return true
}

func cutFirstEndSlash(urlPath string) string {
	if len(urlPath) == 0 {
		return urlPath
	}
	if end := urlPath[len(urlPath)-1]; end == '/' {
		urlPath = urlPath[:len(urlPath)-1]
	}
	if len(urlPath) == 0 {
		return urlPath
	}
	if start := urlPath[0]; start == '/' {
		urlPath = urlPath[1:]
	}
	return urlPath
}
