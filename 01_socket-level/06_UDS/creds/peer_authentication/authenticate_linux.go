package auth

import (
	"golang.org/x/sys/unix"
	"log"
	"net"
	"os/user"
	"strconv"
)

// works on Linux and Darwin
func Allowed(conn *net.UnixConn, groups map[string]struct{}) bool {
	if conn == nil || groups == nil || len(groups) == 0 {
		return false
	}

	file, _ := conn.File()
	defer func() { _ = file.Close() }()

	var (
		err   error
		ucred *unix.Ucred
	)

	for {
		ucred, err = unix.GetsockoptUcred(int(file.Fd()), unix.SOL_SOCKET, unix.SO_PEERCRED)
		if err == unix.EINTR {
			continue // syscall interrupted, try again
		}

		if err != nil {
			log.Println(err)
			return false
		}

		break
	}
	u, err := user.LookupId(strconv.FormatUint(uint64(ucred.Uid), 10))
	if err != nil {
		log.Println(err)
		return false
	}

	gids, err := u.GroupIds()
	if err != nil {
		log.Println(err)
		return false
	}

	for _, gid := range gids {
		if _, ok := groups[gid]; ok {
			return true
		}
	}
	return false
}
