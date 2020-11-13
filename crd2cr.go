package crd2cr

import (
	"fmt"

	"sigs.k8s.io/yaml"

	yaml2json "github.com/ghodss/yaml"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Convert(data []byte) ([]byte, error) {
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
	res.SetName(fmt.Sprintf("%sInstance", s.Spec.Names.Singular))
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
	case "array":
		if value.Items.Schema != nil {
			value := getValue(*value.Items.Schema)
			return []interface{}{value}
		}
		return ""
	case "string":
		if len(value.Enum) > 0 {
			return string(value.Enum[0].Raw)
		}
		return ""
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
