# [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) to CR


## Generate simple `Custom Resource` based on `Custom Resource Definition`

It's very hard to write your own first `CR` based on a `CRD` openAPIV3Schema, that's why I built this simple service that helps to do that by generating simple `CR` where you can put needed values. Supported types `array`, `string`, `integer`, `boolean` and `object`, `enum` also supported.

-----
only `apiextensions.k8s.io/v1` version supported


# Usage

```
Generate simple CR from CRD

Usage:
  crd2cr [flags]

Flags:
      --file string   file path. use STDIN by default
  -h, --help          help for crd2cr
```

[![asciicast](https://asciinema.org/a/aE9iwwmaETNyX9SboSCmVZY0R.png)](https://asciinema.org/a/aE9iwwmaETNyX9SboSCmVZY0R)


