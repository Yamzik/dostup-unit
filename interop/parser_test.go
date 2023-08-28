package interop

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

const input = `interface: wg0
public key: 4Vai3DvLEqHGxTSSJIhbPvsutVM+Ftqr7wTxSNmAKBI=
private key: (hidden)
listening port: 51820

peer: 67kGUWckjYpKvgDP1mt/r5sf1JFBhe+72pyqjB+WgFk=
preshared key: (hidden)
endpoint: 192.168.48.1:48774
allowed ips: 10.8.0.2/32
latest handshake: 19 seconds ago
transfer: 28.41 KiB received, 35.62 KiB sent

peer: 67kGUWckjYpKvgDP1mt/r5sf1JFBhe+72pyqjB+WgFk=
preshared key: (hidden)
endpoint: 192.168.48.1:48774
allowed ips: 10.8.0.2/32
latest handshake: 19 seconds ago
transfer: 28.41 KiB received, 35.62 KiB sent`

func equal(t *testing.T, expected, actual any) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	_, fn, line, _ := runtime.Caller(1)
	t.Errorf("Failed equals at %s:%d\nactual   %#v\nexpected %#v", fn, line, actual, expected)
	return false
}

func lenTest(t *testing.T, actualO any, expected int) bool {
	actual := reflect.ValueOf(actualO).Len()
	if reflect.DeepEqual(expected, actual) {
		return true
	}
	_, fn, line, _ := runtime.Caller(1)
	t.Errorf("Wrong length at %s:%d\nactual   %#v\nexpected %#v", fn, line, actual, expected)
	return false
}

func TestParseShow(t *testing.T) {
	result, _ := parseShow(input)
	lenTest(t, result.Peers, 2)
	equal(t, result.Peers[0].PublicKey, "67kGUWckjYpKvgDP1mt/r5sf1JFBhe+72pyqjB+WgFk=")
	equal(t, result.Peers[1].LatestHandshake, "19 seconds ago")
	fmt.Printf("%v", result)
}
