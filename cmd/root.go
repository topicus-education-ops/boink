package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var rootCmd = &cobra.Command{
	Use:   "applicationScaler",
	Short: "ApplicationScaler is a helper command to stop and start Kubernetes Deployments and statefulSets.",
	Long: `ApplicationScaler is a helper command to stop and start Kubernetes Deployments and statefulSets.
	       It can remember previous scale settings prior to stopping.`,
}

var Clientset *kubernetes.Clientset

var Namespace string

var PathToConfig string

var Selectors string

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true})

	// Output to stdout instead of the default stderr, could also be a file.
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	rootCmd.PersistentFlags().StringVarP(&PathToConfig, "config", "c", "", "config file (default is $KUBECONFIG)")

	rootCmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "default", "The namespace where boink will look for the deployments and statefulSets.")

	rootCmd.PersistentFlags().StringVarP(&Selectors, "label", "l", "default", "The filter deployments and statefulSets based on kubernetes selector")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getClient() (*kubernetes.Clientset, error) {
	logrus.Info("Calling getClient()")
	var config *rest.Config
	var err error
	if PathToConfig == "" {
		logrus.Info("Using in cluster config")
		config, err = rest.InClusterConfig()
		// in cluster access
	} else {
		logrus.Info("Using out of cluster config")
		config, err = clientcmd.BuildConfigFromFlags("", PathToConfig)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getDeployments() (*v1.DeploymentList, error) {
	deploymentClient := Clientset.AppsV1().Deployments(Namespace)
	var listOptions metav1.ListOptions

	if Selectors != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: Selectors,
			Limit:         100,
		}

	} else {
		listOptions = metav1.ListOptions{}
	}
	deployments, err := deploymentClient.List(listOptions)
	if err != nil {
		logrus.Warnf("Failed to find deployments: %v", err)
		return nil, err
	}

	return deployments, nil
}

func getStatefulSets() (*v1.StatefulSetList, error) {
	statefulSetClient := Clientset.AppsV1().StatefulSets(Namespace)
	var listOptions metav1.ListOptions

	if Selectors != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: Selectors,
			Limit:         100,
		}

	} else {
		listOptions = metav1.ListOptions{}
	}
	statefulSets, err := statefulSetClient.List(listOptions)
	if err != nil {
		logrus.Warnf("Failed to find statefulSets: %v", err)
		return nil, err
	}

	return statefulSets, nil
}
