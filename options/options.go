package options

import (
    "flag"
    "os"
)

// Each of the instance will represent a command line argment in the function below
type optionValue struct {
    ptr interface{}
    comment string
}

type optionType interface {
    bool | string
}

func newOption[V optionType](default_value V, comment string) optionValue {
    return optionValue{
        ptr: &default_value,
        comment: comment,
    }
}

func GetValuePtr[V optionType](v *optionValue) *V {
    return v.ptr.(*V)
}

type Options struct {
    list map[string]optionValue
}

var RedirectInputFlag = "qemu_redirect_input"
var ActionFileFlag = "qemu_action_file"
var DriveFile1Flag = "drive_file1"
var SockNetFlag = "sock_net"

func NewOptions() *Options {
    return &Options{
        list: map[string]optionValue{
            RedirectInputFlag: newOption(
                false,
                "redirect so that it can run script"),

            ActionFileFlag: newOption(
                "",
                "path of action script"),

            DriveFile1Flag: newOption(
                "/tmp/1.qcow2",
                "path of the first hard drive"),

            SockNetFlag: newOption(
                "",
                "use socket network, either listen=:1234, or connect=:1234"),
        },
    }
}

func GetOptionsPtr[V optionType](o *Options, name string) *V {
    if v, ok := o.list[name]; ok {
        return GetValuePtr[V](&v)
    } else {
        return nil
    }
}

func GetOptionsValue[V optionType](o *Options, name string) (V, bool) {
    if p := GetOptionsPtr[V](o, name); p != nil {
        return *p, true
    } else  {
        var r V
        return r, false
    }
}

func GetOptionsInfo[V optionType](o *Options, name string) (*V, string) {
    if v, ok := o.list[name]; ok {
        return GetValuePtr[V](&v), v.comment
    } else {
        return nil, ""
    }
}
        
// Always extract information from env first.
func optionsFromEnv(opts *Options) {
    if _, ok := os.LookupEnv(RedirectInputFlag); ok {
        ptr := GetOptionsPtr[bool](opts, RedirectInputFlag)
        *ptr = true
    }

    if val, ok := os.LookupEnv(ActionFileFlag); ok {
       ptr := GetOptionsPtr[string](opts, ActionFileFlag)
       *ptr = val
    }

    if val, ok := os.LookupEnv(DriveFile1Flag); ok {
       ptr := GetOptionsPtr[string](opts, DriveFile1Flag)
       *ptr = val
    }

    if val, ok := os.LookupEnv(SockNetFlag); ok {
       ptr := GetOptionsPtr[string](opts, SockNetFlag)
       *ptr = val
    }
}


// Commandline flags shall override values from env.
func optionsFromFlags(opts *Options) {
    if ptr, comment := GetOptionsInfo[bool](opts, RedirectInputFlag); true {
        default_val := *ptr
        flag.BoolVar(ptr, RedirectInputFlag, default_val, comment)
    }

    if ptr, comment := GetOptionsInfo[string](opts, ActionFileFlag); true {
        default_val := *ptr
        flag.StringVar(ptr, ActionFileFlag, default_val, comment)
    }

    if ptr, comment := GetOptionsInfo[string](opts, DriveFile1Flag); true {
        default_val := *ptr
        flag.StringVar(ptr, DriveFile1Flag, default_val, comment)
    }

    if ptr, comment := GetOptionsInfo[string](opts, SockNetFlag); true {
        default_val := *ptr
        flag.StringVar(ptr, SockNetFlag, default_val, comment)
    }

    // This function automatically handles parsing error and may exit the program as well.
    flag.Parse()
}

func GetOptionsOnce(opts *Options) {
    optionsFromEnv(opts)
    optionsFromFlags(opts)
}
