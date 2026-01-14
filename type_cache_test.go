package mongoose_test

import (
	"reflect"
	"testing"

	"github.com/tinh-tinh/mongoose/v2"
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

// BenchmarkModel for testing cached type info performance
type BenchmarkModel struct {
	mongoose.BaseSchema `bson:"inline"`
	Title               string `bson:"title"`
	Author              string `bson:"author"`
	Content             string `bson:"content"`
	Category            string `bson:"category"`
	Tags                string `bson:"tags"`
}

func (b BenchmarkModel) CollectionName() string {
	return "benchmarks"
}

// BenchmarkGetCollectionName_WithCache benchmarks collection name lookup with caching
func BenchmarkGetCollectionName_WithCache(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = mongoose.GetCachedCollectionName[BenchmarkModel]()
	}
}

// BenchmarkGetCollectionName_NoCache benchmarks collection name lookup without caching
func BenchmarkGetCollectionName_NoCache(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var model BenchmarkModel
		v := reflect.ValueOf(&model).Elem()
		fnc := v.MethodByName("CollectionName")
		if fnc.IsValid() {
			_ = fnc.Call(nil)[0].String()
		} else {
			_ = common.GetStructName(model)
		}
	}
}

// BenchmarkGetTypeInfo_WithCache benchmarks type info lookup with caching
func BenchmarkGetTypeInfo_WithCache(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = mongoose.GetTypeInfo[BenchmarkModel]()
	}
}

// BenchmarkGetTypeInfo_NoCache benchmarks type info computation without caching
func BenchmarkGetTypeInfo_NoCache(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var model BenchmarkModel
		t := reflect.TypeOf(model)
		// Simulate the type info computation
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			_ = field.Name
			_ = field.Tag.Get("bson")
			_ = field.Tag.Get("mongoose")
			_ = field.Type.Name()
		}
	}
}

// BenchmarkFieldIteration_Cached benchmarks field iteration with cached type info
func BenchmarkFieldIteration_Cached(b *testing.B) {
	b.ReportAllocs()
	typeInfo := mongoose.GetTypeInfo[BenchmarkModel]()
	for i := 0; i < b.N; i++ {
		for _, field := range typeInfo.Fields {
			_ = field.Name
			_ = field.BsonTag
		}
	}
}

// BenchmarkFieldIteration_Reflect benchmarks field iteration with reflection
func BenchmarkFieldIteration_Reflect(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var model BenchmarkModel
		t := reflect.TypeOf(model)
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			_ = field.Name
			_ = field.Tag.Get("bson")
		}
	}
}

// ModelWithRef for testing ref path caching
type ModelWithRef struct {
	mongoose.BaseSchema `bson:"inline"`
	Name                string `bson:"name"`
	AuthorID            string `bson:"authorId" ref:"authorId->authors"`
}

func (m ModelWithRef) CollectionName() string {
	return "models_with_ref"
}

// ModelWithInvalidRef has ref tag but no bson tag
type ModelWithInvalidRef struct {
	Name     string `bson:"name"`
	AuthorID string `ref:"authorId->authors"` // Missing bson tag
}

func (m ModelWithInvalidRef) CollectionName() string {
	return "models_invalid_ref"
}

// ModelNoRef has no ref tags
type ModelNoRef struct {
	mongoose.BaseSchema `bson:"inline"`
	Name                string `bson:"name"`
}

func (m ModelNoRef) CollectionName() string {
	return "models_no_ref"
}

// TestRefPath_ValidRef tests that valid ref paths are cached correctly
func TestRefPath_ValidRef(t *testing.T) {
	typeInfo := mongoose.GetTypeInfo[ModelWithRef]()

	refPath := typeInfo.RefPaths["authorId"]
	if refPath == nil {
		t.Fatal("Expected refPath to be non-nil for valid ref tag")
	}

	if refPath.From != "authors" {
		t.Errorf("Expected From='authors', got '%s'", refPath.From)
	}
	if refPath.ForeignKey != "authorId" {
		t.Errorf("Expected ForeignKey='authorId', got '%s'", refPath.ForeignKey)
	}
	if refPath.As != "authorId" {
		t.Errorf("Expected As='authorId', got '%s'", refPath.As)
	}
}

// TestRefPath_InvalidRefName tests that invalid ref names return nil
func TestRefPath_InvalidRefName(t *testing.T) {
	typeInfo := mongoose.GetTypeInfo[ModelWithRef]()

	// Try to get a ref path that doesn't exist
	refPath := typeInfo.RefPaths["nonexistent"]
	if refPath != nil {
		t.Error("Expected nil refPath for non-existent ref name")
	}
}

// TestRefPath_MissingBsonTag tests that ref tags without bson tags are ignored
func TestRefPath_MissingBsonTag(t *testing.T) {
	typeInfo := mongoose.GetTypeInfo[ModelWithInvalidRef]()

	// AuthorID has ref tag but no bson tag, should not be in RefPaths
	refPath := typeInfo.RefPaths["authorId"]
	if refPath != nil {
		t.Error("Expected nil refPath when bson tag is missing")
	}
}

// TestRefPath_NoRefTags tests model with no ref tags
func TestRefPath_NoRefTags(t *testing.T) {
	typeInfo := mongoose.GetTypeInfo[ModelNoRef]()

	if len(typeInfo.RefPaths) != 0 {
		t.Errorf("Expected empty RefPaths, got %d entries", len(typeInfo.RefPaths))
	}
}
