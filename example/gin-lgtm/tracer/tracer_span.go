package tracer

import (
	"encoding/json"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (S *Span) SetNative(kv attribute.KeyValue) {
	S.span.SetAttributes(kv)
}

func (S *Span) SetStatus(code codes.Code, desc string) {
	S.span.SetStatus(code, desc)
}

func (S *Span) SetBoolTag(key string, value bool) {

	S.span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.BoolValue(value),
	})

	return
}

func (S *Span) SetStringTag(key, value string) {
	S.span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.StringValue(value),
	})

	return
}

func (S *Span) SetIntTag(key string, value int) {
	S.span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.IntValue(value),
	})

	return
}

func (S *Span) SetInt64Tag(key string, value int64) {
	S.span.SetAttributes(attribute.KeyValue{
		Key:   attribute.Key(key),
		Value: attribute.Int64Value(value),
	})

	return
}

func (S *Span) SetObjectTag(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	S.SetStringTag(key, string(jsonData))
	return nil
}
