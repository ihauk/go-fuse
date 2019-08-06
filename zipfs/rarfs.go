package zipfs

import (
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/mholt/archiver"
	"github.com/nwaples/rardecode"
	"io/ioutil"
)

type RarFile struct {
	*archiver.File
	//content []byte
	rarPath string
}

func (f *RarFile) Stat(out *fuse.Attr) {
	out.Mode = fuse.S_IFREG | uint32(f.File.Mode())
	out.Size = uint64(f.File.Size())
	out.Mtime = uint64(f.File.ModTime().Unix())
	out.Atime = out.Mtime
	out.Ctime = out.Mtime
}

func (f *RarFile) Data() []byte {
	zf := f.File
	//rc, err :=
	fmt.Println(" ====content file add:", zf)
	//content, err := ioutil.ReadAll(*zf)
	//if err != nil {
	//	panic(err)
	//}
	trh, ok := f.Header.(*rardecode.FileHeader)
	if !ok {
		fmt.Errorf("expected header")
		return nil
	}
	//bytes, err := ioutil.ReadFile(rh.Name)
	r := archiver.NewRar()
	r.Password = "123456"
    var rst []byte
	err := r.Walk(f.rarPath, func(af archiver.File) error {
		rh, ok := af.Header.(*rardecode.FileHeader)
		if !ok {
			return fmt.Errorf("expected header")
		}
		//fmt.Println("FileName:", rh.Name)
		if trh.Name == rh.Name {
			content, e:= ioutil.ReadAll(af)

			if e != nil {
				return e
			}
			rst = content
			fmt.Println(rst)

		}


		return nil
	})
	if err != nil {
		return nil
	}

	return rst
}

// NewZipTree creates a new file-system for the zip file named name.
func NewRARTree(name string) (map[string]MemFile, error) {
	r := archiver.NewRar()
	r.Password = "123456"
	out := map[string]MemFile{}
	err := r.Walk(name, func(f archiver.File) error {
		rh, ok := f.Header.(*rardecode.FileHeader)
		if !ok {
			return fmt.Errorf("expected header")
		}
		//fmt.Println("FileName:", rh.Name)

		//content, err:= ioutil.ReadAll(f)
		//if err != nil {
		//	return err
		//}

		fmt.Println(rh.Name,"file add:", &f)
		out[rh.Name] = &RarFile{&f,name}


		return nil
	})

	//file := out["models/durmodel-pingan-online-20190426/conf"].(*RarFile)
	//content ,err := ioutil.ReadAll(file.File)
	//fmt.Println(err)
	//fmt.Println(content)


	fmt.Println(err)
	return out, nil
}
