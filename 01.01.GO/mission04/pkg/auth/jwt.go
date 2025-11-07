package auth

import (
	"mission04/internal/config"
	"mission04/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	// RegisteredClaims 是 JWT 标准声明结构体，包含了 RFC 7519 规范中定义的预注册声明字段。
	// 这些声明提供了关于 JWT 的基本信息，包括签发者、过期时间、主题等标准属性。
	//
	// 该结构体嵌入到自定义的声明结构体中，用于创建符合标准的 JWT token。
	// 标准声明字段包括：
	// - Issuer (iss): 签发者标识
	// - Subject (sub): 主题标识
	// - Audience (aud): 接收者标识
	// - Expiration Time (exp): 过期时间戳
	// - Not Before (nbf): 生效时间戳
	// - Issued At (iat): 签发时间戳
	// - JWT ID (jti): 唯一标识符
	//
	// 使用时通常作为匿名字段嵌入到自定义 Claims 结构体中，例如：
	// type MyClaims struct {
	//     jwt.RegisteredClaims
	//     CustomField string
	// }

	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(user model.User) (string, error) {
	// 计算令牌过期时间
	// 从配置中获取令牌有效期，然后基于当前时间计算出具体的过期时间戳
	tokenExpiry := config.GetConfig().Auth.TokenExpiry
	expirationTime := time.Now().Add(time.Duration(tokenExpiry) * time.Second)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		// RegisteredClaims 是 JWT 标准声明结构体的初始化
		// 该结构体包含了 JWT 规范中定义的标准声明字段
		//
		// ExpiresAt: 指定令牌的过期时间，使用 Unix 时间戳格式
		//            通过 jwt.NewNumericDate 将 expirationTime 转换为标准的 NumericDate 格式
		//
		// 此处仅设置了过期时间，其他标准声明字段（如 Issuer、Subject、Audience 等）使用默认值
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// 生成JWT令牌并使用密钥进行签名
	//
	// 该函数创建一个新的JWT令牌，使用HS256签名方法和提供的声明信息，
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 然后使用配置中的JWT密钥对令牌进行签名并返回签名后的字符串。
	//
	// 返回值：
	//   - string: 签名后的JWT令牌字符串
	//   - error: 签名过程中可能发生的错误
	return token.SignedString([]byte(config.GetConfig().Auth.JwtSecret))
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	// parseJWTToken 解析JWT令牌字符串并验证其有效性
	// tokenString: 需要解析的JWT令牌字符串
	// claims: 用于存储解析后声明信息的结构体指针
	// 返回值: 解析后的token对象和可能的错误信息
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 提供JWT签名密钥的回调函数，用于验证令牌签名
		return []byte(config.GetConfig().Auth.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
