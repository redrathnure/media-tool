package mtp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_wpdFile_relPath(t *testing.T) {
	tests := []struct {
		name     string
		sut      wpdFile
		basepath string
		want     string
	}{
		{
			name:     "EmptyBasePath",
			sut:      buildFile("\\tmp\\test\\file1"),
			basepath: "",
			want:     "\\tmp\\test\\file1",
		},
		{
			name:     "NormalCase",
			sut:      buildFile("\\tmp\\test\\file1"),
			basepath: "\\tmp",
			want:     "test\\file1",
		},
		{
			name:     "DifferentRoots",
			sut:      buildFile("\\tmp1\\test\\file1"),
			basepath: "\\tmp2",
			want:     "..\\tmp1\\test\\file1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			r.Equal(tt.want, tt.sut.relPath(tt.basepath))
		})
	}
}

func Test_isIgnored_Ignored(t *testing.T) {
	r := require.New(t)

	r.True(isIgnored("System Volume Information"))
	r.True(isIgnored("$RECYCLE.BIN"))
}

func Test_isIgnored_Normal(t *testing.T) {
	r := require.New(t)

	r.False(isIgnored("DCIM"))
	r.False(isIgnored("test"))
}

func buildFile(filePath string) wpdFile {
	return wpdFile{filePath: filePath}
}
