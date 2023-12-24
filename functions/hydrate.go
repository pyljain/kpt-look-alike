package functions

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func Hydrate(filename string, args []string) error {

	substitutions := make(map[string]string)

	for _, input := range args {
		fragments := strings.Split(input, "=")
		if len(fragments) != 2 {
			return fmt.Errorf("make sure to pass in arguments with the argument name, example : image=nginx")
		}

		substitutions[fragments[0]] = fragments[1]
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(file)
	var node yaml.Node

	// Decode the next token
	err = decoder.Decode(&node)
	if err != nil {
		return err
	}

	visitNode(&node, substitutions)
	manifestBytes, err := yaml.Marshal(node.Content[0])
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, manifestBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func visitNode(node *yaml.Node, substitutions map[string]string) {

	if node.Kind == yaml.ScalarNode {

		if strings.Contains(node.LineComment, "kla-set:") {
			klaSetter := strings.Trim(strings.Split(node.LineComment, "kla-set:")[1], " ")
			tmpl, err := template.New("kla-template").Parse(klaSetter)
			if err != nil {
				log.Printf("Error parsing template %s", err)
				return
			}

			buf := bytes.NewBuffer([]byte{})

			err = tmpl.Execute(buf, substitutions)
			if err != nil {
				return
			}

			node.Value = buf.String()

		}
	}

	if node.Content == nil {
		return
	}

	for _, child := range node.Content {
		visitNode(child, substitutions)
	}
}

// func convertKind(kind yaml.Kind) string {
// 	switch kind {
// 	case yaml.DocumentNode:
// 		return "DocumentNode"
// 	case yaml.MappingNode:
// 		return "MappingNode"
// 	case yaml.ScalarNode:
// 		return "ScalarNode"
// 	case yaml.SequenceNode:
// 		return "SequenceNode"
// 	default:
// 		return "Unknown"
// 	}
// }
