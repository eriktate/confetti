package confetti

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// FromEnv returns a type T hydrated by the environment using [ApplyEnv].
func FromEnv[T any]() (T, error) {
	var target T
	return target, ApplyEnv(&target)
}

// FromFiles returns a type T hydrated by the files at the given files using
// [ApplyFiles].
func FromFiles[T any](paths ...string) (T, error) {
	var target T
	return target, ApplyFiles(&target, paths...)
}

// ApplyEnv attempts to coerce matching environment variables into struct fields. It
// matches using the `conf` struct field tag if present, falling back to the struct
// field name otherwise.
func ApplyEnv(target any) error {
	targetType, targetVal, err := getTarget(target)
	if err != nil {
		return err
	}

	targetName := targetType.Name()
	for i := range targetType.NumField() {
		field := targetType.Field(i)

		confKey := field.Tag.Get("conf")
		if confKey == "" {
			confKey = field.Name
		}

		val := os.Getenv(confKey)
		if val == "" {
			continue
		}

		if err := coerceValue(field, targetVal.Field(i), val); err != nil {
			return fmt.Errorf("applying env to %q: %w", targetName, err)
		}
	}

	return nil
}

// ApplyFiles reads .env formatted files and attempts to apply them to the given target.
// Files are applied in order with the latter taking precedence. It matches on keys using
// the `conf` struct field tag if present, falling back to the struct field name
// otherwise.
func ApplyFiles(target any, paths ...string) error {
	for _, path := range paths {
		if err := applyFile(target, path); err != nil {
			return err
		}
	}

	return nil
}

func applyFile(target any, path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	var done bool
	for !done {
		line, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("reading config file: %w", err)
			}

			done = true
		}

		key, val, found := strings.Cut(string(line), "=")
		if !found {
			// skip lines with bogus config values
			continue
		}

		if err := applyKeyVal(
			target,
			strings.Trim(key, " \t\n"),
			strings.Trim(val, " \t\n"),
		); err != nil {
			return fmt.Errorf("applying %q: %w", path, err)
		}
	}

	return nil
}

func getTarget(target any) (reflect.Type, reflect.Value, error) {
	ptrType := reflect.TypeOf(target)
	if ptrType.Kind() != reflect.Pointer {
		return nil,
			reflect.Value{},
			errors.New("confetti can only parse into pointer types")
	}

	targetType := ptrType.Elem()
	if targetType.Kind() != reflect.Struct {
		return nil,
			reflect.Value{},
			errors.New("confetti can only parse into struct types")
	}

	return targetType, reflect.ValueOf(target).Elem(), nil
}

func coerceValue(field reflect.StructField, val reflect.Value, str string) error {
	switch field.Type.Kind() {
	case reflect.String:
		val.SetString(str)
	case reflect.Bool:
		switch strings.ToLower(str) {
		case "true", "t", "yes", "1", "on":
			val.SetBool(true)
		case "", "false", "f", "no", "0", "off":
			val.SetBool(false)
		default:
			return fmt.Errorf("could not assign %q to bool %q", str, field.Name)
		}
	case reflect.Int:
		intVal, err := strconv.Atoi(str)
		if err != nil {
			return fmt.Errorf("could not assign %q to int %q: %w", str, field.Name, err)
		}
		val.SetInt(int64(intVal))
	case reflect.Uint:
		uintVal, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return fmt.Errorf("could not assign %q to uint %q: %w", str, field.Name, err)
		}
		val.SetUint(uint64(uintVal))
	case reflect.Slice:
		if field.Type.Elem().Kind() == reflect.Uint8 {
			val.Set(reflect.ValueOf([]byte(str)))
			break
		}

		return fmt.Errorf(
			"could not assign %q to slice %q: only byte slices are supported",
			str,
			field.Name,
		)
	}

	return nil
}

func applyKeyVal(target any, key, value string) error {
	targetType, targetVal, err := getTarget(target)
	if err != nil {
		return err
	}

	targetName := targetType.Name()
	for i := range targetType.NumField() {
		field := targetType.Field(i)

		confKey := field.Tag.Get("conf")
		if confKey == "" {
			confKey = field.Name
		}

		if confKey == key {
			fieldVal := targetVal.Field(i)
			if err := coerceValue(field, fieldVal, value); err != nil {
				return fmt.Errorf("applying config to %q: %w", targetName, err)
			}
		}
	}

	return nil
}
