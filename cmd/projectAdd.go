//nolint:gochecknoglobals,gochecknoinits
package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nikoksr/proji/messages"

	"github.com/nikoksr/proji/storage/models"
	"github.com/pkg/errors"

	"github.com/nikoksr/proji/util"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add LABEL PATH",
	Short: "Add an existing project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("missing label or path")
		}

		path, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		if !util.DoesPathExist(path) {
			return fmt.Errorf("path %s does not exist", path)
		}

		label := strings.ToLower(args[0])

		err = addProject(label, path)
		if err != nil {
			return errors.Wrap(err, "failed to add project")
		}
		messages.Success("successfully added project at path %s", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addProject(label, path string) error {
	name := filepath.Base(path)
	pkg, err := session.StorageService.LoadPackage(label)
	if err != nil {
		return errors.Wrap(err, "failed to load package")
	}

	project := models.NewProject(name, path, pkg)
	err = session.StorageService.SaveProject(project)
	if err != nil {
		return errors.Wrap(err, "failed to save package")
	}
	return nil
}
