package initwebhook

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/go-logr/logr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"math/big"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

const (
	certKey          = "tls.crt"
	keyKey           = "tls.key"
	injectSecretName = "chaosmeta-inject-webhook-server-cert"
	defaultNamespace = "chaosmeta"
	namespaceEnv     = "DEFAULTNAMESPACE"
	secretPath       = "/tmp/k8s-webhook-server/serving-certs/"
)

func isExistAndValid(client client.Client) (*v1.Secret, bool) {
	secret := &v1.Secret{}
	err := client.Get(context.Background(), types.NamespacedName{Name: injectSecretName, Namespace: defaultNamespace}, secret)
	if err == nil {
		// is it expired
		cert, err1 := x509.ParseCertificate(secret.Data[certKey])
		if err1 != nil {
			return secret, false
		}
		now := time.Now()
		if now.After(cert.NotAfter) || now.Before(cert.NotBefore) {
			return secret, false
		}
		return secret, true
	}
	return nil, false
}

func InitCert(log logr.Logger) error {
	curNamespace := os.Getenv(namespaceEnv)
	if curNamespace == "" {
		curNamespace = defaultNamespace
	}
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		return err
	}

	if oldSecret, valid := isExistAndValid(cl); valid == true {
		err = saveSecretToFile(log, oldSecret.Data[certKey], oldSecret.Data[keyKey])
		if err != nil {
			return err
		}
		err = updateWebhookConfig(log, cl, oldSecret.Data[certKey])
		if err != nil {
			return err
		}
		return nil
	}
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Organization: []string{"chaosmeta.io"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// CA private key
	caPrivateKey, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Self signed CA certificate
	caBytes, err := x509.CreateCertificate(cryptorand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// PEM encode CA cert
	caPEM := new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	dnsNames := []string{"chaosmeta-inject-webhook-service." + defaultNamespace + ".svc"}

	// server cert config
	cert := &x509.Certificate{
		DNSNames:     dnsNames,
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"chaosmeta.io"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	// server private key
	serverPrivateKey, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// sign the server cert
	serverCertBytes, err := x509.CreateCertificate(cryptorand.Reader, cert, ca, &serverPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// PEM encode the server cert and key
	serverCertPEM := new(bytes.Buffer)
	_ = pem.Encode(serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})

	serverPrivateKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(serverPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivateKey),
	})

	secret := &v1.Secret{
		Type: v1.SecretTypeTLS,
		ObjectMeta: v12.ObjectMeta{
			Name:      "chaosmeta-inject-webhook-server-cert",
			Namespace: curNamespace,
		},
		Data: map[string][]byte{
			certKey: serverCertPEM.Bytes(),
			keyKey:  serverPrivateKeyPEM.Bytes(),
		}}
	// remove first
	oldSecret := &v1.Secret{}
	secretIndex := types.NamespacedName{Namespace: curNamespace, Name: "chaosmeta-inject-webhook-server-cert"}
	if err = cl.Get(context.Background(), secretIndex, oldSecret); err == nil {
		err = cl.Delete(context.Background(), oldSecret)
		if err != nil {
			return err
		}
	}
	err = cl.Create(context.Background(), secret)
	if err != nil {
		log.Error(err, "create secret failed")
		return err
	}
	err = saveSecretToFile(log, serverCertPEM.Bytes(), serverPrivateKeyPEM.Bytes())
	if err != nil {
		return err
	}
	err = updateWebhookConfig(log, cl, serverCertPEM.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func saveSecretToFile(log logr.Logger, serverCertBytes []byte, serverPrivateKeyBytes []byte) error {
	err := os.MkdirAll(secretPath, 755)
	if err != nil {
		log.Error(err, "create secret dir failed")
		return err
	}
	err = os.WriteFile(secretPath+certKey, serverCertBytes, 0600)
	if err != nil {
		log.Error(err, "create secret file failed")
		return err
	}

	err = os.WriteFile(secretPath+keyKey, serverPrivateKeyBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func updateWebhookConfig(log logr.Logger, cl client.Client, serverCertBytes []byte) error {
	mutatingWebhookConfig := &admissionregistrationv1.MutatingWebhookConfiguration{}
	err := cl.Get(context.Background(), types.NamespacedName{Name: "chaosmeta-inject-mutating-webhook-configuration"}, mutatingWebhookConfig)
	if err != nil {
		log.Error(err, "failed to get mutatingWebhookConfig")
		return err
	}
	if mutatingWebhookConfig.Webhooks[0].ClientConfig.CABundle == nil {
		mutatingWebhookConfig.Webhooks[0].ClientConfig.CABundle = serverCertBytes
		err = cl.Update(context.Background(), mutatingWebhookConfig)
		if err != nil {
			log.Error(err, "failed to get mutatingWebhookConfig")
			return err
		}
	}

	validatingWebhookConfig := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	err = cl.Get(context.Background(), types.NamespacedName{Name: "chaosmeta-inject-validating-webhook-configuration"}, validatingWebhookConfig)
	if err != nil {
		log.Error(err, "failed to get mutatingWebhookConfig")
		return err
	}
	if validatingWebhookConfig.Webhooks[0].ClientConfig.CABundle == nil {
		validatingWebhookConfig.Webhooks[0].ClientConfig.CABundle = serverCertBytes
		err = cl.Update(context.Background(), validatingWebhookConfig)
		if err != nil {
			log.Error(err, "failed to get validatingWebhookConfig")
			return err
		}
	}
	return nil
}
