// +build windows

package mimefix

import "golang.org/x/sys/windows/registry"

type mimeFixer struct{}

// Fix implementation which sets the registry value
// 'Content Type' of the registry key 'HKCR\\.js' to
// the passed expected mime type string.
func (m *mimeFixer) Fix(expectedMime string) (err error) {
	k, err := registry.OpenKey(registry.CLASSES_ROOT, ".js", registry.SET_VALUE)
	if err != nil {
		return
	}
	defer k.Close()

	err = k.SetStringValue("Content Type", expectedMime)
	return nil
}
