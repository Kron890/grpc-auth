package envcfg

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Load заполняет структуру значениями из переменных окружения согласно тегам `env` и `default`.
// Пример тега: `env:"SERVER_PORT,required"` или с значением по умолчанию: `default:"8080"`.
func Load(target any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("target must be a non-nil pointer to a struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("target must point to a struct")
	}

	rt := rv.Type()
	var missing []string

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		parts := strings.Split(envTag, ",")
		envName := strings.TrimSpace(parts[0])
		required := false
		for _, p := range parts[1:] {
			if strings.TrimSpace(p) == "required" {
				required = true
			}
		}

		val, ok := os.LookupEnv(envName)
		if !ok || val == "" {
			def := field.Tag.Get("default")
			if def != "" {
				val = def
				ok = true
			}
		}

		if (!ok || val == "") && required {
			missing = append(missing, envName)
			continue
		}

		if ok {
			fv := rv.Field(i)
			if !fv.CanSet() {
				continue
			}
			if err := setWithString(fv, val); err != nil {
				return fmt.Errorf("failed to set field %s from %s: %w", field.Name, envName, err)
			}
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return nil
}

func setWithString(fv reflect.Value, s string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(s)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(n)
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		fv.SetBool(b)
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(f)
		return nil
	default:
		return fmt.Errorf("unsupported kind: %s", fv.Kind())
	}
}
