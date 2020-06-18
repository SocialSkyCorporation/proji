package cmd

import (
	"fmt"

	"github.com/nikoksr/proji/pkg/helper"

	"github.com/nikoksr/proji/pkg/proji/storage/item"

	"github.com/spf13/cobra"
)

var removeAllClasses, forceRemoveClasses bool

var classRmCmd = &cobra.Command{
	Use:   "rm LABEL [LABEL...]",
	Short: "Remove one or more classes",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Collect classes that will be removed
		var classes []*item.Class

		if removeAllClasses {
			var err error
			classes, err = projiEnv.Svc.LoadAllClasses()
			if err != nil {
				return err
			}
		} else {
			if len(args) < 1 {
				return fmt.Errorf("missing class label")
			}

			for _, label := range args {
				classID, err := projiEnv.Svc.LoadClassIDByLabel(label)
				if err != nil {
					return err
				}
				class, err := projiEnv.Svc.LoadClass(classID)
				if err != nil {
					return err
				}
				classes = append(classes, class)
			}
		}

		// Remove the classes
		for _, class := range classes {
			// Skip default classes
			if class.IsDefault {
				continue
			}
			// Ask for confirmation if force flag was not passed
			if !forceRemoveClasses {
				if !helper.WantTo(
					fmt.Sprintf("Do you really want to remove class '%s (%s)'?", class.Name, class.Label),
				) {
					continue
				}
			}
			err := projiEnv.Svc.RemoveClass(class.ID)
			if err != nil {
				fmt.Printf("> Removing '%s' failed: %v\n", class.Label, err)
				return err
			}
			fmt.Printf("> '%s' was successfully removed\n", class.Label)
		}
		return nil
	},
}

func init() {
	classCmd.AddCommand(classRmCmd)
	classRmCmd.Flags().BoolVarP(&removeAllClasses, "all", "a", false, "Remove all classes")
	classRmCmd.Flags().BoolVarP(&forceRemoveClasses, "force", "f", false, "Don't ask for confirmation")
}