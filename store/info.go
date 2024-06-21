package store


type ImageInfo struct {
    Name string
    IsoPath string
    Disk1Path string
    ParentName string
}

type VMInfo struct {
    Name string
    StartTime string
    State string
    ImageName string
}

type MetaStore interface {
    GetImage(imageName string) (ImageInfo, bool)
    GetVM(vmName string) (VMInfo, bool)
    ListImages() []string
    ListVMs() []string
    SaveImage(name string, info ImageInfo)
    UpdateVM(name string, info VMInfo)
}

func NewImageInfo() ImageInfo {
    return ImageInfo{}
}

func NewVMInfo() VMInfo {
    return VMInfo{}
}
