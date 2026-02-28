package main

import (
	"fmt"
	"net/http"
)

// feedTemplate loads and executes a template found under label using
// the given data parameter (used by [template.Template.Execute] to
// fill in the template.)
func feedTemplate(w http.ResponseWriter, label string, data any) error {
	// Load the template.
	t, ok := tpls[label+".gohtml"]
	if !ok {
		return fmt.Errorf("no template under label '%s'", label)
	}

	// Use the template.
	if err := t.Execute(w, data); err != nil {
		return fmt.Errorf("can't execute template under label '%s': %w", label, err)
	}

	return nil
}
