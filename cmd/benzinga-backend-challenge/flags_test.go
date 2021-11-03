package main

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestInitFlagsFromEnvNoEnvSet(t *testing.T) {
	flg := &flags{
		http:          ":8080",
		batchSize:     10,
		batchInterval: 10 * time.Second,
	}
	initFlagsFromEnv(flg)
	expectedFlg := &flags{
		http:          ":8080",
		batchSize:     10,
		batchInterval: 10 * time.Second,
	}
	if !reflect.DeepEqual(flg, expectedFlg) {
		t.Fatalf("unexpected flags got: %#v want: %#v", flg, expectedFlg)
	}
}

func TestInitFlagsFromEnvPostEndpoint(t *testing.T) {
	flg := &flags{}
	err := os.Setenv(postEndPointEnvVar, "https://request.bin/rand_123")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(postEndPointEnvVar) //nolint:errcheck
	initFlagsFromEnv(flg)
	expectedFlg := &flags{
		postEndpoint: "https://request.bin/rand_123",
	}
	if !reflect.DeepEqual(flg, expectedFlg) {
		t.Fatalf("unexpected flags got: %#v want: %#v", flg, expectedFlg)
	}
}

func TestInitFlagsFromEnvPostEndpointErr(t *testing.T) {
	flg := &flags{}
	err := os.Setenv(postEndPointEnvVar, "://invalid.request.bin/rand_123")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(postEndPointEnvVar) //nolint:errcheck
	initFlagsFromEnv(flg)
	expectedFlg := &flags{
		postEndpoint: "",
	}
	if !reflect.DeepEqual(flg, expectedFlg) {
		t.Fatalf("unexpected flags got: %#v want: %#v", flg, expectedFlg)
	}
}

func TestInitFlagsFromEnvBatchSize(t *testing.T) {
	flg := &flags{}
	err := os.Setenv(batchSizeEnvVar, "10")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(batchSizeEnvVar) //nolint:errcheck
	initFlagsFromEnv(flg)
	expectedFlg := &flags{
		batchSize: 10,
	}
	if !reflect.DeepEqual(flg, expectedFlg) {
		t.Fatalf("unexpected flags got: %#v want: %#v", flg, expectedFlg)
	}
}

func TestInitFlagsFromEnvBatchSizeInvalid(t *testing.T) {
	flg := &flags{}
	err := os.Setenv(batchSizeEnvVar, "invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(batchSizeEnvVar) //nolint:errcheck
	initFlagsFromEnv(flg)
	expectedFlg := &flags{
		batchSize: 0,
	}
	if !reflect.DeepEqual(flg, expectedFlg) {
		t.Fatalf("unexpected flags got: %#v want: %#v", flg, expectedFlg)
	}
}

func TestInitFlagsFromEnvBatchInterval(t *testing.T) {
	flg := &flags{}
	err := os.Setenv(batchIntervalEnvVar, "1s")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(batchIntervalEnvVar) //nolint:errcheck
	initFlagsFromEnv(flg)
	expectedFlg := &flags{
		batchInterval: 1 * time.Second,
	}
	if !reflect.DeepEqual(flg, expectedFlg) {
		t.Fatalf("unexpected flags got: %#v want: %#v", flg, expectedFlg)
	}
}

func TestInitFlagsFromEnvBatchIntervalInvalid(t *testing.T) {
	flg := &flags{}
	err := os.Setenv(batchIntervalEnvVar, "invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv(batchIntervalEnvVar) //nolint:errcheck
	initFlagsFromEnv(flg)
	expectedFlg := &flags{
		batchInterval: 0,
	}
	if !reflect.DeepEqual(flg, expectedFlg) {
		t.Fatalf("unexpected flags got: %#v want: %#v", flg, expectedFlg)
	}
}
