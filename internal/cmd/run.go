package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
	"gitlab.collabox.dev/go/nse/nsjet"
	"gopkg.in/yaml.v3"
)

func Run() {
	app := &cli.App{
		Name:    "yampl",
		Usage:   "utility for generating files with using template engine",
		Version: "v0.0.1",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "config.yml",
				DefaultText: "config.yml",
			},
			&cli.BoolFlag{
				Name:        "tmp",
				Value:       true,
				DefaultText: "true",
			},
			&cli.PathFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "output",
				DefaultText: "output",
			},
		},
		Action: func(context *cli.Context) error {
			return render(
				context.Path("config"),
				context.Path("output"),
				context.Bool("tmp"),
			)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

var (
	yamlErrReg  = regexp.MustCompile(`yaml: line (\\d+):`)
	plushErrReg = regexp.MustCompile(`line (\\d+):`)
)

func render(configPath, outPath string, tmp bool) error {
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return errFailedToReadConfig(configPath, err)
	}

	configRendered, err := renderTemplate(configPath, string(configContent), nil)
	if err != nil {
		return errFailedToRenderConfig(configPath, err)
	}

	tmpConfigPath := generateTmpConfigPath(configPath)
	var c Config

	err = yaml.Unmarshal([]byte(configRendered), &c)
	if err != nil {
		if !tmp {
			return errFailedToUnmarshalRenderedConfigAsYaml(configPath, err)
		}

		matches := yamlErrReg.FindStringSubmatch(err.Error())
		if len(matches) <= 1 {
			return errFailedToUnmarshalRenderedConfigAsYaml(configPath, err)
		}

		if tmpSaveErr := os.WriteFile(tmpConfigPath, []byte(configRendered), os.ModePerm); tmpSaveErr != nil {
			return errFailedToSaveTmpConfig(tmpConfigPath, tmpSaveErr)
		}

		return errFailedToUnmarshalRenderedConfigAsYaml(
			tmpConfigPath,
			fmt.Errorf("%s:%s:0: %s",
				tmpConfigPath,
				matches[1],
				strings.Replace(err.Error(), matches[0], "", 1),
			),
		)
	}

	_ = os.Remove(tmpConfigPath)

	templateContent, err := os.ReadFile(c.Template)
	if err != nil {
		return errFailedToReadTemplate(c.Template, err)
	}

	templateRendered, err := renderTemplate(c.Template, string(templateContent), c.Params)
	if err != nil {
		matches := plushErrReg.FindStringSubmatch(err.Error())
		if len(matches) <= 1 {
			return errFailedToRenderTemplate(c.Template, err)
		}

		return errFailedToRenderTemplate(
			c.Template, fmt.Errorf("%s:%s:0:%s",
				c.Template,
				matches[1],
				strings.Replace(err.Error(), matches[0], "", 1),
			))
	}

	if err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
		return errFailedToSaveResult(outPath, err)
	}

	if err = os.WriteFile(outPath, []byte(templateRendered), os.ModePerm); err != nil {
		return errFailedToSaveResult(outPath, err)
	}

	return nil
}

func renderTemplate(path, content string, vars map[string]any) (string, error) {
	l := nsjet.NewInMemLoader()

	l.Set(path, content)

	tempalte, err := nsjet.NewSet(l).GetTemplate(path)
	if err != nil {
		log.Println("тут")
		panic(err)
	}

	ctx := map[string]reflect.Value{}
	for k, v := range vars {
		ctx[k] = reflect.ValueOf(v)
	}

	result := bytes.NewBuffer(nil)
	if err := tempalte.Execute(result, ctx, nil); err != nil {
		return "", err
	}

	return result.String(), nil
}

func generateTmpConfigPath(configPath string) string {
	tmpConfigBase := filepath.Base(configPath)
	if parts := strings.Split(tmpConfigBase, "."); len(parts) > 1 {
		parts = append(parts[0:len(parts)-1], "tmp", parts[len(parts)-1])
		tmpConfigBase = strings.Join(parts, ".")
	} else {
		tmpConfigBase = tmpConfigBase + ".tmp"
	}

	return filepath.Join(
		filepath.Dir(configPath),
		tmpConfigBase,
	)
}
