package mongoose

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFilter_DetectsDangerousOperators(t *testing.T) {
	testCases := []struct {
		name     string
		filter   interface{}
		expected string // expected operator in error, empty if no error expected
	}{
		{
			name:     "simple $ne operator",
			filter:   map[string]interface{}{"password": map[string]interface{}{"$ne": ""}},
			expected: "$ne",
		},
		{
			name:     "simple $gt operator",
			filter:   map[string]interface{}{"age": map[string]interface{}{"$gt": 18}},
			expected: "$gt",
		},
		{
			name:     "$where operator (very dangerous)",
			filter:   map[string]interface{}{"$where": "this.password == ''"},
			expected: "$where",
		},
		{
			name:     "$or operator",
			filter:   map[string]interface{}{"$or": []interface{}{map[string]interface{}{"admin": true}}},
			expected: "$or",
		},
		{
			name:     "nested $in operator",
			filter:   map[string]interface{}{"role": map[string]interface{}{"$in": []string{"admin", "root"}}},
			expected: "$in",
		},
		{
			name:     "$regex operator",
			filter:   map[string]interface{}{"email": map[string]interface{}{"$regex": ".*"}},
			expected: "$regex",
		},
		{
			name:     "$expr operator",
			filter:   map[string]interface{}{"$expr": map[string]interface{}{"$eq": []string{"$password", ""}}},
			expected: "$expr",
		},
		{
			name:     "deeply nested dangerous operator",
			filter:   map[string]interface{}{"user": map[string]interface{}{"profile": map[string]interface{}{"age": map[string]interface{}{"$gte": 0}}}},
			expected: "$gte",
		},
		{
			name: "struct with embedded map containing operator (bypass attempt)",
			filter: struct {
				Query map[string]interface{}
			}{
				Query: map[string]interface{}{"$where": "malicious"},
			},
			expected: "$where",
		},
		{
			name: "$function operator (CVE-2025-10061)",
			filter: map[string]interface{}{"$function": map[string]interface{}{
				"body": "function() { return true }",
				"args": []interface{}{},
				"lang": "js",
			}},
			expected: "$function",
		},
		{
			name:     "$accumulator operator (server-side JS)",
			filter:   map[string]interface{}{"$accumulator": "malicious"},
			expected: "$accumulator",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SanitizeFilter(tc.filter)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expected)
			assert.True(t, IsDangerousOperatorError(err))
		})
	}
}

func TestSanitizeFilter_AllowsSafeFilters(t *testing.T) {
	testCases := []struct {
		name   string
		filter interface{}
	}{
		{
			name:   "nil filter",
			filter: nil,
		},
		{
			name:   "simple string value",
			filter: map[string]interface{}{"name": "john"},
		},
		{
			name:   "multiple safe fields",
			filter: map[string]interface{}{"name": "john", "age": 25, "active": true},
		},
		{
			name:   "nested safe map",
			filter: map[string]interface{}{"user": map[string]interface{}{"name": "john", "email": "test@test.com"}},
		},
		{
			name:   "struct-based filter",
			filter: struct{ Name string }{Name: "john"},
		},
		{
			name:   "pointer to struct filter",
			filter: &struct{ Name string }{Name: "john"},
		},
		{
			name:   "safe field with dollar sign in value (not key)",
			filter: map[string]interface{}{"price": "$100"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SanitizeFilter(tc.filter)
			assert.Nil(t, err)
		})
	}
}

func TestSanitizeFilter_AllDangerousOperators(t *testing.T) {
	// Test that all operators in the DangerousOperators list are detected
	for _, op := range DangerousOperators {
		t.Run(op, func(t *testing.T) {
			filter := map[string]interface{}{op: "test"}
			err := SanitizeFilter(filter)
			require.Error(t, err)
			assert.True(t, IsDangerousOperatorError(err))

			var opErr *ErrDangerousOperator
			require.ErrorAs(t, err, &opErr)
			assert.Equal(t, op, opErr.Operator)
		})
	}
}

