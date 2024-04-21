package cmd

import "fmt"

func errFailedToReadConfig(path string, err error) error {
	return fmt.Errorf("failed to read config %s\n%s", path, err)
}

func errFailedToRenderConfig(path string, err error) error {
	return fmt.Errorf("failed to render config %s\n%s", path, err)
}

func errFailedToSaveTmpConfig(path string, err error) error {
	return fmt.Errorf("failed to save tmp config %s\n%s", path, err)
}

func errFailedToUnmarshalRenderedConfigAsYaml(path string, err error) error {
	return fmt.Errorf("failed to unmarshal rendered config as yaml %s\n%s", path, err)
}

func errFailedToReadTemplate(path string, err error) error {
	return fmt.Errorf("failed to read template %s\n%s", path, err)
}

func errFailedToRenderTemplate(path string, err error) error {
	return fmt.Errorf("failed to render template %s\n%s", path, err)
}

func errFailedToSaveResult(path string, err error) error {
	return fmt.Errorf("failed to save result to %s\n%s", path, err)
}
