package clientset

import (
	"chaosmeta-platform/util/log"

	"context"
	"k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"time"
)

type KubeLoadModeStr string

func (k KubeLoadModeStr) IsEmpty() bool {
	return k == ""
}

func (k KubeLoadModeStr) String() string {
	return string(k)
}

//func GetKubeRestConf(mode KubeLoadMode, cfgInfo string) (rfg *rest.Config, err error) {
//	switch mode {
//	case KubeLoadFromLocal:
//		rfg, err = clientcmd.BuildConfigFromFlags("", cfgInfo)
//	case KubeLoadInCluster:
//		rfg, err = rest.InClusterConfig()
//	case KubeLoadFromBase64Stream:
//		var bts []byte
//		bts, err = base64.StdEncoding.DecodeString(cfgInfo)
//		if err != nil {
//			break
//		}
//		var ccfg clientcmd.ClientConfig
//		ccfg, err = clientcmd.NewClientConfigFromBytes(bts)
//		if err != nil {
//			break
//		}
//		rfg, err = ccfg.ClientConfig()
//	default:
//		err = fmt.Errorf("invalid kube config load mode: %s or empty", mode)
//	}
//	if err != nil {
//		return
//	}
//
//	return
//}

func GetLoadFile(loadFile string) string {
	homeDir := "~"
	if strings.HasPrefix(loadFile, homeDir) {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Error(err)
		}
		return strings.Replace(loadFile, homeDir, dirname, 1)
	}
	return loadFile
}

type Operator struct {
	opClient  versioned.Interface
	clientset Interface
}

func Init() error {
	kubeCfg, err := GetKubeRestConf("LoadFromLocal", "/Users")
	if err != nil {
		return err
	}

	opClient, err := versioned.NewForConfig(kubeCfg)
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		return err

	}

	kubeInformerFactory := informers.NewSharedInformerFactory(kubeClient, time.Second*30)
	cs := NewClientset(opClient)

	operator := NewOperator(
		opClient,
		cs,
	)

	run := func(ctx context.Context) {
		kubeInformerFactory.Start(ctx.Done())
		if err := operator.Run(ctx); err != nil {
			log.Panic(err)
		}
	}

	run(context.Background())
	return nil
}

func NewOperator(
	opClient versioned.Interface,
	clientset Interface,
) *Operator {

	return &Operator{
		opClient:  opClient,
		clientset: clientset,
	}
}

func (o *Operator) Run(ctx context.Context) error {
	defer runtime.HandleCrash()
	// wait cluster synced
	go o.clientset.RunWorker()
	// start controller
	go o.clientset.Run(ctx, false)

	select {}
}
