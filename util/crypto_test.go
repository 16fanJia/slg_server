package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAesCBCEncrypt(t *testing.T) {
	testCases := []struct {
		name        string
		src         []byte
		key         []byte
		paddingMode string

		wantRes string
		wantErr error
	}{
		{
			name:        "cbcEncrypt",
			src:         []byte("加密测试数据"),
			key:         []byte("1234567812345678"),
			paddingMode: PKCS7_PADDING,
			wantErr:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := AesCBCEncrypt(tc.src, tc.key, tc.paddingMode)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
		})
	}

}
