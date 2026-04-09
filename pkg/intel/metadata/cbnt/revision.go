// Copyright 2017-2026 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cbnt

import (
	"encoding/binary"
	"fmt"
	"io"
)

// FIXME: we could stick to having only 1.0 and 2.0, and treat
// "2.1" as 2.0. After all the, in the "hacky" temp fix, the logic
// of handling the header is the same (afaik Intel didn't changed
// anything there), and we only have to treat headers from 21 to 25
// as 20.

type BootGuardVersion uint8

func (bgv BootGuardVersion) String() string {
	switch bgv {
	case Version10:
		return "1.0"
	case Version20:
		return "2.0"
	case Version21:
		return "2.1"
	}
	return "unknown"
}

func DetectBGV(r io.ReadSeeker) (BootGuardVersion, error) {
	// We could take StructInfoBG here as well since version
	// is under the saem offset, so it does not really matter.
	// Plus we just have it here for version detection, so it won't
	// hurt even if read version is actually 0x10.
	var s StructInfoCBNT
	err := binary.Read(r, endianess, &s)
	if err != nil {
		return 0, fmt.Errorf("unable to read field 'ID': %w", err)
	}
	_, err = r.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	switch s.Version {
	case 0x10:
		return Version10, nil
	case 0x20, 0x21:
		return Version20, nil
	case 0x22, 0x23, 0x24, 0x25:
		return Version21, nil
	default:
		return 0, fmt.Errorf("couldn't detect version 0x%x", s.Version)
	}
}
