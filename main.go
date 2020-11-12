package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"sigs.k8s.io/yaml"

	yaml2json "github.com/ghodss/yaml"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func handler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := convert(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, string(data))
}

func webServer() error {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on localhost:8080")
	return http.ListenAndServe(":8080", nil)
}

func main() {
	if err := webServer(); err != nil {
		log.Fatal(err)
	}
	// data, err := ioutil.ReadAll(bufio.NewReader(os.Stdin))
	// data, err := ioutil.ReadFile("schema.yaml")
	// if err != nil {
	// log.Fatal(err)
	// }
	// data, err = convert(data)
	// if err != nil {
	// log.Fatal(err)
	// }
	// fmt.Println(string(data))
}

func convert(data []byte) ([]byte, error) {
	s := v1.CustomResourceDefinition{}
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	res := unstructured.Unstructured{}
	gvk := schema.GroupVersionKind{}
	gvk.Kind = s.Spec.Names.Kind
	gvk.Group = s.Spec.Group
	for _, version := range s.Spec.Versions {
		if !version.Served || !version.Storage {
			continue
		}
		gvk.Version = version.Name
		spec := getSpec(version.Schema.OpenAPIV3Schema)
		s := parseObject(spec.Properties)
		content := res.UnstructuredContent()
		unstructured.SetNestedMap(content, s, "spec")
		res.SetUnstructuredContent(content)
	}
	res.SetName(fmt.Sprintf("%s_instance", s.Spec.Names.Singular))
	res.SetGroupVersionKind(gvk)
	j, err := res.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return yaml2json.JSONToYAML(j)
}

func parseObject(obj map[string]v1.JSONSchemaProps) map[string]interface{} {
	res := map[string]interface{}{}
	for key, value := range obj {
		res[key] = getValue(value)
	}
	return res
}

func getValue(value v1.JSONSchemaProps) interface{} {
	switch value.Type {
	case "object":
		return parseObject(value.Properties)
	case "integer":
		fmt.Println(value.Minimum)
		if value.Minimum != nil {
			return *value.Minimum
		}
		return 0.0
	case "boolean":
		return false
	default:
		return ""
	}
}

func getSpec(schema *v1.JSONSchemaProps) v1.JSONSchemaProps {
	for key, value := range schema.Properties {
		if key == "spec" {
			return value
		}
	}
	return v1.JSONSchemaProps{}
}
