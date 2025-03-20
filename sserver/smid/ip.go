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
	ErrInvalid = errors.New("invalid ip")

	_ranges []uint32
)

func CheckIp(ips []string) gin.HandlerFunc {
	for _, ip := range ips {
		r, err := parseRange(ip)
		if err != nil {
			slog.Errorf("Parse ip failed, ip: %s, error: %v.", ip, err)
			continue
		}
		_ranges = append(_ranges, r)
	}

	return func(ctx *gin.Context) {
		// ip := ctx.RemoteIP() // 这个 ip 地址是客户端与服务器建立 tcp 连接时的原始 ip
		// ip := ctx.ClientIP() // 先检查 X-Forwarded-For、X-Real-IP 等 header，如果没有代理信息或代理不受信任则退回到 RemoteIP()

		ip := ctx.RemoteIP()

		iv, err := ip2int(ip)
		if err != nil {
			if lg := sserver.GetLogger(); lg != nil {
				lg.Warnf("Parse ip failed, ip: %s, error: %v.", ip, err)
			}
			sresp.AbortError(ctx, serror.ForbiddenIp)
			return
		}

		for _, r := range _ranges {
			if iv&r == r {
				return
			}
		}

		sresp.AbortError(ctx, serror.ForbiddenIp)
	}
}

func ip2int(ip string) (uint32, error) {
	ipv4 := net.ParseIP(ip).To4()
	if len(ipv4) == 0 {
		return 0, ErrInvalid
	}

	i := uint32(0)
	for _, b := range ipv4 {
		i = i << 8
		i |= uint32(b)
	}

	return i, nil
}

func parseRange(ip string) (uint32, error) {
	im := strings.Split(ip, "/")
	if len(im) == 0 || len(im) > 2 {
		return 0, ErrInvalid
	}

	mk := 32
	if len(im) == 2 {
		i, err := strconv.Atoi(im[1])
		if err != nil {
			return 0, ErrInvalid
		}
		ip = im[0]
		mk = i
	}

	iv, err := ip2int(ip)
	if err != nil {
		return 0, err
	}

	return iv & (uint32(0xffffffff) << (32 - mk)), nil
}
