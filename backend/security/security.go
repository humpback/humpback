package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"humpback/pkg/utils"
	"log/slog"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// 证书相关常量
	caCertValidity   = 365 * 24 * time.Hour // 1年
	certValidity     = 90 * 24 * time.Hour  // 3个月
	certOrganization = "Humpback"

	// JWT相关常量
	tokenValidity    = 24 * time.Hour // 24小时
	tokenRefreshTime = 1 * time.Hour  // 在token还有1小时过期时刷新
)

// JWTClaims 自定义的JWT声明
type JWTClaims struct {
	WorkerID string `json:"wid"`
	jwt.RegisteredClaims
}

// CertificateBundle 包含完整的证书和密钥
type CertificateBundle struct {
	Cert     *x509.Certificate
	PrivKey  *ecdsa.PrivateKey
	CertPool *x509.CertPool // CA证书池
	CertPEM  []byte         // PEM编码的证书
	KeyPEM   []byte         // PEM编码的私钥
}

// SecurityManager 安全管理器
type SecurityManager struct {
	caCert       *x509.Certificate
	caPrivateKey *ecdsa.PrivateKey
	jwtSecret    []byte
	cacheFolder  string
}

var sm *SecurityManager

// NewSecurityManager 创建安全管理器实例
func InitSecurityManager(cacheFolder string) error {
	sm = &SecurityManager{}

	sm.cacheFolder = cacheFolder
	// 生成安全的JWT密钥
	sm.jwtSecret = make([]byte, 32) // 256-bit key
	if _, err := rand.Read(sm.jwtSecret); err != nil {
		return fmt.Errorf("failed to generate JWT secret: %w", err)
	}

	return nil
}

func GenerateWebsiteCert(certPath, keyPath string) (*CertificateBundle, error) {

	if utils.FileExist(certPath) && utils.FileExist(keyPath) {
		ca, key, err := LoadCertificateAndKey(certPath, keyPath)
		if err == nil {
			slog.Info("[Cert] Certs already exist, skip generating new certs.")

			// PEM编码
			certPEM := pem.EncodeToMemory(&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: ca.Raw,
			})

			keyBytes, err := x509.MarshalECPrivateKey(key)
			if err != nil {
				return nil, err
			}

			keyPEM := pem.EncodeToMemory(&pem.Block{
				Type:  "EC PRIVATE KEY",
				Bytes: keyBytes,
			})

			return &CertificateBundle{
				Cert:    ca,
				PrivKey: key,
				CertPEM: certPEM,
				KeyPEM:  keyPEM,
			}, nil
		}
	}

	return CreateCertificateBundle("humpback-website")
}

// 从文件加载证书和私钥
func LoadCertificateAndKey(certFile, keyFile string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	// 读取证书文件
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return nil, nil, fmt.Errorf("read CA cert file failed: %v", err)
	}

	// 读取私钥文件
	keyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("read CA key failed: %v", err)
	}

	// 解析证书
	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return nil, nil, fmt.Errorf("invalid CA cert PEM data")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CA cert failed: %v", err)
	}

	// 解析私钥
	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, nil, fmt.Errorf("invalid CA key PEM data")
	}

	// 尝试解析为ECDSA私钥
	privKey, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err == nil {
		return cert, privKey, nil
	}

	return nil, nil, fmt.Errorf("invalid CA key")
}

func generateFileName(commonName string) (string, string) {
	// 生成证书和私钥的文件名
	certFile := fmt.Sprintf("%s/%s.crt", sm.cacheFolder, commonName)
	keyFile := fmt.Sprintf("%s/%s.key", sm.cacheFolder, commonName)
	return certFile, keyFile
}

