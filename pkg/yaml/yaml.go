package yaml

import (
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

func (r *Request) getEnvVars() *[]string {
	re := regexp.MustCompile(URL_VAR_PATTERN)

	matches := re.FindAllStringSubmatch(r.URL, -1)

	if len(matches) == 0 {
		log.Println("No variables found in URL")
		return nil
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
	return &uniqueList
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

func (gc *GlobalConfig) resolveRequest(requestName string) []map[string]string {
	var req *Request
	fmt.Printf("REQ BEFORE: %s\n", req)
	var endpointEnv map[string]string
	for _, eg := range gc.Endpoints {
		k, ok := eg.Requests[requestName]
		if ok {
			req = &k
			endpointEnv = eg.Env
			break
		}
	}
	fmt.Printf("REQ AFTER: %s\n", req)
	envVars := req.getEnvVars()
	envVarsFound := 0
	resolvedEnvVars := map[string]string{}
	for _, ev := range *envVars {
		if envVarsFound == len(*envVars) {
			break
		}
		// Check if env vars are in local Request Scope
		if varible, found := req.Env[ev]; found == true {
			resolvedEnvVars[ev] = varible
			envVarsFound += 1
			continue
		}

		// Check if env vars are in Endpoint Group scope
		if varible, found := endpointEnv[ev]; found == true {
			resolvedEnvVars[ev] = varible
			envVarsFound += 1
			continue
		}

		// Check if env vars are in global scope
		if varible, found := gc.GlobalEnv[ev]; found == true {
			resolvedEnvVars[ev] = varible
			envVarsFound += 1
			continue
		}

		log.Fatalf("Env var (%s) not found - aborting\n", ev)
		// if not in any scope log error
	}
	fmt.Printf("envVars: %v+\n", envVars)
	fmt.Printf("resolvedMaps: %v+\n", resolvedEnvVars)

	fmt.Printf("endpointEnv AFTER: %s\n", endpointEnv)

	// Take resolved env Vars and substitute them in to the URL string for the correct values
	realUrl := req.resolveUrlWithEnvVars(resolvedEnvVars)
	fullUrl := fmt.Sprintf("%s/%s", gc.BaseURL, realUrl) // TODO: parse the BaseURL for envs vars
	fmt.Printf("ResolvedURL BEFORE: %s\n", req.URL)
	fmt.Printf("ResolvedURL AFTER: %s\n", realUrl)
	fmt.Printf("FullUrl: %s\n", fullUrl)

	return nil
}

func ParseYaml(yamlData *string) {
	var cfg GlobalConfig
	if err := yaml.Unmarshal([]byte(*yamlData), &cfg); err != nil {
		log.Fatalf("failed unmarshalling yaml data: %s\n", err)
	}

	log.Printf("Successfully unmarshalled yamldata")

	cfg.resolveRequest("get_all_users")
}
