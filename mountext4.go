package gokrazy

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	BTRFS_CSUM_SIZE = 32
	BTRFS_FSID_SIZE = 16
	BTRFS_MAGIC     = 0x4D5F53665248425F
)

type btrfs_super_block struct {
	/* The first 4 fields must match struct btrfs_header */
	Csum [BTRFS_CSUM_SIZE]byte
	/* FS specific UUID, visible to user */
	Fsid [BTRFS_FSID_SIZE]byte
	/* This block number */
	Bytenr uint64
	Flags  uint64

	/* Allowed to be different from the btrfs_header from here own down */
	Magic      uint64
	Generation uint64
	Root       uint64
	Chunk_root uint64
	Log_root   uint64
}

func probeBTRFS(r io.ReadSeeker) (string, error) {
	// probe btrfs
	const btrfsSuperblockOffset = 0x10000
	if _, err := r.Seek(btrfsSuperblockOffset, io.SeekStart); err != nil {
		return "", err
	}
	var sb btrfs_super_block
	if err := binary.Read(r, binary.LittleEndian, &sb); err != nil {
		return "", err
	}
	if sb.Magic != BTRFS_MAGIC {
		return "", fmt.Errorf("no btrfs superblock found")
	}
	buf := sb.Fsid
	return fmt.Sprintf(
		"%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		buf[0], buf[1], buf[2], buf[3],
		buf[4], buf[5],
		buf[6], buf[7],
		buf[8], buf[9],
		buf[10], buf[11], buf[12], buf[13], buf[14], buf[15]), nil
}

func probeExt4(r io.ReadSeeker) (string, error) {
	// probe ext4
	const extSuperblockOffset = 0x400
	if _, err := r.Seek(extSuperblockOffset, io.SeekStart); err != nil {
		return "", err
	}
	var sb ext2SuperBlock
	if err := binary.Read(r, binary.LittleEndian, &sb); err != nil {
		return "", err
	}
	if sb.Magic != 0xef53 {
		return "", fmt.Errorf("no ext4 superblock found")
	}
	buf := sb.UUID
	return fmt.Sprintf(
		"%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		buf[0], buf[1], buf[2], buf[3],
		buf[4], buf[5],
		buf[6], buf[7],
		buf[8], buf[9],
		buf[10], buf[11], buf[12], buf[13], buf[14], buf[15]), nil
}

type ext2SuperBlock struct {
	InodesCount         uint32
	BlocksCount         uint32
	ReservedBlocksCount uint32
	FreeBlocksCount     uint32
	FreeInodesCount     uint32
	FirstDataBlock      uint32
	LogBlockSize        uint32
	LogFragSize         int32
	BlocksPerGroup      uint32
	FragsPerGroup       uint32
	InodesPerGroup      uint32
	Mtime               uint32
	Wtime               uint32
	MountCount          uint16
	MaxMountCount       int16
	Magic               uint16
	State               uint16
	Errors              uint16
	MinorRevLevel       uint16
	LastCheck           uint32
	CheckInterval       uint32
	CreatorOS           uint32
	RevLevel            uint32
	DefResuid           uint16
	DefResgid           uint16
	FirstIno            uint32
	InodeSize           uint16
	BlockGroupNr        uint16
	FeatureCompat       uint32
	FeatureIncompat     uint32
	FeatureRoCompat     uint32
	UUID                [16]uint8

	// remaining fields elided (irrelevant for probing)
}