// GenerateCA 生成根CA
func GenerateCA() error {

	slog.Info("cache folder", "path", sm.cacheFolder)

	commonName := "Humpback-Root-CA"
	caFile, caKey := generateFileName(commonName)

	if utils.FileExist(caFile) && utils.FileExist(caKey) {
		caCert, caKey, err := LoadCertificateAndKey(caFile, caKey)
		if err == nil {
			slog.Info("[Cert] CA Certs already exist, skip generating new certs.")
			sm.caCert = caCert
			sm.caPrivateKey = caKey
			return nil
		}
	}

	// 生成私钥
	caPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	// 创建证书模板
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{certOrganization},
			CommonName:   commonName,
		},
		NotBefore:             time.Now().Add(-10 * time.Minute), // 10分钟前开始生效
		NotAfter:              time.Now().Add(caCertValidity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}

	// 自签名生成CA证书
	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return err
	}

	// 解析证书对象
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return err
	}

	sm.caCert = caCert
	sm.caPrivateKey = caPrivateKey

	SaveToFile(caCert, caPrivateKey, caFile, caKey)

	return nil
}

func GetRootCA() *x509.Certificate {
	if sm == nil || sm.caCert == nil {
		return nil
	}
	return sm.caCert
}

// 保存证书到文件（可选）
func SaveToFile(cert *x509.Certificate, privKey *ecdsa.PrivateKey, certPath, keyPath string) error {
	// 保存证书
	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	certBlock := &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}
	if err := pem.Encode(certFile, certBlock); err != nil {
		return err
	}

	// 保存私钥
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	keyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return err
	}

	keyBlock := &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}
	if err := pem.Encode(keyFile, keyBlock); err != nil {
		return err
	}

	return nil
}

// CreateCertificateBundle 为服务创建证书包
func CreateCertificateBundle(commonName string) (*CertificateBundle, error) {
	if sm.caCert == nil || sm.caPrivateKey == nil {
		return nil, errors.New("CA not initialized")
	}

	certFile, certKey := generateFileName(commonName)

	var cert *x509.Certificate
	var privateKey *ecdsa.PrivateKey
	var err error

	if utils.FileExist(certFile) && utils.FileExist(certKey) {
		cert, privateKey, err = LoadCertificateAndKey(certFile, certKey)
		if err == nil {
			slog.Info("[Cert] Certs already exist, skip generating new certs.", "commonName", commonName)
		}
	} else {
		// 生成私钥
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}

		// 创建CSR模板
		csrTemplate := &x509.Certificate{
			SerialNumber: big.NewInt(time.Now().UnixNano()),
			Subject: pkix.Name{
				CommonName:   commonName,
				Organization: []string{certOrganization},
			},
			NotBefore:   time.Now(),
			NotAfter:    time.Now().Add(certValidity),
			KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		}

		// 使用CA签名生成证书
		certBytes, err := x509.CreateCertificate(rand.Reader, csrTemplate, sm.caCert, &privateKey.PublicKey, sm.caPrivateKey)
		if err != nil {
			return nil, err
		}

		// 解析生成的证书
		cert, err = x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, err
		}

		SaveToFile(cert, privateKey, certFile, certKey)
	}

	// 创建证书池
	certPool := x509.NewCertPool()
	certPool.AddCert(sm.caCert)

	// PEM编码
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})

	return &CertificateBundle{
		Cert:     cert,
		PrivKey:  privateKey,
		CertPool: certPool,
		CertPEM:  certPEM,
		KeyPEM:   keyPEM,
	}, nil
}

// CreateTLSConfig 从证书包创建TLS配置
func (bundle *CertificateBundle) CreateTLSConfig(isWebsite bool) *tls.Config {
	cert, _ := tls.X509KeyPair(bundle.CertPEM, bundle.KeyPEM)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      bundle.CertPool,
		// ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs: bundle.CertPool,
	}
	if isWebsite {
		config.ClientAuth = tls.NoClientCert
	} else {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return config
}

// GenerateWorkerToken 为worker生成JWT token
func GenerateWorkerToken(workerID string) (string, error) {
	claims := JWTClaims{
		WorkerID: workerID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "humpback-server",
			Subject:   "worker-auth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenValidity)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(sm.jwtSecret)
}

// VerifyWorkerToken 验证worker token
func VerifyWorkerToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return sm.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
