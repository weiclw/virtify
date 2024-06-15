package runner

import (
    "virtify/options"
)

type QemuOptions struct {
    ImagePath string
    DiskPath string
}

func NewQemuOptions(opt *options.Options) QemuOptions {
    return QemuOptions{
        ImagePath: "/tmp/alpine-virt-3.18.4-aarch64.iso",
        DiskPath: "/tmp/1.qcow2",
    }
}
