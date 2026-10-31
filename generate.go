package elogik

// openapi.json is the untouched spec from https://oms.ekan-roissy.fr/api/doc/swagger;
// fix_spec.py repairs what oapi-codegen cannot parse.
//go:generate sh -c "python3 -I scripts/fix_spec.py openapi.json > spec.fixed.json && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.1 -config oapi-codegen.yaml spec.fixed.json; status=$?; rm -f spec.fixed.json; exit $status"
