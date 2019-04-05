# ApplicationScaler
ApplicationScaler (Boink) is a simple Go client application that can handle stopping and starting Kubernetes `Deployments` and `StatefulSets`.
It works by selecting `Deployments` and `StatefulSets` based on labels.  It can also remember the previous known replicas, unlike a standard `kubectl scale` command where you need to specify the replicas manually.

This tool can be helpful when you have certain applications which needs to be stopped during certain period of time.
Pair this tool with kubernetes `CronJob` to automatically stop or start a `Deployment` or `StatefulSet`


## How to run the application

1. Outside of the cluster - This is a normal executable program to interact with kubernetes cluster.  You can use minikube.

    Arguments:
    - --config - The location of $KUBECONFIG
    - --namespace - The namespace to use.
    - --label - Specify the selectors.
    - --action - Action to do, possible values `stop`, `start`

    To Stop:

       applicationScaler --config $KUBECONFIG --namespace test --label app=nginx --action stop

    To Start:
        
       applicationScaler --config $KUBECONFIG --namespace test --label app=nginx --action start

2.  In cluster

    Make sure you have the right permission in the cluster.  See samples in `manifest/` folder.
    
    1. Create the `ServiceAccount`

    ```
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: applicationScaler
      namespace: test
    ```

    2. Create the `Role`

    ```
    apiVersion: rbac.authorization.k8s.io/v1
    kind: Role
    metadata:
      name: applicationScaler
      namespace: test
    rules:
    - apiGroups: ["extensions","apps"]
      resources: ["deployments","statefulsets"]
      verbs: ["get", "list", "update"]

    ```

    3.  Bind the Role with the service account

    ```
    apiVersion: rbac.authorization.k8s.io/v1
    kind: RoleBinding
    metadata:
        name: applicationScaler 
        namespace: test       
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: Role
      name: applicationScaler
    subjects:
    - kind: ServiceAccount
      name: applicationScaler
      namespace: test
    ```

    4.  Create the `CronJob`, Stop `Deployment` or `StatefulSet` with label `app=nginx`.

    ```
    apiVersion: batch/v1beta1
    kind: CronJob
    metadata:
      name: nginx-starter
      namespace: test
    spec:
      #this is in UTC  will run at 6:46 AM everyday
      schedule: "46 6 * * *"
      startingDeadlineSeconds: 10
      concurrencyPolicy: Forbid
      jobTemplate:
        spec:      
          template:
            spec:
              serviceAccountName: applicationScaler
              containers:
              - name: applicationScaler
                image: applicationScaler:1.0
                command: ["/applicationScaler"]
                args: ["--namespace","test", "--label", "app=nginx", "--action" , "start"]
              restartPolicy: OnFailure
    ```



### Tools:
1. command line library using `github.com/urfave:v1.18.0`.  This library makes it easy to create a command line application in `Go`.

2.  Uses client-go to interact with kubernetes cluster.

3.  `github.com/Sirupsen/logrus` - Standard logging mechanism for `Go`.

Check the `go.mod` to see all the dependencies.

## To build and run 
1. Make sure you enable `$GO111MODULE` to `on`
2. Go to the working directory `$GOPATH/src/boink/`.
3. Update all the dependencies using `dep ensure` to update all the dependencies.
4. Do `go build`
5. Do `go test ./... -cover` to run unit test with code coverage.
6. Finally to run the application `applicationScaler --config &KUBECONFIG --namespace test --label app=nginx --action stop`.  


If you are using skaffold, there is `skaffold.yaml` included at the root of the project.  Simply do a `skaffold dev` from the `$GOPATH/src/boink` and you are good to go.

