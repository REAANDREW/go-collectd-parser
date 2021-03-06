package collectd

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

var packetBytes []byte
var err error

func TestMain(m *testing.M) {
	packetBytes, err = ioutil.ReadFile("cpu_disk_packet_5.4.0.dat")
	if err != nil {
		fmt.Errorf("error encountered %v", err)
	}

	os.Exit(m.Run())
}

func FindStringParts(partTypeId uint16, parts []interface{}) (stringParts []StringPart, err error) {
	var returnParts []StringPart
	for _, value := range parts {
		switch value := value.(type) {
		case StringPart:
			if value.Header.Type == partTypeId {
				returnParts = append(returnParts, value)
			}
		}
	}
	return returnParts, nil
}

func FindNumericParts(partTypeId uint16, parts []interface{}) (numericParts []NumericPart, err error) {
	var returnParts []NumericPart
	for _, value := range parts {
		switch value := value.(type) {
		case NumericPart:
			if value.Header.Type == partTypeId {
				returnParts = append(returnParts, value)
			}
		}
	}
	return returnParts, nil
}

func FindValueParts(parts []interface{}) (valueParts []ValuePart, err error) {
	var returnParts []ValuePart
	for _, value := range parts {
		switch value := value.(type) {
		case ValuePart:
			if value.Header.Type == VALUE {
				returnParts = append(returnParts, value)
			}
		}
	}
	return returnParts, nil
}

func Test_parsesTheHostname(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	string_parts, _ := FindStringParts(HOSTNAME, parts)
	assert.Equal(t, 1, len(string_parts), "number of parts is not equal to 1")
	assert.Equal(t, string_parts[0].Content, "localhost", "contents does not equal localhost")
}

func Test_parsesTheHighDefinitionTime(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	numeric_parts, _ := FindNumericParts(HIGH_DEF_TIME, parts)
	assert.Equal(t, 26, len(numeric_parts), "number of parts is not equal to 26")
	assert.Equal(t, 1419765641, numeric_parts[0].Content, "contents does not equal expected")
}

func Test_parsesThePlugin(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	string_parts, _ := FindStringParts(PLUGIN, parts)
	assert.Equal(t, 2, len(string_parts), "number of parts is not equal to 2")
	assert.Equal(t, string_parts[0].Content, "disk", "plugin content does not equal expected")
	assert.Equal(t, string_parts[1].Content, "cpu", "plugin content does not equal expected")
}

func Test_parsesThePluginInstance(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	string_parts, _ := FindStringParts(PLUGIN_INSTANCE, parts)
	assert.Equal(t, 6, len(string_parts), "number of parts is not equal to 6")
	assert.Equal(t, "sda1", string_parts[0].Content, "plugin content does not equal expected")
	assert.Equal(t, "sda2", string_parts[1].Content, "plugin content does not equal expected")
	assert.Equal(t, "sda5", string_parts[2].Content, "plugin content does not equal expected")
	assert.Equal(t, "dm-0", string_parts[3].Content, "plugin content does not equal expected")
	assert.Equal(t, "dm-1", string_parts[4].Content, "plugin content does not equal expected")
	assert.Equal(t, "0", string_parts[5].Content, "plugin content does not equal expected")
}

func Test_parsesTheType(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	string_parts, _ := FindStringParts(TYPE, parts)
	assert.Equal(t, 21, len(string_parts), "number of parts is not equal to 21")
	for counter := 0; counter < 20; counter += 4 {
		assert.Equal(t, "disk_octets", string_parts[counter].Content, "plugin content does not equal expected")
		assert.Equal(t, "disk_ops", string_parts[counter+1].Content, "plugin content does not equal expected")
		assert.Equal(t, "disk_time", string_parts[counter+2].Content, "plugin content does not equal expected")
		assert.Equal(t, "disk_merged", string_parts[counter+3].Content, "plugin content does not equal expected")
	}
	assert.Equal(t, "cpu", string_parts[20].Content, "plugin content does not equal expected")
}

func Test_parsesTheTypeInstance(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	string_parts, _ := FindStringParts(TYPE_INSTANCE, parts)
	assert.Equal(t, 6, len(string_parts), "number of parts is not equal to 6")
	assert.Equal(t, "user", string_parts[0].Content, "plugin content does not equal expected")
	assert.Equal(t, "nice", string_parts[1].Content, "plugin content does not equal expected")
	assert.Equal(t, "system", string_parts[2].Content, "plugin content does not equal expected")
	assert.Equal(t, "idle", string_parts[3].Content, "plugin content does not equal expected")
	assert.Equal(t, "wait", string_parts[4].Content, "plugin content does not equal expected")
	assert.Equal(t, "interrupt", string_parts[5].Content, "plugin content does not equal expected")
}

func Test_parsesTheHighDefinitionInterval(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	numeric_parts, _ := FindNumericParts(HIGH_DEF_INTERVAL, parts)
	assert.Equal(t, 1, len(numeric_parts), "number of parts is not equal to 1")
	assert.Equal(t, 2, numeric_parts[0].Content, "contents does not equal expected")
}

func Test_parsesTheValuesPart(t *testing.T) {
	buffer := bytes.NewBuffer(packetBytes)
	parts := parseParts(buffer)
	value_parts, _ := FindValueParts(parts)
	assert.Equal(t, 26, len(value_parts), "number of parts is not equal to 1")
}
