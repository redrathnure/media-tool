package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasFileFilter_Valid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")

	sut := HasFileFilter("test")

	assert.True(t, sut.Accept(&dev))
}

func TestHasFileFilter_Invalid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")

	sut := HasFileFilter("not_test.jpg")

	assert.False(t, sut.Accept(&dev))
}

func TestHasDeviceNameFilter_Valid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")

	sut := HasDeviceNameFilter("Mock")

	assert.True(t, sut.Accept(&dev))
}

func TestHasDeviceNameFilter_Invalid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")

	sut := HasDeviceNameFilter("Non Mock")

	assert.False(t, sut.Accept(&dev))
}

func TestNotFilter_Valid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")
	f1 := HasDeviceNameFilter("Mock")

	sut := NotFilter(f1)

	assert.False(t, sut.Accept(&dev))
}

func TestAndFilter_Valid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")
	f1 := HasDeviceNameFilter("Mock")
	f2 := HasDeviceNameFilter("Mock")

	sut := AndFilter(f1, f2)

	assert.True(t, sut.Accept(&dev))
}

func TestAndFilter_Invalid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")
	f1 := HasDeviceNameFilter("Mock")
	f2 := HasDeviceNameFilter("Not Mock")

	sut := AndFilter(f1, f2)

	assert.False(t, sut.Accept(&dev))
}

func TestOrFilter_Valid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")
	f1 := HasDeviceNameFilter("Mock")
	f2 := HasDeviceNameFilter("Not Mock")

	sut := OrFilter(f1, f2)

	assert.True(t, sut.Accept(&dev))
}

func TestOrFilter_Invalid(t *testing.T) {
	dev := NewMockDev("Mock", "test.jpg")
	f1 := HasDeviceNameFilter("Not Mock")
	f2 := HasDeviceNameFilter("Not Mock")

	sut := OrFilter(f1, f2)

	assert.False(t, sut.Accept(&dev))
}

func TestGoProFilter_Valid(t *testing.T) {
	tests := []RemovableDevice{
		NewMockDev("Mock", GOPRO_DIR),
		NewMockDev("HERO", GOPRO_DIR),
		NewMockDev("GoPro", GOPRO_DIR),
		NewMockDev("HERO_empty"),
		NewMockDev("GoPro_empty"),
	}

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			assert.True(t, GoProFilter.Accept(&tt))
		})
	}
}

func TestGoProFilter_Invalid(t *testing.T) {
	tests := []RemovableDevice{
		NewMockDev("Mock_empty"),
		NewMockDev("Mock_DCIM", DCIM_DIR),
	}

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			assert.False(t, GoProFilter.Accept(&tt))
		})
	}
}

func TestCamFilter_Valid(t *testing.T) {
	tests := []RemovableDevice{
		NewMockDev("CAM_CAM_FILES_DIR", CAM_FILES_DIR),
	}

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			assert.True(t, CamFilter.Accept(&tt))
		})
	}
}

func TestCamFilter_Invalid(t *testing.T) {
	tests := []RemovableDevice{
		NewMockDev("Mock_empty"),
		NewMockDev("Mock_DCIM", DCIM_DIR),
		NewMockDev("Mock_FILES_DIR", CAM_FILES_DIR),
		NewMockDev("CAM_empty"),
		NewMockDev("CAM_DCIM", DCIM_DIR),
	}

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			assert.False(t, CamFilter.Accept(&tt))
		})
	}
}

func TestSdPhotosFilter_Valid(t *testing.T) {
	tests := []RemovableDevice{
		NewMockDev("Mock_DCIM", DCIM_DIR),
		NewMockDev("CAM_DCIM", DCIM_DIR),
	}

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			assert.True(t, SdPhotosFilter.Accept(&tt))
		})
	}
}

func TestSdPhotosFilter_Invalid(t *testing.T) {
	tests := []RemovableDevice{
		NewMockDev("Mock_empty"),
		NewMockDev("Mock_FILES_DIR", CAM_FILES_DIR),
		NewMockDev("CAM_empty"),

		//VideoCam
		NewMockDev("CAM_CAM_FILES_DIR", CAM_FILES_DIR),

		//GoPro
		NewMockDev("Mock", GOPRO_DIR),
		NewMockDev("HERO", GOPRO_DIR),
		NewMockDev("GoPro", GOPRO_DIR),
		NewMockDev("HERO_empty"),
		NewMockDev("GoPro_empty"),
	}

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			assert.False(t, SdPhotosFilter.Accept(&tt))
		})
	}
}
