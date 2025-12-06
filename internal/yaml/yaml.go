package yaml

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const URL_VAR_PATTERN = `\$\{([^}]+)\}`

type Request struct {
	Name string            `yaml:"name"`
	Type string            `yaml:"type"`
	URL  string            `yaml:"url"`
	Body string            `yaml:"body"`
	Env  map[string]string `yaml:"env"`
}

func (r *Request) getEnvVars() []string {
	re := regexp.MustCompile(URL_VAR_PATTERN)

	matches := re.FindAllStringSubmatch(r.URL, -1)

	if len(matches) == 0 {
		log.Println("No variables found in URL")
		return []string{}
	}

	extractedVars := make(map[string]struct{})

	for _, match := range matches {
		varName := match[1]
		extractedVars[varName] = struct{}{}
	}

	var uniqueList []string
	for v := range extractedVars {
		uniqueList = append(uniqueList, v)
	}
	return uniqueList
}

func (r *Request) resolveUrlWithEnvVars(envVars map[string]string) string {
	url := r.URL
	for envVar, realVal := range envVars {
		varToReplace := fmt.Sprintf("${%s}", envVar)
		url = strings.ReplaceAll(url, varToReplace, realVal)
	}
	return url
}

type EndpointGroup struct {
	Env      map[string]string  `yaml:"env"`
	Requests map[string]Request `yaml:"requests"`
}

type GlobalConfig struct {
	Project   string                   `yaml:"project"`
	BaseURL   string                   `yaml:"base_url"`
	GlobalEnv map[string]string        `yaml:"env"`
	Endpoints map[string]EndpointGroup `yaml:"endpoints"`
}

// Global Yaml file cache
var cfg GlobalConfig

func (gc *GlobalConfig) GetEndpoints() map[string]EndpointGroup {
	if gc == nil || len(gc.Endpoints) == 0 {
		return nil
	}

	// for _, v := range gc.Endpoints {
	// 	fmt.Printf("\nK::%s\n", v.Requests)
	// }

	return gc.Endpoints
}

// func (gc *GlobalConfig) GetEndpointsEnvs() map[string]string {
// 	var endpointEnv map[string]string
// 	for _, eg := range gc.Endpoints {
// 		r, ok := eg.Requests[requestName]
// 		if ok {
// 			req = &r
// 			endpointEnv = eg.Env
// 			break
// 		}
// 	}
// }

func (gc *GlobalConfig) ResolveRequest(requestName string) (string, error) {
	var req *Request
	var endpointEnv map[string]string
	for _, eg := range gc.Endpoints {
		r, ok := eg.Requests[requestName]
		if ok {
			req = &r
			endpointEnv = eg.Env
			break
		}
	}
	if req == nil {
		return "", errors.New(fmt.Sprintf("No requests found matching name %s\n", requestName))
	}
	envVars := req.getEnvVars()
	envVarsFound := 0
	resolvedEnvVars := map[string]string{}
	for _, ev := range envVars {
		if envVarsFound == len(envVars) {
			break
		}
		// Check if env vars are in local Request Scope
		if variable, found := req.Env[ev]; found == true {
			resolvedEnvVars[ev] = variable
			envVarsFound += 1
			continue
		}

		// Check if env vars are in Endpoint Group scope
		if variable, found := endpointEnv[ev]; found == true {
			resolvedEnvVars[ev] = variable
			envVarsFound += 1
			continue
		}

		// Check if env vars are in global scope
		if variable, found := gc.GlobalEnv[ev]; found == true {
			resolvedEnvVars[ev] = variable
			envVarsFound += 1
			continue
		}

		log.Fatalf("Env var (%s) not found - aborting\n", ev)
		// if not in any scope log error
	}
	fullUrl := req.resolveUrlWithEnvVars(resolvedEnvVars)
	if len(gc.BaseURL) > 0 {
		fullUrl = fmt.Sprintf("%s/%s", gc.BaseURL, fullUrl)
	}

	return fullUrl, nil
}

func ParseYaml(yamlData *string) GlobalConfig {
	if err := yaml.Unmarshal([]byte(*yamlData), &cfg); err != nil {
		log.Fatalf("failed unmarshalling yaml data: %s\n", err)
	}

	log.Printf("Successfully unmarshalled yaml data")
	cfg.GetEndpoints()
	return cfg
}
