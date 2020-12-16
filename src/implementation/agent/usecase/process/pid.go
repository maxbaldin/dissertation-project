package process

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

type ProcessesIndex map[int]Process

func (p ProcessesIndex) InodeExist(inodeId int) (proc Process, b bool) {
	if val, ok := p[inodeId]; ok {
		return val, ok
	}
	return proc, false
}

type Process struct {
	PID  int
	User int
	Name string
}

func GetProcMap() (ProcessesIndex, error) {
	fh, err := os.Open("/proc/")
	if err != nil {
		return nil, err
	}

	dirNames, err := fh.Readdirnames(-1)
	fh.Close()
	if err != nil {
		return nil, err
	}

	var (
		res  = ProcessesIndex{}
		stat syscall.Stat_t
	)
	for _, dirName := range dirNames {
		pid, err := strconv.Atoi(dirName)
		if err != nil {
			continue
		}

		fdBase := filepath.Join("/proc/", dirName, "fd")
		dfh, err := os.Open(fdBase)
		if err != nil {
			continue
		}

		fdNames, err := dfh.Readdirnames(-1)
		_ = dfh.Close()
		if err != nil {
			continue
		}

		err = syscall.Lstat(filepath.Join("/proc/", dirName, "/ns/net"), &stat)
		if err != nil {
			continue
		}

		var name string
		for _, fdName := range fdNames {
			err = syscall.Stat(filepath.Join(fdBase, fdName), &stat)
			if err != nil {
				continue
			}

			if stat.Mode&syscall.S_IFMT != syscall.S_IFSOCK {
				continue
			}

			if name == "" {
				if name = procName(filepath.Join("/proc/", dirName)); name == "" {
					break
				}
			}

			res[int(stat.Ino)] = Process{
				PID:  pid,
				User: int(stat.Uid),
				Name: name,
			}
		}
	}
	return res, nil
}

func procName(base string) string {
	fh, err := os.Open(filepath.Join(base, "/comm"))
	if err != nil {
		return ""
	}

	name := make([]byte, 64)
	l, err := fh.Read(name)
	fh.Close()
	if err != nil {
		return ""
	}

	if l < 2 {
		return ""
	}

	return string(name[:l-1])
}
