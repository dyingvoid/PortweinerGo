package infrastructure

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-envparse"
	"gopkg.in/yaml.v3"
)

func getValue(variables map[string]string, key string) string {
	if value, ok := variables[key]; ok {
		return value
	}
	value := ""

	return value
}

func updateEnvironmentVariables(service map[string]interface{}, variables map[string]string) {
	if env, ok := service["environment"]; ok {
		if envMap, ok := env.(map[string]interface{}); ok {
			for key, value := range envMap {
				valueStr := value.(string)
				if strings.HasPrefix(valueStr, "${") && strings.HasSuffix(valueStr, "}") {
					envVarName := strings.Trim(valueStr, "${}")
					envValue := getValue(variables, envVarName)
					if envValue != "" {
						envMap[key] = envValue
					} else {
						log.Printf("Warning: Environment variable %s not set", envVarName)
					}
				}
			}
		}
	}
}

func ProduceFile(envPath string, filePath string, stackName string) error {
	file, err := os.Open(envPath)
	if err != nil {
		return fmt.Errorf("could not open file %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	if err != nil {
		return fmt.Errorf("could not close the file %v", err)
	}

	variables, err := envparse.Parse(file)
	if err != nil {
		return fmt.Errorf("could not read env variables. %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file. %v", err)
	}

	var compose map[string]interface{}
	err = yaml.Unmarshal(data, &compose)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	services, ok := compose["services"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse services section")
	}

	for _, svc := range services {
		if service, ok := svc.(map[string]interface{}); ok {
			updateEnvironmentVariables(service, variables)
		}
	}

	updatedData, err := yaml.Marshal(&compose)
	if err != nil {
		return fmt.Errorf("failed to marshal updated YAML: %v", err)
	}

	if err := os.MkdirAll("produced/"+stackName, 0777); err != nil {
		return fmt.Errorf("failed to create dir for stack %s, error: %v", stackName, err)
	}

	err = os.WriteFile("produced/"+stackName+"/docker-compose.yml", updatedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
