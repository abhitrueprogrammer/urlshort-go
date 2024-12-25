package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dest, exist := pathsToUrls[r.URL.Path]
		fmt.Print(dest)
		if exist {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(w, r)
	})
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type URLMapping struct {
    Source string 
    Destination  string 
}
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var decodedYml []URLMapping

	err := yaml.Unmarshal(yml, &decodedYml) 
	if err != nil{
		return nil, fmt.Errorf("error marshalling the yaml: %v", err);
	}
	sourceDestinationMap := ConvertYamlToMap(decodedYml)
	return MapHandler(sourceDestinationMap, fallback),nil;
}

func ConvertYamlToMap(yml [] URLMapping) (map[string]string) {
	sourceDestinationMap := make(map[string]string, 0)
	for _, v := range yml {
		sourceDestinationMap[v.Source] = v.Destination
	}

	return sourceDestinationMap
}
