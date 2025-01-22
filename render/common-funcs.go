package render

import (
	"os"
	"path/filepath"
	// "fmt"
)

// Read the file.
func readFile(file string) []byte {
	data, err := os.ReadFile(file)
	if err != nil {
		errLog.Fatalf("ERROR: \"%v\"}", err)
	}
	return data
}

// Recursively merge right Values into left one
func mergeKeys(left, right Values) Values {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			if _, ok := leftVal.(Values); ok {
				left[key] = mergeKeys(leftVal.(Values), rightVal.(Values))
			} else {
				left[key] = rightVal
			}
		} else {
			if left == nil {
				left = make(Values)
			}
			left[key] = rightVal
		}
	}
	return left
}

// If glob is used. Then map folder tree as yaml structure.
func DirrectoryMapping(dir []string, val Values) Values {
	data := make(Values)
	if len(dir) > 1 {
		data[dir[0]] = DirrectoryMapping(dir[1:], val)
	} else if dir[0] != "." {
		data[dir[0]] = val
	} else {
		data = val
	}
	return data
}

// Make file path absolut depend on main file.
func absolutePath(basefile, depenfile string) string {
	dir := filepath.Dir(basefile)
	return filepath.Join(dir, depenfile)
}
