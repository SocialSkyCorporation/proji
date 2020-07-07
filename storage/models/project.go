package models

import (
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/copy"
	"gorm.io/gorm"
)

// Project represents a project that was created by proji. It holds tags for gorm and toml defining its storage and
// export/import behaviour.
type Project struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_unq_project_path_deletedat,unique;"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"size:64"`
	Path      string         `gorm:"index:idx_unq_project_path_deletedat,unique;not null"`
	Class     *Class         `gorm:"ForeignKey:ID;References:ID"`
}

// NewProject returns a new project.
func NewProject(name, path string, class *Class) *Project {
	return &Project{
		Name:  name,
		Path:  path,
		Class: class,
	}
}

// Create starts the creation of a project.
func (p *Project) Create(cwd, baseConfigPath string) (err error) {
	err = p.createProjectFolder()
	if err != nil {
		return err
	}

	// Chdir into the new project folder and defer chdir back to old cwd
	err = os.Chdir(p.Path)
	if err != nil {
		return err
	}

	// Append a slash if not exists. Out of convenience.
	if cwd[:len(cwd)-1] != "/" {
		cwd += "/"
	}
	defer func() {
		newErr := os.Chdir(cwd)
		if newErr != nil {
			err = newErr
		}
	}()

	err = p.preRunPlugins(baseConfigPath)
	if err != nil {
		return err
	}

	err = p.createFilesAndFolders(baseConfigPath)
	if err != nil {
		return err
	}

	return p.postRunPlugins(baseConfigPath)
}

// createProjectFolder tries to create the main project folder.
func (p *Project) createProjectFolder() error {
	return os.Mkdir(p.Path, os.ModePerm)
}

func (p *Project) createFilesAndFolders(baseConfigPath string) error {
	baseTemplatesPath := filepath.Join(baseConfigPath, "/templates/")
	for _, template := range p.Class.Templates {
		if len(template.Path) > 0 {
			// Copy template file or folder
			err := copy.Copy(filepath.Join(baseTemplatesPath, template.Path), template.Destination)
			if err != nil {
				return err
			}
		}
		if template.IsFile {
			// Create file
			_, err := os.Create(template.Destination)
			if err != nil {
				return err
			}
		} else {
			// Create folder
			err := os.MkdirAll(template.Destination, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Project) preRunPlugins(baseConfigPath string) error {
	basePluginsPath := filepath.Join(baseConfigPath, "plugins")
	for _, plugin := range p.Class.Plugins {
		if plugin.ExecNumber >= 0 {
			continue
		}
		// Plugin path is relative by default to make it shareable. We have to make it an absolute path here,
		// so that we can execute it.
		p.Path = filepath.Join(basePluginsPath, p.Path)
		err := plugin.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) postRunPlugins(baseConfigPath string) error {
	basePluginsPath := filepath.Join(baseConfigPath, "plugins")
	for _, plugin := range p.Class.Plugins {
		if plugin.ExecNumber <= 0 {
			continue
		}
		// Plugin path is relative by default to make it shareable. We have to make it an absolute path here,
		// so that we can execute it.
		p.Path = filepath.Join(basePluginsPath, p.Path)
		err := plugin.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
