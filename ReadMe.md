#Timeslot App

## Run App
1) Run `docker compose build`
2) Run `docker push wizkaley/go-timeslot-app:tag`
3) 

## Monitor
1) `minikube start`
note:
download dashboard 
- `helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard`
2) Run `kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443`

note:

create cluster role and role binding
- `kubectl create role pod-reader-all --verb=get --verb=list --verb=watch --resource=*`
- `kubectl create rolebinding default-pod-reader-all --role=pod-reader-all --serviceaccount=default:default --namespace=default`
3) Get the token for service account to login to dashboard `kubectl create token eshan`
4) Once in install the helm chart `cd timeslot-app/helm/timeslot-app/`with the command `helm install timeslot-app .`
5) If the chart is installed delete it first with `helm delete timeslot-app`
6) `minikube service --all` to check the IP's and ports of the running services.
7) See - http://127.0.0.1:56321/swagger/index.html#/ for API reference.