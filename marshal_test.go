package marshaler

import (
	"log"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type Basic struct {
		StringValue string
		IntValue    int32
		BoolValue   bool
		FloatValue  float32
	}
	type BasicNested struct {
		StringValue string
		IntValue    int32
		BoolValue   bool
		FloatValue  float32
		DeepNested  Basic
	}
	type testStruct struct {
		StringValue string
		IntValue    int32
		BoolValue   bool
		FloatValue  float32
		Nested      BasicNested
	}

	type pointer struct {
		Pointer *Basic
	}

	type recursive struct {
		StringValue string
		Recursive   *recursive
	}

	data := []struct {
		args any
		want []byte

		name string
	}{
		{name: "Marshals unnested", args: Basic{StringValue: "MyString", IntValue: 50, BoolValue: true, FloatValue: 50.5}, want: []byte("StringValue=MyString IntValue=50 BoolValue=true FloatValue=50.5")},
		{name: "Marshals nested", args: testStruct{StringValue: "MyString", IntValue: 50, BoolValue: true, FloatValue: 5.0, Nested: BasicNested{StringValue: "myNestedString", IntValue: 100, BoolValue: true, FloatValue: 10.5, DeepNested: Basic{StringValue: "myNestedString", IntValue: 100, BoolValue: true, FloatValue: 10.5}}}, want: []byte("StringValue=MyString IntValue=50 BoolValue=true FloatValue=5 Nested.StringValue=myNestedString Nested.IntValue=100 Nested.BoolValue=true Nested.FloatValue=10.5 Nested.DeepNested.StringValue=myNestedString Nested.DeepNested.IntValue=100 Nested.DeepNested.BoolValue=true Nested.DeepNested.FloatValue=10.5")},
		{name: "Marshals with pointer", args: pointer{Pointer: &Basic{StringValue: "MyString", IntValue: 50, BoolValue: true, FloatValue: 50.5}}, want: []byte("Pointer.StringValue=MyString Pointer.IntValue=50 Pointer.BoolValue=true Pointer.FloatValue=50.5")},
		{name: "Marshals empty", args: struct{}{}, want: []byte("")},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			result := Marshal(tt.args)
			if !reflect.DeepEqual(result, tt.want) {
				t.Fatalf("\nExpected:\t %v \ngot:\t\t %v", string(tt.want), string(result))
				t.Fatalf("\nExpected:\t %v \ngot:\t\t %v", tt.want, result)
			}

			log.Println(string(result))
		})
	}
}
