package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	degit "github.com/qiushiyan/degit/pkg"
	"github.com/spf13/cobra"
)

var templates = []string{"learlab/itell-strapi-demo/apps/cttc-poe"}

var createCmd = &cobra.Command{
	Use:   "create <dest> [-t|--template] <template>",
	Short: "Creates a new textbook volume with template",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := degit.ParseRepo(Template)
		if err != nil {
			err := survey.AskOne(&survey.Select{
				Message: fmt.Sprintf(
					"Template `%s` is invalid. You can chose from the standard itell templates listed below.",
					Template,
				),
				Options: templates,
			}, &Template)
			if err != nil {
				return err
			}
			repo, err = degit.ParseRepo(Template)
			if err != nil {
				return err
			}
		}

		source := path.Join(repo.URL, repo.Subdir)
		var confirm bool

		dst := args[0]
		err = survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("cloning %s to %s?", source, dst),
			Default: true,
		}, &confirm)

		if err != nil {
			return err
		}

		if !confirm {
			return nil
		}

		if err := repo.Clone(dst, Force, Verbose); err != nil {
			return err
		}

		var title string
		err = survey.AskOne(
			&survey.Input{Message: "Textbook title"},
			&title,
			survey.WithValidator(survey.Required),
		)
		if err != nil {
			return err
		}

		qs := []*survey.Question{
			{
				Name: "packagename",
				Prompt: &survey.Input{
					Message: "package name",
					Default: fmt.Sprintf(
						"@itell/%s",
						strings.ToLower(strings.ReplaceAll(title, " ", "-")),
					),
					Help: "The `name` field in package.json, e.g., @itell/cttc-poe",
				},
				Validate: survey.Required,
			},
			{
				Name: "description",
				Prompt: &survey.Input{
					Message: "description",
				},
				Validate: survey.Required,
			},
			{
				Name: "port",
				Prompt: &survey.Input{
					Message: "dev server port",
					Default: "3000",
				},
				Validate: func(val any) error {
					if _, err := strconv.Atoi(val.(string)); err != nil {
						return errors.New("port must be a number")
					}

					return nil
				},
			},
			{
				Name: "keepcontent",
				Prompt: &survey.Confirm{
					Message: "Keep the content in the template? (if no, you will need to add content yourself before starting the dev server)",
					Default: true,
					Help:    "If false, this will try to delete all files in the content/textbook folder",
				},
			},
			{
				Name: "env",
				Prompt: &survey.Multiline{
					Message: "set environmental variables, or leave empty to edit later",
				},
			},
		}

		var response struct {
			PackageName string
			Description string
			Port        string
			KeepContent bool
			Env         string
		}
		err = survey.Ask(qs, &response)
		if err != nil {
			return err
		}

		err = updatePackageJson(dst, response.PackageName, response.Port)
		if err != nil {
			return err
		}

		err = updateEnv(dst, response.Env)
		if err != nil {
			return err
		}

		err = updateConfig(dst, title, response.Description)
		if err != nil {
			return err
		}

		if !response.KeepContent {
			err = removeTextbookContent(dst)
			if err != nil {
				return err
			}
		}

		fmt.Printf(`
Textbook created successfully. Here are the next steps:

1. cd into %s and run "pnpm install"
2. Make necessary changes to .env. You will likely need to create a new supabase instance and set the DATABASE_URL environmental variable.
3. If the database is new, run "pnpm drizzle-kit push" to initialize the tables
4. Edit content/home.mdx to update the homepage text
5. Run "pnpm dev" to start the dev server

`, dst)
		return nil
	},
}

var Template string
var Force bool
var Verbose bool

func init() {
	createCmd.Flags().
		StringVarP(&Template, "template", "t", "learlab/itell-strapi-demo/apps/cttc-poe", "template directory")
	createCmd.Flags().BoolVarP(&Force, "force", "f", false, "overwrite destination")
	createCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "verbose mode")
	rootCmd.AddCommand(createCmd)
}

func updatePackageJson(dir string, pkgName string, port string) error {
	pkgPath := path.Join(dir, "package.json")
	if !exists(pkgPath) {
		return errors.New("no package.json found")
	}

	f, err := os.OpenFile(pkgPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var data = make(map[string]any)
	s, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(s, &data)
	if err != nil {
		return err
	}
	data["name"] = pkgName

	var scripts map[string]any
	if _, ok := data["scripts"].(map[string]any); ok {
		scripts = data["scripts"].(map[string]any)
		scripts["dev"] = fmt.Sprintf("next dev -p %s", port)
	} else {
		return errors.New("invalid package.json, no scripts found")
	}

	s, _ = json.Marshal(data)

	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	_, err = f.Write(s)
	return err
}

func updateConfig(dir string, title string, description string) error {
	path := path.Join(dir, "src", "config", "site.ts")
	s, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	titlePattern := regexp.MustCompile(`title:\s*"([^"]*)"`)
	descriptionPattern := regexp.MustCompile(`description:\s*"([^"]*)"`)

	contents := titlePattern.ReplaceAllString(string(s), `title: "`+title+`"`)
	contents = descriptionPattern.ReplaceAllString(
		contents,
		`description: "`+description+`"`,
	)

	return os.WriteFile(path, []byte(contents), 0644)
}

func updateEnv(dir string, data string) error {
	f, err := os.OpenFile(path.Join(dir, ".env"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	return err
}

func removeTextbookContent(dir string) error {
	return os.RemoveAll(path.Join(dir, "content", "textbook"))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
