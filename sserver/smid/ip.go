package smid

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/slib/serror"
	"github.com/yyliziqiu/slib/slog"
	"github.com/yyliziqiu/slib/sserver"
	"github.com/yyliziqiu/slib/sserver/sresp"
)

var (
	ErrInvalidIp = errors.New("invalid ip")

	_whiteList []uint32
)

func CheckIp(ips []string) gin.HandlerFunc {
	for _, ip := range ips {
		iv, err := parseIp(ip)
		if err != nil {
			slog.Errorf("Parse ip failed, ip: %s, error: %v.", ip, err)
			continue
		}
		_whiteList = append(_whiteList, iv)
	}

	return func(ctx *gin.Context) {
		// ip := ctx.RemoteIP() // 这个 ip 地址是客户端与服务器建立 tcp 连接时的原始 ip
		// ip := ctx.ClientIP() // 先检查 X-Forwarded-For、X-Real-IP 等 header，如果没有代理信息或代理不受信任则退回到 RemoteIP()

		ip := ctx.RemoteIP()

		iv, err := ip2int(ip)
		if err != nil {
			logger := sserver.GetLogger()
			if logger != nil {
				logger.Warnf("Parse ip failed, ip: %s, error: %v.", ip, err)
			}
			sresp.AbortError(ctx, serror.ForbiddenIp)
			return
		}

		for _, match := range _whiteList {
			if iv&match == match {
				return
			}
		}

		sresp.AbortError(ctx, serror.ForbiddenIp)
	}
}

func ip2int(ip string) (uint32, error) {
	ipv4 := net.ParseIP(ip).To4()
	if len(ipv4) == 0 {
		return 0, ErrInvalidIp
	}

	i := uint32(0)
	for _, b := range ipv4 {
		i = i << 8
		i |= uint32(b)
	}

	return i, nil
}

func parseIp(ip string) (uint32, error) {
	pm := strings.Split(ip, "/")
	if len(pm) == 0 || len(pm) > 2 {
		return 0, ErrInvalidIp
	}

	mask := 32
	if len(pm) == 2 {
		mk, err := strconv.Atoi(pm[1])
		if err != nil {
			return 0, ErrInvalidIp
		}
		ip = pm[0]
		mask = mk
	}

	i, err := ip2int(ip)
	if err != nil {
		return 0, err
	}

	return i & (uint32(0xffffffff) << (32 - mask)), nil
}
