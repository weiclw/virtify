package runner

import (
    "fmt"
    "math/rand"
    "os"
    "strings"
    "virtify/options"
)

type qemuOptions struct {
    imagePath string
    drive1Path string
    sockNet string
}

func genUniqueMac() string {
    mac := fmt.Sprintf("%02x", 0x2)
    for i := 1; i < 6; i++ {
        mac = mac + fmt.Sprintf(":%02x", byte(rand.Intn(256)))
    }

    return mac
}

func NewQemuOptions(opt *options.Options) qemuOptions {
    drive1Path := ""
    if val, got := options.GetOptionsValue[string](opt, options.DriveFile1Flag); got {
        drive1Path = val
    }

    sockNet := ""
    if val, got := options.GetOptionsValue[string](opt, options.SockNetFlag); got {
        sockNet = val
    }

    return qemuOptions{
        imagePath: "/tmp/alpine-virt-3.18.4-aarch64.iso",
        drive1Path: drive1Path,
        sockNet: sockNet,
    }
}

func (opts *qemuOptions) genInitialCommand() []string {
    return []string{
        "/opt/homebrew/bin/qemu-system-aarch64",
        "-cpu",
        "host",
        "-smp",
        "cpus=2,sockets=1,cores=2,threads=1",
        "-machine",
        "virt,",
        "-accel",
        "hvf",
        "-accel",
        "tcg,tb-size=512",
    }
}

func (opts *qemuOptions) genNetworkPart() []string {
    if len(opts.sockNet) == 0 {
        return []string{}
    } else if strings.HasPrefix(opts.sockNet, "listen=") {
        return []string{
            "-netdev",
            "user,id=n1",
            "-netdev",
            "hubport,hubid=1,netdev=n1,id=h1",
            "-netdev",
            "socket,id=n2," + opts.sockNet,
            "-netdev",
            "hubport,hubid=1,netdev=n2,id=h2",
            "-netdev",
            "hubport,hubid=1,id=h3",
            "-device",
            "virtio-net-pci,netdev=h3,mac=" + genUniqueMac(),
        }
    } else if strings.HasPrefix(opts.sockNet, "connect=") {
        return []string{
            "-netdev",
            "socket,id=n1," + opts.sockNet,
            "-device",
            "virtio-net-pci,netdev=n1,mac=" + genUniqueMac(),
        }
    } else {
        fmt.Println("Incorrect sock_net option: " + opts.sockNet)
        os.Exit(1)
        return []string{}
    }
}

func (opts *qemuOptions) genRest() []string {
    drive1opt := "if=none,media=disk,id=drive1,file=" + opts.drive1Path
    return []string{
        "-drive",
        "if=pflash,format=raw,unit=0,file=/Applications/UTM.app/Contents/Resources/qemu/edk2-aarch64-code.fd,readonly=on",
        "-drive",
        "if=pflash,unit=1,file=/Applications/UTM.app/Contents/Resources/qemu/edk2-arm-vars.fd,readonly=on",
        "-nographic",
        "-boot",
        "menu=on",
        "-m",
        "2048",
        "-device",
        "intel-hda",
        "-device",
        "virtio-blk-pci,drive=drive0,bootindex=0",
        "-drive",
        "if=none,media=disk,id=drive0,file=/tmp/alpine-virt-3.18.4-aarch64.iso,readonly=on",
        "-device",
        "virtio-blk-pci,drive=drive1,bootindex=1",
        "-drive",
        drive1opt,
    }
}

func (opts *qemuOptions) getCommand() []string {
    cmd := opts.genInitialCommand()
    cmd = append(cmd, opts.genNetworkPart()...)
    cmd = append(cmd, opts.genRest()...)
    return cmd
}
