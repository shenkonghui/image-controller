# image-controller
The image-controller can solve the domestic problem that quay.io and gcr.io cannot be pulled. It can also solve problems such as cross-cluster deployment mirror synchronization.

## Scenes

1. gcr.io and quay.io cannot be accessed in China, you can use it to modify it as a domestic cache
2. In private cloud scenarios, it is often necessary to replace the image with a private address, which can be used to automate the replacement



## Installation

install crd、operator 、clusterrole、role and serveraccount.

```
kubectl apply -f deploy/crds/github.com_imageconfigs_crd.yaml

kubectl apply -f deploy/operator.yaml
kubectl apply -f deploy/admin.yaml
```



## Demo

A demo install a deployment use image quay.io/app-sre/nginx. 

Because China cannot access quay.io. Here replace the address in the image address with the domestic address of the University of Science and Technology of China.



####  install the ImageConfig cr

```
# cat example/quay.yaml
apiVersion: github.com/v1
kind: ImageConfig
metadata:
  name: quay
spec:
  repo: quay.io
  newRepo: quay.mirrors.ustc.edu.cn
  
# kubectl apply -f example/quay.yaml
```


#### install the deployment 
```
# cat example/nginx.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: quay.io/app-sre/nginx
          ports:
            - containerPort: 80
            
# kubectl apply -f example/nginx.yaml
```



View image changes in real time, you can see  image address "quay.io"  change to "quay.mirrors.ustc.edu.cn", That is success.

```
## print the images every 1s
~/src/github/image-controller#  while true;do kubectl get deployment  -o jsonpath='{.items[0].spec.template.spec.containers[0].image}' ; sleep 1; echo "" ;done
quay.io/app-sre/nginx
quay.io/app-sre/nginx
quay.mirrors.ustc.edu.cn/app-sre/nginx
quay.mirrors.ustc.edu.cn/app-sre/nginx
quay.mirrors.ustc.edu.cn/app-sre/nginx
quay.mirrors.ustc.edu.cn/app-sre/nginx
```

