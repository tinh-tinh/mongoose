package mongoose

import (
	"reflect"
	"sync"

	"github.com/tinh-tinh/tinhtinh/v2/common"
)

// FieldInfo caches reflection metadata for a single struct field
type FieldInfo struct {
	Index       int    // Field index in struct (direct index for top-level fields)
	Name        string // Go field name
	BsonTag     string // bson tag value
	MongooseTag string // mongoose tag value (e.g., "readonly")
	TypeName    string // Field type name (e.g., "BaseSchema")
	RefTag      string // ref tag value for population
	IndexPath   []int  // Full index path for nested fields (e.g., [0, 1] for embedded)
}

// TypeInfo caches reflection metadata for a struct type
type TypeInfo struct {
	CollectionName string                // Cached result of CollectionName() method
	Fields         []FieldInfo           // Ordered list of struct fields
	FieldsByName   map[string]*FieldInfo // Lookup by field name (includes promoted fields)
	FieldsByBson   map[string]*FieldInfo // Lookup by bson tag
	RefPaths       map[string]*RefPath   // Cached ref paths by foreign key
}

// TypeCache is a thread-safe cache for type metadata
type TypeCache struct {
	mu    sync.RWMutex
	cache map[reflect.Type]*TypeInfo
}

// Global type cache instance
var globalTypeCache = &TypeCache{
	cache: make(map[reflect.Type]*TypeInfo),
}

// GetTypeInfo returns cached TypeInfo for a given type, computing it if not cached
func GetTypeInfo[M any]() *TypeInfo {
	var m M
	t := reflect.TypeOf(m)

	// Fast path: check with read lock
	globalTypeCache.mu.RLock()
	info, ok := globalTypeCache.cache[t]
	globalTypeCache.mu.RUnlock()

	if ok {
		return info
	}

	// Slow path: compute and cache with write lock
	globalTypeCache.mu.Lock()
	defer globalTypeCache.mu.Unlock()

	// Double-check after acquiring write lock
	if info, ok = globalTypeCache.cache[t]; ok {
		return info
	}

	info = computeTypeInfo[M]()
	globalTypeCache.cache[t] = info
	return info
}

// computeTypeInfo computes TypeInfo for a type using reflection
func computeTypeInfo[M any]() *TypeInfo {
	var model M
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(&model).Elem()

	info := &TypeInfo{
		Fields:       make([]FieldInfo, 0),
		FieldsByName: make(map[string]*FieldInfo),
		FieldsByBson: make(map[string]*FieldInfo),
		RefPaths:     make(map[string]*RefPath),
	}

	// Get collection name via CollectionName method or struct name
	fnc := v.MethodByName("CollectionName")
	if fnc.IsValid() {
		info.CollectionName = fnc.Call(nil)[0].String()
	} else {
		info.CollectionName = common.GetStructName(model)
	}

	// Phase 1: Collect all fields recursively (slice may reallocate)
	collectFieldsRecursive(t, info, []int{})

	// Phase 2: Build maps after slice is stable (no more appends)
	for i := range info.Fields {
		field := &info.Fields[i]

		// Add to FieldsByName (first occurrence wins for promoted fields)
		if _, exists := info.FieldsByName[field.Name]; !exists {
			info.FieldsByName[field.Name] = field
		}

		// Add to FieldsByBson
		if field.BsonTag != "" {
			if _, exists := info.FieldsByBson[field.BsonTag]; !exists {
				info.FieldsByBson[field.BsonTag] = field
			}
		}

		// Parse ref tags for population
		if field.RefTag != "" {
			refPath := parseRefPath(*field)
			if refPath != nil {
				info.RefPaths[refPath.ForeignKey] = refPath
			}
		}
	}

	return info
}

// collectFieldsRecursive collects fields including promoted fields from embedded structs
func collectFieldsRecursive(t reflect.Type, info *TypeInfo, indexPath []int) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		currentIndex := append(append([]int{}, indexPath...), i)

		fieldInfo := FieldInfo{
			Index:       i,
			Name:        field.Name,
			BsonTag:     field.Tag.Get("bson"),
			MongooseTag: field.Tag.Get("mongoose"),
			TypeName:    field.Type.Name(),
			RefTag:      field.Tag.Get("ref"),
			IndexPath:   currentIndex,
		}
		info.Fields = append(info.Fields, fieldInfo)

		// Recursively collect promoted fields from embedded structs
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			collectFieldsRecursive(field.Type, info, currentIndex)
		}
	}
}

// parseRefPath parses a ref tag into a RefPath struct
func parseRefPath(field FieldInfo) *RefPath {
	if field.RefTag == "" {
		return nil
	}

	// BsonTag is required for the 'as' field in $lookup aggregation
	if field.BsonTag == "" {
		return nil
	}

	// ref tag format: "foreignKey->collectionName"
	var foreignKey, foreignCol string
	for i, c := range field.RefTag {
		if c == '-' && i+1 < len(field.RefTag) && field.RefTag[i+1] == '>' {
			foreignKey = field.RefTag[:i]
			foreignCol = field.RefTag[i+2:]
			break
		}
	}

	if foreignKey == "" {
		return nil
	}

	return &RefPath{
		From:       foreignCol,
		ForeignKey: foreignKey,
		As:         field.BsonTag,
	}
}

// GetCachedCollectionName returns the cached collection name for type M
func GetCachedCollectionName[M any]() string {
	return GetTypeInfo[M]().CollectionName
}
