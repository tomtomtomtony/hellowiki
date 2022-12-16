package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"hellowiki/api/result"
	"hellowiki/config"
	"log"
	"strings"
	"time"
)

type JWT struct {
	SignKey    []byte        // 秘钥
	MaxRefresh time.Duration // 刷新token的最大过期时间
}

// CustomJWTClaims 自定义Payload信息
type CustomJWTClaims struct {
	userName     string
	ExpireAtTime int64 // 过期时间

	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwt.RegisteredClaims
}

func NewJWT() *JWT {
	maxRefreshTime := config.Cfg.JwtConfig.MaxRefreshTime

	return &JWT{
		SignKey:    []byte(config.Cfg.JwtConfig.SecretKey),
		MaxRefresh: maxRefreshTime,
	}
}

// ParserToken 解析 Token，中间件中调用
func (j *JWT) ParserToken(ctx *gin.Context) (*CustomJWTClaims, int) {
	// 从header获取token
	tokenString, parseErr := j.getTokenFromHeader(ctx)
	if parseErr != result.SUCCSE {
		return nil, parseErr
	}

	// 解析token
	token, err := j.parseTokenString(tokenString)
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok {
			if validationErr.Errors == jwt.ValidationErrorMalformed {
				return nil, result.ERROR_TOKEN_MALFORMED
			} else if validationErr.Errors == jwt.ValidationErrorExpired {
				return nil, result.ERROR_TOKEN_EXPIRED
			}
		}
		return nil, result.ERROR_TOKEN_INVALID
	}

	if claims, ok := token.Claims.(*CustomJWTClaims); ok && token.Valid {
		return claims, result.SUCCSE
	}
	return nil, result.ERROR_TOKEN_INVALID
}

// RefreshToken 更新 Token，用以提供 refresh token 接口
func (j *JWT) RefreshToken(ctx *gin.Context) (string, int) {
	// 1. 从 Header 里获取 token
	tokenString, parseErr := j.getTokenFromHeader(ctx)
	if parseErr != result.SUCCSE {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := j.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		// 满足 refresh 的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return "", result.ERROR
		}
	}

	// 4. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*CustomJWTClaims)

	// 5. 检查是否过了『最大允许刷新的时间』
	t := time.Now().Add(-j.MaxRefresh).Unix()
	// 首次签名时间 > (当前时间 - 最大允许刷新时间)
	if claims.IssuedAt.Unix() > t {
		claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(j.expireAtTime())
		res, err := j.createToken(*claims)
		if err != nil {
			return "", result.ERROR
		}
		return res, result.SUCCSE
	}

	return "", result.ERROR_TOKEN_EXPIRED_MaxRefresh
}

// IssueToken 生成  Token，在登录成功时调用
func (j *JWT) IssueToken(userName string) string {
	// 构造自定义Payload信息
	expireTime := j.expireAtTime()
	claims := CustomJWTClaims{
		// 用户信息
		userName: userName,
		// 过期时间
		ExpireAtTime: expireTime.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now()), // 签名生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 首次签名时间（后续刷新 Token 不会更新）
			ExpiresAt: jwt.NewNumericDate(expireTime), // 签名过期时间
			Issuer:    config.Cfg.JwtConfig.Issuer,    // 签名颁发者
		},
	}

	// 根据 claims 生成token对象
	token, err := j.createToken(claims)
	if err != nil {
		log.Fatalf("未能生成token:{%v}", err)
		return ""
	}
	return token
}

// createToken 创建 Token，内部使用，外部请调用 IssueToken
func (j *JWT) createToken(claims CustomJWTClaims) (string, error) {
	// 使用HS256算法进行token生成
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.SignKey)
}

// token过期时间
func (j *JWT) expireAtTime() time.Time {
	expireTime := config.Cfg.JwtConfig.TokenExpireDuration
	return time.Now().Add(expireTime)
}

// getTokenFromHeader 使用 ParseWithClaims 解析 Token
// Authorization:Bearer xxxxx
func (j *JWT) getTokenFromHeader(ctx *gin.Context) (string, int) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", result.ERROR_HEADER_EMPTY
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", result.ERROR_HEADER_MALFORMED
	}
	return parts[1], result.SUCCSE
}

// parseTokenString 解析token
func (j *JWT) parseTokenString(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JwtConfig.SecretKey), nil
	})
}
