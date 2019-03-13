// Copyright (c) 2018-2019, Sylabs Inc. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/sylabs/singularity/internal/pkg/test"
)

func imagePull(library string, imagePath string, sourceSpec string, force, unauthenticated bool) ([]byte, error) {
	var argv []string
	argv = append(argv, "pull")
	if force {
		argv = append(argv, "--force")
	}
	if unauthenticated {
		argv = append(argv, "--allow-unauthenticated")
	}
	if library != "" {
		argv = append(argv, "--library", library)
	}
	if imagePath != "" {
		argv = append(argv, imagePath)
	}
	argv = append(argv, sourceSpec)

	return exec.Command(cmdPath, argv...).CombinedOutput()
}

func TestPull(t *testing.T) {
	test.DropPrivilege(t)

	imagePath := "./test_pull.sif"

	tests := []struct {
		name            string
		sourceSpec      string
		force           bool
		unauthenticated bool
		library         string
		imagePath       string
		success         bool
	}{
		{"Pull_Library", "library://alpine:3.8", false, false, "", imagePath, true}, // https://cloud.sylabs.io/library
		{"Force", "library://alpine:3.8", true, false, "", imagePath, true},
		{"Unsigned_image", "library://sylabs/tests/unsigned:1.0.0", true, true, "", imagePath, true},
		{"Unsigned_image_fail", "library://sylabs/tests/unsigned:1.0.0", true, false, "", imagePath, false}, // pull a unsigned image; should fail
		{"Pull_Docker", "docker://alpine:3.8", true, false, "", imagePath, true},                            // https://hub.docker.com/
		{"Pull_Shub", "shub://GodloveD/busybox", true, false, "", imagePath, true},                          // https://singularity-hub.org/
		{"PullWithHash", "library://sylabs/tests/signed:sha256.5c439fd262095766693dae95fb81334c3a02a7f0e4dc6291e0648ed4ddc61c6c", true, false, "", imagePath, true},
		{"PullWithoutTransportProtocol", "alpine:3.8", true, false, "", imagePath, true},
		{"PullNonExistent", "library://this_should_not/exist/not_exist", true, false, "", imagePath, false}, // pull a non-existent container
	}
	defer os.Remove(imagePath)
	for _, tt := range tests {
		t.Run(tt.name, test.WithoutPrivilege(func(t *testing.T) {
			b, err := imagePull(tt.library, tt.imagePath, tt.sourceSpec, tt.force, tt.unauthenticated)
			if tt.success {
				if err != nil {
					t.Log(string(b))
					t.Fatalf("unexpected failure: %v", err)
				}
			} else {
				if err == nil {
					t.Log(string(b))
					t.Fatalf("unexpected success: command should have failed")
				}
			}
			imageVerify(t, tt.imagePath, false)
		}))
	}
}
