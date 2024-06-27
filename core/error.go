package core

import "fmt"

func WrapConstraintError(err error) error {
	return fmt.Errorf("ent constraint error: %w", err)
}
