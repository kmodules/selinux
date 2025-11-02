/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package selinux

import (
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestGetMountSELinuxLabel(t *testing.T) {
	seLinuxOpts1 := v1.SELinuxOptions{
		Level: "s0:c123,c456",
	}
	seLinuxOpts2 := v1.SELinuxOptions{
		Level: "s0:c234,c567",
	}
	seLinuxOpts3 := v1.SELinuxOptions{
		Level: "s0:c345,c678",
	}
	label1 := "system_u:object_r:container_file_t:s0:c123,c456"

	tests := []struct {
		name               string
		seLinuxOptions     []*v1.SELinuxOptions
		expectError        bool
		expectedMountLabel string
	}{
		// Tests with no labels
		{
			name:               "no label",
			seLinuxOptions:     nil,
			expectError:        false,
			expectedMountLabel: "",
		},
		// Tests with one label and RWOP volume
		{
			name:               "one label, Recursive change policy, no feature gate",
			seLinuxOptions:     []*v1.SELinuxOptions{&seLinuxOpts1},
			expectError:        false,
			expectedMountLabel: label1,
		},
		// Corner cases
		{
			name:               "multiple same labels",
			seLinuxOptions:     []*v1.SELinuxOptions{&seLinuxOpts1, &seLinuxOpts1, &seLinuxOpts1, &seLinuxOpts1},
			expectError:        false,
			expectedMountLabel: label1,
		},
		// Error cases
		{
			name:               "multiple different labels",
			seLinuxOptions:     []*v1.SELinuxOptions{&seLinuxOpts1, &seLinuxOpts2, &seLinuxOpts3},
			expectError:        true,
			expectedMountLabel: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			seLinuxTranslator := NewFakeSELinuxLabelTranslator()

			// Act
			info, err := GetMountSELinuxLabel(tt.seLinuxOptions, seLinuxTranslator)

			// Assert
			if err != nil {
				if !tt.expectError {
					t.Errorf("GetMountSELinuxLabel() unexpected error: %v", err)
				}
				return
			}
			if tt.expectError {
				t.Errorf("GetMountSELinuxLabel() expected error, got none")
				return
			}

			if info != tt.expectedMountLabel {
				t.Errorf("GetMountSELinuxLabel() expected %+v, got %+v", tt.expectedMountLabel, info)
			}
		})
	}
}