func TestIsDangerousOperator(t *testing.T) {
	testCases := []struct {
		key      string
		expected bool
	}{
		{"$ne", true},
		{"$gt", true},
		{"$where", true},
		{"name", false},
		{"$custom", false}, // Not in our list
		{"", false},
		{"$", false}, // Just dollar sign, not in list
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			result := IsDangerousOperator(tc.key)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestErrDangerousOperator_Error(t *testing.T) {
	err := &ErrDangerousOperator{Operator: "$ne"}
	assert.Equal(t, "dangerous MongoDB operator detected in filter: $ne", err.Error())
}

func TestSanitizeFilter_SliceWithDangerousOperator(t *testing.T) {
	// Test detection in array elements
	filter := []interface{}{
		map[string]interface{}{"name": "safe"},
		map[string]interface{}{"password": map[string]interface{}{"$ne": ""}},
	}
	err := SanitizeFilter(filter)
	require.Error(t, err)
	assert.True(t, IsDangerousOperatorError(err))
}

// TestAuthBypassAttempt simulates a real authentication bypass attack
func TestAuthBypassAttempt(t *testing.T) {
	// This simulates what an attacker might try to send via JSON body
	// POST /login with body: {"username": "admin", "password": {"$ne": ""}}
	maliciousFilter := map[string]interface{}{
		"username": "admin",
		"password": map[string]interface{}{"$ne": ""},
	}

	err := SanitizeFilter(maliciousFilter)
	require.Error(t, err)

	var opErr *ErrDangerousOperator
	require.ErrorAs(t, err, &opErr)
	assert.Equal(t, "$ne", opErr.Operator)
}

// Test Model.sanitizeFilter helper method
func TestModel_SanitizeFilter(t *testing.T) {
	// Test with StrictFilters disabled (default)
	modelDisabled := NewModel[testStruct]()
	err := modelDisabled.sanitizeFilter(map[string]interface{}{"$ne": "test"})
	assert.Nil(t, err, "should allow dangerous operators when StrictFilters is disabled")

	// Test with StrictFilters enabled
	modelEnabled := NewModel[testStruct](ModelOptions{StrictFilters: true})
	err = modelEnabled.sanitizeFilter(map[string]interface{}{"name": map[string]interface{}{"$ne": ""}})
	assert.Error(t, err, "should reject dangerous operators when StrictFilters is enabled")
	assert.True(t, IsDangerousOperatorError(err))

	// Test with safe filter when StrictFilters is enabled
	err = modelEnabled.sanitizeFilter(map[string]interface{}{"name": "safe"})
	assert.Nil(t, err, "should allow safe filters when StrictFilters is enabled")

	// Test with nil filter
	err = modelEnabled.sanitizeFilter(nil)
	assert.Nil(t, err, "should allow nil filters")

	// Test with struct filter (always safe)
	err = modelEnabled.sanitizeFilter(&testStruct{Name: "test"})
	assert.Nil(t, err, "should allow struct-based filters")
}

type testStruct struct {
	BaseSchema `bson:"inline"`
	Name       string `bson:"name"`
}

func (t testStruct) CollectionName() string {
	return "test_sanitize"
}

// TestStrictFilters_AllOperators ensures all dangerous operators are blocked
func TestStrictFilters_AllOperators(t *testing.T) {
	model := NewModel[testStruct](ModelOptions{StrictFilters: true})

	for _, op := range DangerousOperators {
		t.Run(op, func(t *testing.T) {
			filter := map[string]interface{}{op: "value"}
			err := model.sanitizeFilter(filter)
			require.Error(t, err, "operator %s should be blocked", op)
			assert.True(t, IsDangerousOperatorError(err))
		})
	}
}

// TestStrictFilters_NestedOperators tests deeply nested dangerous operators
func TestStrictFilters_NestedOperators(t *testing.T) {
	model := NewModel[testStruct](ModelOptions{StrictFilters: true})

	testCases := []struct {
		name   string
		filter interface{}
	}{
		{
			name:   "nested in map value",
			filter: map[string]interface{}{"field": map[string]interface{}{"$gt": 10}},
		},
		{
			name:   "deeply nested",
			filter: map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"$ne": ""}}},
		},
		{
			name:   "in slice element",
			filter: []interface{}{map[string]interface{}{"$or": []interface{}{}}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := model.sanitizeFilter(tc.filter)
			require.Error(t, err)
			assert.True(t, IsDangerousOperatorError(err))
		})
	}
}

// TestSanitizeFilter_EdgeCases tests edge cases for the sanitization logic
func TestSanitizeFilter_EdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		filter      interface{}
		shouldError bool
	}{
		{
			name:        "empty map",
			filter:      map[string]interface{}{},
			shouldError: false,
		},
		{
			name:        "empty slice",
			filter:      []interface{}{},
			shouldError: false,
		},
		{
			name:        "nil pointer",
			filter:      (*map[string]interface{})(nil),
			shouldError: false,
		},
		{
			name:        "primitive string",
			filter:      "simple string",
			shouldError: false,
		},
		{
			name:        "primitive int",
			filter:      42,
			shouldError: false,
		},
		{
			name:        "key starts with $ but not in list",
			filter:      map[string]interface{}{"$unknownOp": "value"},
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SanitizeFilter(tc.filter)
			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
