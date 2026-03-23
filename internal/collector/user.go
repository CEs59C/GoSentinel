package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

type UserInfo struct {
	User     string
	Host     string
	Started  string
	Terminal string
}

func GetUserInfo() ([]UserInfo, error) {
	users, err := host.Users()
	if err != nil {
		return []UserInfo{}, fmt.Errorf("failed to get users info: %w", err)
	}

	var ul = make([]UserInfo, len(users))
	for i, us := range users {
		ul[i].User = us.User
		ul[i].Terminal = us.Terminal
		ul[i].Host = us.Host
		ul[i].Started = time.Unix(int64(us.Started), 0).Format("15:04:05")
	}

	//for _, u := range ul {
	//	fmt.Println(u)
	//}
	return ul, nil
}

func (u UserInfo) String() string {
	return fmt.Sprintf("User:\t\t%s, Host=%s, Started=%s, Terminal=%s.",
		u.User, u.Host, u.Started, u.Terminal)
}
