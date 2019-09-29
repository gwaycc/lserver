package etc

import (
	"testing"
)

func TestEtc(t *testing.T) {
	Etc.GetString("applet/cms", "cookie_name")
}
