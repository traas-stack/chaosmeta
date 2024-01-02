package initwebhook

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/go-logr/logr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	certificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

const (
	certKey                     = "tls.crt"
	keyKey                      = "tls.key"
	secretName                  = "chaosmeta-%s-webhook-server-cert"
	csrName                     = "chaosmeta-%s-csr"
	validatingWebhookConfigName = "chaosmeta-%s-validating-webhook-configuration"
	mutatingWebhookConfigName   = "chaosmeta-%s-mutating-webhook-configuration"
	webhookServiceName          = "chaosmeta-%s-webhook-service.%s.svc"
	defaultNamespace            = "chaosmeta"
	namespaceEnv                = "DEFAULTNAMESPACE"
	secretPath                  = "/tmp/k8s-webhook-server/serving-certs/"
)

func isExistAndValid(client client.Client, component, curNamespace string) (*v1.Secret, bool) {
	secret := &v1.Secret{}
	err := client.Get(context.Background(), types.NamespacedName{Name: fmt.Sprintf(secretName, component), Namespace: curNamespace}, secret)
	if err == nil {
		// is it expired
		cert, err1 := tls.X509KeyPair(secret.Data[certKey], secret.Data[keyKey])
		if err1 != nil {
			return nil, false
		}
		x509Cert, err1 := x509.ParseCertificate(cert.Certificate[0])
		if err1 != nil {
			return secret, false
		}
		now := time.Now()
		if now.After(x509Cert.NotAfter) || now.Before(x509Cert.NotBefore) {
			return secret, false
		}
		return secret, true
	}
	return nil, false
}

func InitCert(log logr.Logger, component string) error {
	curNamespace := os.Getenv(namespaceEnv)
	if curNamespace == "" {
		curNamespace = defaultNamespace
	}
	clientSet, err := kubernetes.NewForConfig(config.GetConfigOrDie())
	if err != nil {
		fmt.Println("failed to create clientSet")
		return err
	}
	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		fmt.Println("failed to create client")
		return err
	}

	if oldSecret, valid := isExistAndValid(cl, component, curNamespace); valid == true {
		err = saveSecretToFile(log, oldSecret.Data[certKey], oldSecret.Data[keyKey])
		if err != nil {
			return err
		}
		err = updateWebhookConfig(log, cl, oldSecret.Data[certKey], component)
		if err != nil {
			return err
		}
		return nil
	}
	// CA private key
	privateKey, err := rsa.GenerateKey(cryptorand.Reader, 4096)
	if err != nil {
		fmt.Println(err)
		return err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})

	subj := pkix.Name{
		CommonName:   fmt.Sprintf("system:node:%s.%s.svc", "chaosmeta", curNamespace),
		Organization: []string{"system:nodes"},
	}
	dnsNames := []string{fmt.Sprintf(webhookServiceName, component, curNamespace)}

	// CSR generation
	csrTemplate := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
		DNSNames:           dnsNames,
	}

	csrBytes, err := x509.CreateCertificateRequest(cryptorand.Reader, &csrTemplate, privateKey)
	if err != nil {
		return err
	}

	// PEM encode the server cert and key
	serverCertPEM := new(bytes.Buffer)
	_ = pem.Encode(serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: csrBytes,
	})

	// delete original csr if there is
	originalCsr := &certificatesv1.CertificateSigningRequest{}
	err = cl.Get(context.Background(), types.NamespacedName{Name: fmt.Sprintf(csrName, component), Namespace: curNamespace}, originalCsr)
	if err == nil {
		err1 := cl.Delete(context.Background(), originalCsr)
		if err1 != nil {
			return err1
		}
	}
	csr := &certificatesv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf(csrName, component),
			Namespace: curNamespace,
		},
		Spec: certificatesv1.CertificateSigningRequestSpec{
			Request:    pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}),
			SignerName: "kubernetes.io/kubelet-serving",
			Usages: []certificatesv1.KeyUsage{
				certificatesv1.UsageDigitalSignature,
				certificatesv1.UsageKeyEncipherment,
				certificatesv1.UsageServerAuth,
			},
			Groups:   []string{"system:nodes"},
			Username: "system:nodes:chaosmeta",
		},
	}
	err = cl.Create(context.Background(), csr)
	if err != nil {
		log.Error(err, "create csr failed")
		return err
	}
	csr.Status.Conditions = append(csr.Status.Conditions, certificatesv1.CertificateSigningRequestCondition{
		Type:           certificatesv1.CertificateApproved,
		Status:         corev1.ConditionTrue,
		Reason:         "ChaosmetaApprove",
		Message:        "approve by chaosmeta",
		LastUpdateTime: metav1.Now(),
	})
	// approve csr
	_, err = clientSet.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.Background(), csr.Name, csr, metav1.UpdateOptions{})
	if err != nil {
		log.Error(err, "approve csr failed")
		return err
	}
	// wait for csr's Status.Certificate
	time.Sleep(time.Second * 3)
	newCsr := &certificatesv1.CertificateSigningRequest{}
	err = cl.Get(context.Background(), types.NamespacedName{Name: fmt.Sprintf(csrName, component), Namespace: curNamespace}, newCsr)
	if err != nil {
		log.Error(err, "get csr failed")
		return err
	}

	if newCsr.Status.Certificate == nil || len(newCsr.Status.Certificate) == 0 {
		log.Error(err, "csr's Status.Certificate is nil ")
		return err
	}
	secret := &v1.Secret{
		Type: v1.SecretTypeTLS,
		ObjectMeta: v12.ObjectMeta{
			Name:      fmt.Sprintf(secretName, component),
			Namespace: curNamespace,
		},
		Data: map[string][]byte{
			certKey: newCsr.Status.Certificate,
			keyKey:  keyPem,
		}}
	// remove old secret first
	oldSecret := &v1.Secret{}
	secretIndex := types.NamespacedName{Namespace: curNamespace, Name: fmt.Sprintf(secretName, component)}
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
	err = saveSecretToFile(log, newCsr.Status.Certificate, keyPem)
	if err != nil {
		return err
	}
	err = updateWebhookConfig(log, cl, newCsr.Status.Certificate, component)
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

func updateWebhookConfig(log logr.Logger, cl client.Client, serverCertBytes []byte, component string) error {
	mutatingWebhookConfig := &admissionregistrationv1.MutatingWebhookConfiguration{}
	err := cl.Get(context.Background(), types.NamespacedName{Name: fmt.Sprintf(mutatingWebhookConfigName, component)}, mutatingWebhookConfig)
	if err != nil {
		log.Error(err, "failed to get mutatingWebhookConfig")
		return err
	}

	mutatingWebhookConfig.Webhooks[0].ClientConfig.CABundle = serverCertBytes
	err = cl.Update(context.Background(), mutatingWebhookConfig)
	if err != nil {
		log.Error(err, "failed to get mutatingWebhookConfig")
		return err
	}

	validatingWebhookConfig := &admissionregistrationv1.ValidatingWebhookConfiguration{}
	err = cl.Get(context.Background(), types.NamespacedName{Name: fmt.Sprintf(validatingWebhookConfigName, component)}, validatingWebhookConfig)
	if err != nil {
		log.Error(err, "failed to get mutatingWebhookConfig")
		return err
	}

	validatingWebhookConfig.Webhooks[0].ClientConfig.CABundle = serverCertBytes
	err = cl.Update(context.Background(), validatingWebhookConfig)
	if err != nil {
		log.Error(err, "failed to get validatingWebhookConfig")
		return err
	}

	return nil
}
