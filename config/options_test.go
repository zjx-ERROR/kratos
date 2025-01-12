package config

import (
	"reflect"
	"testing"
)

func TestDefaultDecoder(t *testing.T) {
	src := &KeyValue{
		Key:    "service",
		Value:  []byte("config"),
		Format: "",
	}
	target := make(map[string]interface{})
	err := defaultDecoder(src, target)
	if err != nil {
		t.Fatal("err is not nil")
	}
	if !reflect.DeepEqual(target, map[string]interface{}{"service": []byte("config")}) {
		t.Fatal(`target is not equal to map[string]interface{}{"service": "config"}`)
	}

	src = &KeyValue{
		Key:    "service.name.alias",
		Value:  []byte("2233"),
		Format: "",
	}
	target = make(map[string]interface{})
	err = defaultDecoder(src, target)
	if err != nil {
		t.Fatal("err is not nil")
	}
	if !reflect.DeepEqual(map[string]interface{}{
		"service": map[string]interface{}{
			"name": map[string]interface{}{
				"alias": []byte("2233"),
			},
		},
	}, target) {
		t.Fatal(`target is not equal to map[string]interface{}{"service": map[string]interface{}{"name": map[string]interface{}{"alias": []byte("2233")}}}`)
	}
}

func TestDefaultResolver(t *testing.T) {
	var (
		portString       = "8080"
		countInt         = 10
		rateFloat        = 0.9
		decimals         = 0.1314
		binary           = 0b111010
		minusBinary      = -0b111010
		octal            = 0o61
		minusOctal       = -0o61
		hexadecimal      = 0xF3B
		minusHexadecimal = -0xF3B
	)

	data := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"notexist":         "${NOTEXIST:100}",
				"port":             "${PORT:8081}",
				"count":            "${COUNT:0}",
				"enable":           "${ENABLE:false}",
				"rate":             "${RATE}",
				"empty":            "${EMPTY:foobar}",
				"url":              "${URL:http://example.com}",
				"decimals":         "${DECIMALS}",
				"binary":           "${BINARY}",
				"minusBinary":      "${MINUSBINARY}",
				"hexadecimal":      "${HEXADECIMAL}",
				"minusHexadecimal": "${MINUSHEXADECIMAL}",
				"octal":            "${OCTAL}",
				"minusOctal":       "${MINUSOCTAL}",
				"array": []interface{}{
					"${PORT}",
					map[string]interface{}{"foobar": "${NOTEXIST:8081}"},
				},
				"value1": "${test.value}",
				"value2": "$PORT",
				"value3": "abc${PORT}foo${COUNT}bar",
				"value4": "${foo${bar}}",
			},
		},
		"test": map[string]interface{}{
			"value": "foobar",
		},
		"PORT":             "8080",
		"COUNT":            "10",
		"ENABLE":           "true",
		"RATE":             "0.9",
		"EMPTY":            "",
		"DECIMALS":         ".1314",
		"BINARY":           "0b111010",
		"MINUSBINARY":      "-0b111010",
		"HEXADECIMAL":      "0xF3B",
		"MINUSHEXADECIMAL": "-0xF3B",
		"OCTAL":            "0o61",
		"MINUSOCTAL":       "-0o61",
	}

	tests := []struct {
		name   string
		path   string
		expect interface{}
	}{
		{
			name:   "test not exist int env with default",
			path:   "foo.bar.notexist",
			expect: 100,
		},
		{
			name:   "test string with default",
			path:   "foo.bar.port",
			expect: portString,
		},
		{
			name:   "test int with default",
			path:   "foo.bar.count",
			expect: countInt,
		},
		{
			name:   "test bool with default",
			path:   "foo.bar.enable",
			expect: true,
		},
		{
			name:   "test float without default",
			path:   "foo.bar.rate",
			expect: rateFloat,
		},
		{
			name:   "test empty value with default",
			path:   "foo.bar.empty",
			expect: "",
		},
		{
			name:   "test url with default",
			path:   "foo.bar.url",
			expect: "http://example.com",
		},
		{
			name:   "test array",
			path:   "foo.bar.array",
			expect: []interface{}{8080, map[string]interface{}{"foobar": 8081}},
		},
		{
			name:   "test ${test.value}",
			path:   "foo.bar.value1",
			expect: "foobar",
		},
		{
			name:   "test $PORT",
			path:   "foo.bar.value2",
			expect: "$PORT",
		},
		{
			name:   "test abc${PORT}foo${COUNT}bar",
			path:   "foo.bar.value3",
			expect: "abc8080foo10bar",
		},
		{
			name:   "test ${foo${bar}}",
			path:   "foo.bar.value4",
			expect: "}",
		},
		{
			name:   "test decimals",
			path:   "foo.bar.decimals",
			expect: decimals,
		},
		{
			name:   "test binary",
			path:   "foo.bar.binary",
			expect: binary,
		},
		{
			name:   "test minusBinary",
			path:   "foo.bar.minusBinary",
			expect: minusBinary,
		},
		{
			name:   "test hexadecimal",
			path:   "foo.bar.hexadecimal",
			expect: hexadecimal,
		},
		{
			name:   "test minusHexadecimal",
			path:   "foo.bar.minusHexadecimal",
			expect: minusHexadecimal,
		},
		{
			name:   "test octal",
			path:   "foo.bar.octal",
			expect: octal,
		},
		{
			name:   "test minusOctal",
			path:   "foo.bar.minusOctal",
			expect: minusOctal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := defaultResolver(data)
			if err != nil {
				t.Fatal(`err is not nil`)
			}
			rd := reader{
				values: data,
			}
			if v, ok := rd.Value(test.path); ok {
				var actual interface{}
				switch test.expect.(type) {
				case int:
					if actual, err = v.Int(); err == nil {
						if !reflect.DeepEqual(test.expect.(int), int(actual.(int64))) {
							t.Fatal(`expect is not equal to actual`)
						}
					}
				case string:
					if actual, err = v.String(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal(`expect is not equal to actual`)
						}
					}
				case bool:
					if actual, err = v.Bool(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal(`expect is not equal to actual`)
						}
					}
				case float64:
					if actual, err = v.Float(); err == nil {
						if !reflect.DeepEqual(test.expect, actual) {
							t.Fatal(`expect is not equal to actual`)
						}
					}
				default:
					actual = v.Load()
					if !reflect.DeepEqual(test.expect, actual) {
						t.Logf("expect: %#v, actural: %#v", test.expect, actual)
						t.Fail()
					}
				}
				if err != nil {
					t.Error(err)
				}
			} else {
				t.Error("value path not found")
			}
		})
	}
}
