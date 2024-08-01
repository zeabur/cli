package util

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

func ValidateTemplate(templateSpec []byte) error {
	var templateSpecUnmarshaled any

	err := yaml.NewDecoder(bytes.NewReader(templateSpec)).Decode(&templateSpecUnmarshaled)
	if err != nil {
		return fmt.Errorf("unmarshal schema: %w", err)
	}

	loader := jsonschema.SchemeURLLoader{
		"https": &httpURLLoader{client: http.DefaultClient},
	}

	compiler := jsonschema.NewCompiler()
	compiler.UseLoader(loader)
	schema, err := compiler.Compile("https://schema.zeabur.app/template.json")
	if err != nil {
		return fmt.Errorf("compile schema: %w", err)
	}

	if err := schema.Validate(templateSpecUnmarshaled); err != nil {
		return fmt.Errorf("validate schema: %w", err)
	}

	return nil
}

type httpURLLoader struct {
	client *http.Client
}

func (l *httpURLLoader) Load(url string) (any, error) {
	client := l.client
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}
	defer resp.Body.Close()

	return jsonschema.UnmarshalJSON(resp.Body)
}
