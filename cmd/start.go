package cmd

import (
	"github.com/topicus-education-ops/boink/handler"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/retry"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts deployments and statefulSets",
	Long:  `This command starts kubernetes deployments and statefulSets`,
	Run: func(cmd *cobra.Command, args []string) {
		Clientset, _ = getClient()

		// Deployments
		deploymentClient := Clientset.AppsV1().Deployments(Namespace)
		deployments, err := getDeployments()
		if err != nil {
			panic(err)
		}
		for _, deployment := range deployments.Items {
			retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				handler.HandleDeployment(deployment, deploymentClient, "start")
				return nil
			})
			if retryErr != nil {
				panic(fmt.Errorf("Update failed: %v", retryErr))
			}

		}

		// StatefulSets
		statefulSetClient := Clientset.AppsV1().StatefulSets(Namespace)
		statefulSets, err := getStatefulSets()
		if err != nil {
			panic(err)
		}
		for _, statefulSet := range statefulSets.Items {
			retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				handler.HandleStatefulSet(statefulSet, statefulSetClient, "start")
				return nil
			})
			if retryErr != nil {
				panic(fmt.Errorf("Update failed: %v", retryErr))
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

}
