package pkg

import (
	"errors"
	"fmt"
	"html/template"
	"os"
)

// This method returns either the template with nil or nothing with the error
// based on whether the path provided to it actually exists or not
func GetTemplate(path string) (*template.Template, error) {
	// Check if the path exists or not
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf(
			"fatal: %s, %s",
			err.Error(),
			errors.New("couldn't find the template").Error(),
		)
	}

	// Parse the template
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf(
			"fatal: %s, %s",
			err.Error(),
			errors.New("couldn't parse the template").Error(),
		)
	}

	return tmpl, nil // Return the parsed template
}
