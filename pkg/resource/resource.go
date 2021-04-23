package resource

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type Resource interface {
	TerraformId() string
	TerraformType() string
	CtyValue() *cty.Value
}

var refactoredResources = []string{
	"aws_cloudfront_distribution",
}

func IsRefactoredResource(typ string) bool {
	for _, refactoredResource := range refactoredResources {
		if typ == refactoredResource {
			return true
		}
	}
	return false
}

type AbstractResource struct {
	Id    string
	Type  string
	Attrs *ResourceAttributes
}

func (a *AbstractResource) TerraformId() string {
	return a.Id
}

func (a *AbstractResource) TerraformType() string {
	return a.Type
}

func (a *AbstractResource) CtyValue() *cty.Value {
	return nil
}

type ResourceFactory interface {
	CreateResource(data interface{}, ty string) (*cty.Value, error)
}

type SerializableResource struct {
	Resource
}

type SerializedResource struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func (u SerializedResource) TerraformId() string {
	return u.Id
}

func (u SerializedResource) TerraformType() string {
	return u.Type
}

func (u SerializedResource) CtyValue() *cty.Value {
	return &cty.NilVal
}

func (s *SerializableResource) UnmarshalJSON(bytes []byte) error {
	var res SerializedResource

	if err := json.Unmarshal(bytes, &res); err != nil {
		return err
	}
	s.Resource = res
	return nil
}

func (s SerializableResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(SerializedResource{Id: s.TerraformId(), Type: s.TerraformType()})
}

type NormalizedResource interface {
	NormalizeForState() (Resource, error)
	NormalizeForProvider() (Resource, error)
}

func IsSameResource(rRs, lRs Resource) bool {
	return rRs.TerraformType() == lRs.TerraformType() && rRs.TerraformId() == lRs.TerraformId()
}

func Sort(res []Resource) []Resource {
	sort.SliceStable(res, func(i, j int) bool {
		if res[i].TerraformType() != res[j].TerraformType() {
			return res[i].TerraformType() < res[j].TerraformType()
		}
		return res[i].TerraformId() < res[j].TerraformId()
	})
	return res
}

func ToResourceAttributes(val *cty.Value) *ResourceAttributes {
	if val == nil {
		return nil
	}

	bytes, _ := ctyjson.Marshal(*val, val.Type())
	var attrs map[string]interface{}
	err := json.Unmarshal(bytes, &attrs)
	if err != nil {
		panic(err)
	}

	return &ResourceAttributes{
		attrs,
		val.Type(),
	}
}

type ResourceAttributes struct {
	Attrs map[string]interface{}
	Type  cty.Type
}

func (a *ResourceAttributes) Get(path string) (interface{}, bool) {
	val, exist := a.Attrs[path]
	return val, exist
}

func (a *ResourceAttributes) SafeDelete(path []string) {
	val := a.Attrs
	for i, key := range path {
		if i == len(path)-1 {
			delete(val, key)
			return
		}

		v, exists := val[key]
		if !exists {
			return
		}
		m, ok := v.(map[string]interface{})
		if !ok {
			return
		}
		val = m
	}
}

func (a *ResourceAttributes) SafeSet(path []string, value interface{}) error {
	val := a.Attrs
	for i, key := range path {
		if i == len(path)-1 {
			val[key] = value
			return nil
		}

		v, exists := val[key]
		if !exists {
			val[key] = map[string]interface{}{}
			v = val[key]
		}

		m, ok := v.(map[string]interface{})
		if !ok {
			return errors.Errorf("Path %s cannot be set: %s is not a nested struct", strings.Join(path, "."), key)
		}
		val = m
	}
	return errors.New("Error setting value") // should not happen ?
}

func (a *ResourceAttributes) SetDefault(path []string) error {
	ctyVal, err := gocty.ToCtyValue(a.Attrs, a.Type)
	if err != nil {
		return err
	}
	attrVal := ctyVal.GetAttr(path[0])
	if attrVal.IsNull() || (!attrVal.IsNull() && attrVal.LengthInt() == 0) {
		a.SafeDelete(path)
	}
	return nil
}

func (a *ResourceAttributes) SanitizeDefaults() {
	original := reflect.ValueOf(a.Attrs)
	copy := reflect.New(original.Type()).Elem()
	a.run("", original, copy)
	a.Attrs = copy.Interface().(map[string]interface{})
}

