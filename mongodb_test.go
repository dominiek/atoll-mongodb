package main

import (
  "testing"
  //"fmt"
  "github.com/stretchr/testify/assert"
  "github.com/jeffail/gabs"
)

func TestMongoDBReport(t *testing.T) {
  mongodb := MongoDB{"localhost", 26017};
  data, err := mongodb.Monitor();
  assert.Equal(t, err, nil)

  t.Logf("Report: %v", data)

  jsonParsed, err := gabs.ParseJSON([]byte(data))
  assert.Equal(t, err, nil)

  children, _ := jsonParsed.S("report").S("items").Children();
  assert.Equal(t, len(children), 3)

  status, _ := jsonParsed.S("report").S("status").S("state").Data().(string);
  assert.Equal(t, status, "ok")
}
