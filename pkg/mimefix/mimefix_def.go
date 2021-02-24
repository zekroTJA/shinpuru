// +build !windows

package mimefix

type mimeFixer struct{}

// Fix implementation which skips the fix on all
// systems instead of windows.
func (m *mimeFixer) Fix(expectedMime string) error {
	return nil
}