func (a *ResourceAttributes) run(path string, original, copy reflect.Value) {
	switch original.Kind() {
	case reflect.Ptr:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return
		}
		copy.Set(reflect.New(originalValue.Type()))
		a.run(path, originalValue, copy.Elem())
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		if originalValue.Len() == 0 {
			fmt.Printf("Skipped empty value %s\n", path)
			return
		}
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		a.run(path, originalValue, copyValue)
		copy.Set(copyValue)

	case reflect.Struct:
		fmt.Printf("Reading struct field %s\n", path)
		for i := 0; i < original.NumField(); i += 1 {
			field := original.Field(i)
			a.run(concatenatePath(path, field.String()), field, copy.Field(i))
		}
	case reflect.Slice:
		fmt.Printf("Reading slice field %s\n", path)
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			a.run(concatenatePath(path, strconv.Itoa(i)), original.Index(i), copy.Index(i))
		}
	case reflect.Map:
		fmt.Printf("Reading map field %s\n", path)
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			a.run(concatenatePath(path, key.String()), originalValue, copyValue)
			copy.SetMapIndex(key, copyValue)
		}
	default:
		fmt.Printf("Reading leaf field %s\n", path)
		copy.Set(original)
	}
}

func concatenatePath(path, next string) string {
	if path == "" {
		return next
	}
	return strings.Join([]string{path, next}, ".")
}

func (a *ResourceAttributes) SanitizeDefaultsV2() {
	original := reflect.ValueOf(a.Attrs)
	copy := reflect.New(original.Type()).Elem()
	a.runV2("", original, copy)
	a.Attrs = copy.Interface().(map[string]interface{})
}

func (a *ResourceAttributes) runV2(path string, original, copy reflect.Value) {
	switch original.Kind() {
	case reflect.Ptr:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return
		}
		copy.Set(reflect.New(originalValue.Type()))
		a.runV2(path, originalValue, copy.Elem())
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		a.runV2(path, originalValue, copyValue)
		copy.Set(copyValue)

	case reflect.Struct:
		fmt.Printf("Reading struct field %s\n", path)
		for i := 0; i < original.NumField(); i += 1 {
			field := original.Field(i)
			a.runV2(concatenatePath(path, field.String()), field, copy.Field(i))
		}
	case reflect.Slice:
		fmt.Printf("Reading slice field %s\n", path)
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			a.runV2(concatenatePath(path, strconv.Itoa(i)), original.Index(i), copy.Index(i))
		}
	case reflect.Map:
		fmt.Printf("Reading map field %s\n", path)
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			a.runV2(concatenatePath(path, key.String()), originalValue, copyValue)
			copy.SetMapIndex(key, copyValue)
		}
	default:
		fmt.Printf("Reading leaf field %s\n", path)
		copy.Set(original)
	}
}

func (a *ResourceAttributes) SanitizeDefaultsV3() {
	original := reflect.ValueOf(a.Attrs)
	copy := reflect.New(original.Type()).Elem()
	a.runV3("", original, copy)
	a.Attrs = copy.Interface().(map[string]interface{})
}

func (a *ResourceAttributes) runV3(path string, original, copy reflect.Value) bool {
	switch original.Kind() {
	case reflect.Ptr:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return false
		}
		copy.Set(reflect.New(originalValue.Type()))
		a.runV3(path, originalValue, copy.Elem())
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		if originalValue.Len() == 0 {
			fmt.Printf("Skipped empty value %s\n", path)
			return false
		}
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		a.runV3(path, originalValue, copyValue)
		copy.Set(copyValue)

	case reflect.Struct:
		fmt.Printf("Reading struct field %s\n", path)
		for i := 0; i < original.NumField(); i += 1 {
			field := original.Field(i)
			a.runV3(concatenatePath(path, field.String()), field, copy.Field(i))
		}
	case reflect.Slice:
		fmt.Printf("Reading slice field %s\n", path)
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			a.runV3(concatenatePath(path, strconv.Itoa(i)), original.Index(i), copy.Index(i))
		}
	case reflect.Map:
		fmt.Printf("Reading map field %s\n", path)
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			created := a.runV3(concatenatePath(path, key.String()), originalValue, copyValue)
			if created {
				copy.SetMapIndex(key, copyValue)
			}
		}
	default:
		fmt.Printf("Reading leaf field %s\n", path)
		copy.Set(original)
	}
	return true
}
