package handler

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

const targetReplicasAnnotation string = "applicationScaler.io/target-replicas"

/*
HandleDeployment - Verifies if the deployment contains the "intraday-enabled"annotation,
  it will scale the deployment to zero to stop
  and scale to targetReplica to start.
*/
func HandleDeployment(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface, action string) {
	if action == "stop" {
		scaleDeploymentToZero(deployment, deploymentClient)
	} else {
		scaleDeploymentUp(deployment, deploymentClient)
	}

}

func scaleDeploymentToZero(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface) {
	if *deployment.Spec.Replicas > int32(0) {
		logrus.Infof("Deployment (%s) has the annotation scaling to zero", deployment.ObjectMeta.Name)
		replicas := deployment.Spec.Replicas
		//Not all deployments have annotations
		if deployment.ObjectMeta.Annotations == nil {
			deployment.ObjectMeta.Annotations = make(map[string]string)
		}
		deployment.ObjectMeta.Annotations[targetReplicasAnnotation] = strconv.Itoa(int(*replicas))
		deployment.Spec.Replicas = int32Ptr(0)
		deploymentClient.Update(&deployment)
	} else {
		logrus.Infof("Deployment (%s) is already scaled to (0)", deployment.ObjectMeta.Name)
	}

}

func scaleDeploymentUp(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface) error {
	a := deployment.ObjectMeta.GetAnnotations()
	var replicas = 1
	var err error
	if *deployment.Spec.Replicas == int32(0) {
		if a[targetReplicasAnnotation] != "" {
			logrus.Infof("Deployment (%s) will be scaled up to (%s)", deployment.ObjectMeta.Name, a[targetReplicasAnnotation])
			replicas, err = strconv.Atoi(a[targetReplicasAnnotation])
			if err != nil {
				logrus.Error("Unable to convert replicas to number")
				panic(err)
			}
		}
		deployment.Spec.Replicas = int32Ptr(int32(replicas))
		_, err = deploymentClient.Update(&deployment)
		if err != nil {
			return err
		}

	} else {
		logrus.Infof("Deployment (%s) is already scaled up to (%d)", deployment.ObjectMeta.Name, *deployment.Spec.Replicas)
	}
	return nil
}

func HandleStatefulSet(statefulSet appsv1.StatefulSet, statefulSetClient v1.StatefulSetInterface, action string) {
	if action == "stop" {
		scaleToStatefulSetZero(statefulSet, statefulSetClient)
	} else {
		scaleStatefulSetUp(statefulSet, statefulSetClient)
	}

}

func scaleToStatefulSetZero(statefulSet appsv1.StatefulSet, statefulSetClient v1.StatefulSetInterface) {
	if *statefulSet.Spec.Replicas > int32(0) {
		logrus.Infof("StatefulSet (%s) has the annotation scaling to zero", statefulSet.ObjectMeta.Name)
		replicas := statefulSet.Spec.Replicas
		//Not all statefulSets have annotations
		if statefulSet.ObjectMeta.Annotations == nil {
			statefulSet.ObjectMeta.Annotations = make(map[string]string)
		}
		statefulSet.ObjectMeta.Annotations[targetReplicasAnnotation] = strconv.Itoa(int(*replicas))
		statefulSet.Spec.Replicas = int32Ptr(0)
		statefulSetClient.Update(&statefulSet)
	} else {
		logrus.Infof("StatefulSet (%s) is already scaled to (0)", statefulSet.ObjectMeta.Name)
	}

}

func scaleStatefulSetUp(statefulSet appsv1.StatefulSet, statefulSetClient v1.StatefulSetInterface) error {
	a := statefulSet.ObjectMeta.GetAnnotations()
	var replicas = 1
	var err error
	if *statefulSet.Spec.Replicas == int32(0) {
		if a[targetReplicasAnnotation] != "" {
			logrus.Infof("StatefulSet (%s) will be scaled up to (%s)", statefulSet.ObjectMeta.Name, a[targetReplicasAnnotation])
			replicas, err = strconv.Atoi(a[targetReplicasAnnotation])
			if err != nil {
				logrus.Error("Unable to convert replicas to number")
				panic(err)
			}
		}
		statefulSet.Spec.Replicas = int32Ptr(int32(replicas))
		_, err = statefulSetClient.Update(&statefulSet)
		if err != nil {
			return err
		}

	} else {
		logrus.Infof("StatefulSet (%s) is already scaled up to (%d)", statefulSet.ObjectMeta.Name, *statefulSet.Spec.Replicas)
	}
	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
