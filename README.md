# 背景
kubernetes（k8s）集群的网络状况有时会出现问题。应用容器间有时访问不通，而业务方不确定具体原因，只能找运维人员。
运维人员可能得逐一排查。为减轻运维人员的工作负担，使得定位问题更快速，且能够一键获取集群
网络连通性的情况，实时了解集群网络状况，使得运维更高效，所以有了这个工具。


# 简介
本项目是一个网络检测工具，旨在检测k8s集群的网络连通性。工具名为NetChecker

### 指标：
· 覆盖集群的每个节点
· 测试pod->pod的网络
· 测试pod->node的网络
· 测试pod->service的网络

备注：pod->pod或pod->node不通可视作node->node不通

### 项目结构：

![](https://github.com/yingyan003/netChecker/blob/master/picture/projStruct.png)

### 项目说明：

该项目有3个独立的子模块：netChecker、ping、pong。
netChecker：负责定时去k8s集群获取ping对应的所有pod列表，所有work节点列表， 和pong对应的service（只建一个）。
ping：定时获取netChecker收集的信息。并利用这些信息进行集群网络检测。在pod->pod/node/service类型的检测中，
发出检测行为的源pod，是部署ping对应的pod。而被检测对象有pod/node/service，这里的pod是部署pong对应
pod，node是集群的每个work节点，service是pond对应的service（只部署了一个）。
pong：在pod->pod类型的检测中起作用，仅供部署ping对应的pod来调。

备注：与k8s交互的只有netChecker。

### 部署方式
代码开发完后会上传到gitlab中，gitlab的pipeline编译时会根据每个模块下的Dockerfile文件将每个模块打成镜像，并推到
指定的镜像仓库（当然这其中的配置这里就不细说了）。镜像生成后，就可以用来创建容器了。

netChecker：
使用deployment资源对象部署。只需部署一个，副本数只能有1个，否则频道中的消息就会重复

ping/pong：
使用k8s的资源对象daemonset（保证每个work节点只有一个pod的）部署。保证每个节点上仅有一个
ping和一个pong对应的pod。实现节点全覆盖网络检测。

pod对应的service：
这里一共部署了2个。
· 因为检测类型涉及了pod->service。这里有pong部署了一个service（只部署了一个）。
· 又因为web访问netChecker，请求检查结果时，是集群外访问集群内的方式，需要给netChecker部署一个nodePort类型的service。

# 实现方案

备注：
方案一对应的代码在分支"no-pubsub"
方案二对应的代码在分支"master"

## 方案一

![](https://github.com/yingyan003/netChecker/blob/master/picture/model1.png)

该方案存在一个问题，就是1和3的定时如何设置合理。比如1设置每59s执行一次，3设置为没60s执行一次。
留了1s的时间差等待1将数据写入后，3再执行。虽说1获取的时间一般来说比较快，但是这种方式并不那么
优雅的保证3就一定能够在1写入完成后获取到数据。比如1访问k8s集群由于网络原因比较慢等原因。
鉴于这种情况，决定采用redis的pub/sub消息发布和订阅模式实现。于是就有了方案二。

## 方案二
![](https://github.com/yingyan003/netChecker/blob/master/picture/model2.png)

pub/sub具体实现方案

![](https://github.com/yingyan003/netChecker/blob/master/picture/pubsub.png)

### 问题：
该方案有个点需要注意，就是ping检查完成后，将结果通过REPORT_XX频道发布时，netChecker如何有效且正确的采集到该轮检查
结果并存入Redis中。

应该达到的效果如下：
在下一轮检查结果发布前，netChecker应该已经将本轮检查结果全部订阅完成并存入到Redis中。
原因：如果下一轮检查结果发布后，而该轮的检查结果仍未被全部订阅，那么会造成相应 reids频道中消息堆积，而netChecker每轮
应该只能获取该轮的检查结果，这样一来netChecker每轮获取到的检查结果并不是最新的（因为最新的消息在相应频道的后面排队）。
如果netChecker每轮获取的检查结果是多轮交叉的，会造成结果混乱，无法保证结果覆盖到每个节点都执行3中类型的检测。

解决方案：
为解决上面的问题，采用的处理方案是：
1. netChecker订阅REPORT_XX频道时，一旦订阅到消息，就开始计数，到计数到达k8s集群work节点数量时，表明本轮结果采集
完毕。（因为每种类型的检查，如pod->pod，每个节点都发布自己的检查结果。所以每轮的REPORT_XX频道，最多只有
work节点数量个消息）
2. 有些节点可能在本轮发布检查结果失败，导致方式方案无法满足。故设置一个定时器，当定时一到，表明本轮结果采集完毕。
这个定时的时间设置需要保证：
· 大于所有节点成功发布时，netChecker成功订阅本轮所有消息的时间。
· 保证在下一轮检查结果发布前，netChecker采集本轮消息完毕并能成功存入redis中

通过实践找出定时器的合理时间：
· 每一轮，从netChecker将从k8s获取到的信息发布成功开始，到netChecker成功完毕本轮所有消息，时间差基本为1s。
下面是从日志信息获取的证明：

![](https://github.com/yingyan003/netChecker/blob/master/picture/timer.png)

总结：
每一轮从prepare发布消息成功到receive接收完成该轮的所有数据，相差时间为：
无网络检测失败：基本为1s
网络检测失败：当网络故障时，ping的检测会等待超时时间1s（ping->pong和ping->baidu都会超时1s,一共是2s）
			由于每个节点都部署一个ping,每个节点的ping都会检测包括自己在内的所有节点的网络，如果所有
			网络都不同下，每个ping检测完成的时间最多延迟 2* len(node)s 的时间

# 创建资源对象

### netChecker-deployment.yaml

注意：yaml文件只允许"空格"，不允许"tab键"，否则在k8s集群中创建资源对象时，校验失败。

```
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: netguard
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      name: netChecker
  template:
    metadata:
      labels:
        name: netChecker
    spec:
      containers:
      - name: netguard
        image: img.reg.3g:15000/netChecker:2788
        ports:
        - containerPort: 8080
        env:
        - name: REDIS_HOST
          value: redis-master-svc.zxy:6379
        - name: LOG_LEVEL
          value: "3"
        resources:
          limits:
            cpu: 100m
            memory: 100Mi
          requests:
            cpu: 100m
            memory: 100Mi
```

### netChecker-svc.yaml

```
apiVersion: v1
kind: Service
metadata:
  name: netChecker-svc
spec:
  ports:
  - nodePort: 32089
    port: 8080
    protocol: TCP
  selector:
    name: netChecker
  type: NodePort
```

### ping-daemonset.yaml

```
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: ping
  namespace: zxy
  labels:
    name: ping
    type: ping
spec:
  selector:
    matchLabels:
      name: ping
      type: ping
  template:
    metadata:
      labels:
        name: ping
        type: ping
    spec:
      containers:
      - env:
        - name: REDISHOST
          value: redis-master-svc.zxy:6379
        - name: NODEPORT
          value: "32088"
        - name: LOG_LEVEL
          value: "3"
        name: ping
        image: img.reg.3g:15000/ping:2788
        ports:
        - containerPort: 8080
          protocol: TCP
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
```

### pong-daemonset.yaml

```
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: pong
  namespace: zxy
  labels:
    name: pong
    type: pong
spec:
  selector:
    matchLabels:
      name: pong
      type: pong
  template:
    metadata:
      labels:
        name: pong
        type: pong
    spec:
      containers:
      - env:
        - name: REDIS_HOST
          value: redis-master-svc.zxy:6379
        - name: REDIS_MAX_ACTIVE_CONN
          value: "100"
        name: pong
        image: img.reg.3g:15000/pong:2730
        ports:
        - containerPort: 8080
          protocol: TCP
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
```

### pong-svc.yaml

```
apiVersion: v1
kind: Service
metadata:
  name: netChecker-svc
spec:
  ports:
  - nodePort: 32089
    port: 8080
    protocol: TCP
  selector:
    name: netChecker
  type: NodePort

```

当代码编写测试完成后，推到gitlab编译并打镜像，然后推到镜像仓库。
镜像准备好后，开始编写yaml文件。准备工作完成后就可以开始部署应用了。

备注：kubectl的操作在k8s的master节点上

#### 部署netChecker

部署deployment

```部署deployment
kubectl create -f netChecker-deployment.yaml

```

部署deployment对应的service

```
kubectl create -f netChecker-svc.yaml
```
#### 部署ping

```
kubectl create -f ping-daemonset.yaml
```

### 部署pong

部署daemonSet
```
kubectl create -f pong-daemonset.yaml
```

部署daemonSet对应的service

```
kubectl create -f pong pong-svc.yaml
```
