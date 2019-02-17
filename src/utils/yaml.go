package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ReadYAMLFileAsJSON : String -> String
func ReadYAMLFileAsJSON(yamlFile string) (string, error) {
	yamlString, _ := ioutil.ReadFile(yamlFile)
	yamlBuffer := []byte(yamlString)
	j2, err := yaml.YAMLToJSON(yamlBuffer)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", err
	}
	return string(j2), nil
}
