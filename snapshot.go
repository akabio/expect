package expect

import (
	"os"
	"path/filepath"

	"golang.org/x/exp/slices"
)

func (e Val) ToBeSnapshot(path string) Val {
	e.t.Helper()

	folder := filepath.Dir(path)
	if folder != "." {
		err := os.MkdirAll(folder, 0o755)
		if err != nil {
			e.t.Fatalf("failed to create target folder %v", folder)
		}
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		e.t.Fatalf("failed to read snaphsot %v: %v", path, err)
	}

	currentString, isString := e.value.(string)
	current, isBytes := e.value.([]byte)
	if !isString && !isBytes {
		e.t.Fatalf("value of .ToBeSnaphsot must be of type string or []byte but it is %T", e.value)
	}

	if isString {
		current = []byte(currentString)
	}

	if existing == nil {
		// snapshot does not exist, create it
		err = os.WriteFile(path, current, 0o644)
		if err != nil {
			e.t.Fatalf("failed to write snapshot %v", path)
		}
	} else {
		if slices.Equal(current, existing) {
			// all is well, snapshot is matched, remove a possible current version
			os.RemoveAll(path + ".current")
		} else {
			e.t.Errorf("snapshot for %v does not match current output", path)
			err = os.WriteFile(path+".current", current, 0o644)
			if err != nil {
				e.t.Fatalf("failed to write snapshot %v", path)
			}
		}
	}

	return e
}
