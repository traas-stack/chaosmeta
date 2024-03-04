![](docs/static/logo.png)
[中文版README](README_CN.md)

[Official Document](https://chaosmeta.gitbook.io/chaosmeta-en)
# Introduction

ChaosMeta is a cloud-native chaos engineering platform open sourced by Ant Group. It embodies the methodologies, technologies and products that Ant Group has accumulated over many years in the practice of large-scale red and blue offensive and defensive drills at the company level. With the "Risk Catalog" (internal general risk scenario manual for technical components in various fields) as theoretical guidance, combined with technical practice, it has escorted Ant Group's various promotional activities for many years.

ChaosMeta is a platform dedicated to supporting all stages of fault drills, covering platform capabilities in multiple stages such as access detection, traffic injection, fault injection, fault measurement, fault recovery, and recovery measurement. While liberating productivity for users, it is also pursuing the future form of chaos engineering: one-click automated drills, and even intelligent drills.

# Core advantages
#### Simple and easy to use, provides user interface, low threshold for use
Support visual user interface, Kubernetes API, command line, HTTP API, and other methods.
[![docs/static/componentlink.png
](docs/static/workflow.png)](https://player.bilibili.com/player.html?aid=276433781&bvid=BV1yF411m7b4&cid=1280401525&p=1)

#### Fully verified by a large amount of practical experience, high reliability

The Blue Army team of Ant Group has been deeply involved in the chaos engineering industry for many years. It holds company-level large-scale red and blue offensive and defensive drills every year, facing all the company's businesses, and many businesses also conduct 7X24-hour drills and monthly normal drills

Internal drill object types cover cloud products, Kubernetes, Operator applications, databases (OceanBase, Etcd, etc.), middleware (message queues, distributed scheduling, configuration centers, etc.), business applications (Java applications, C++ applications, Golang applications)

#### High flexibility, supporting a variety of user needs

Whether the user wants a complete chaos engineering platform, or just wants the underlying platform capabilities such as remote injection, orchestration and scheduling, or even just wants the single-machine fault injection capability, or manages and injects targets on or off the cloud Failure, there are corresponding deployment plans to meet
![](docs/static/componentlink.png)

#### Rich fault injection capabilities, cloud native chaos engineering

Because Ant Group attaches great importance to offensive and defensive drills, it has led to large-scale and high-frequency drills, which in turn has promoted the construction of various fault injection capabilities. And because Ant has a huge internal infrastructure scale, coupled with the low fault tolerance of finance, the stability requirements for infrastructure such as Kubernetes and middleware are very high. Therefore, Ant Chaos Engineering has accumulated rich fault capabilities in the cloud-native field. and exercise experience.


#### The platform has powerful capabilities, supports the complete "chaos engineering life cycle", and is oriented towards automation.
ChaosMeta covers access detection, traffic injection, fault injection, fault measurement, fault recovery, recovery measurement and other stages of platform capabilities, as the technical basis of "automated chaos engineering"
![](docs/static/lifecycle_en.png)

In addition to the platform capability support of the exercise process, another big mountain in the automated exercise is the design of the experiment. At present, it is difficult to completely rely on machines to automatically design. However, we can systematically abstract the reusable experience and organize it into a book. When conducting chaos engineering exercises on the same type of components, we can quickly reuse it. This is the original intention of the risk catalog design

<img src="docs/static/riskdir_en.png" width="50%" >

ChaosMeta will realize the automated drill capability of **one-click physical examination** based on the technical foundation of "Chaos Engineering Life Cycle" and the theoretical basis of "Risk Catalog", directly generate the target stability score, and greatly liberate users in chaos

# Architecture overview
![](docs/static/architecture.png)
##### User layer (Client)
The Client layer is mainly composed of **chaosmeta-platform** components. Its main task is to lower the threshold for users to use and provide a visual interface to facilitate users' planning, orchestration, experiment configuration, experiment record details, and Agent management (pods/node of k8s clusters, cross-cluster objects, non-k8s physical machines/containers, etc.) and other platform capabilities.

##### Engine layer (Engine)
The Engine layer includes the core platform capabilities of ChaosMeta and the implementation of some cloud-native fault capabilities, including the following components:
- **chaosmeta-CRD**: ChaosMeta's platform capabilities are developed based on the Operator framework, so each type of capability has a corresponding CRD, and then the corresponding Operator monitors the status and performs the corresponding operations. For example, the CRD of the fault injection capability is experiments.inject.chaosmeta.io and the corresponding monitoring operator is chaosmeta-inject-operator. Therefore, users can create corresponding CR instances through Kubectl or Kubernetes-Client to perform corresponding capabilities;


- **chaosmeta-inject-operator**: Listens to CR instances related to fault injection created by users, compares the actual status of CR in the cluster with the expected status in the control loop to execute relevant fault injection logic and status transfer, and converts the actual status Tune into the desired state. Different operations are performed based on the fault type defined by the CR instance. For example: if it is a system resource fault, remote injection is required through chaosmeta-daemonset or HTTP or command channel; if it is a cloud native fault, injection will be based on Kubernetes APIServer. , and if it involves a dynamic admission failure, chaosmeta-webhook will also be requested to update the tampering rules and interception rules;


- **chaosmeta-webhook**: The API processing process of each APIServer needs to go through authentication, authentication, and admission, and the admission stage will go through the Mutating Admission Webhook (tampering) and Validating Admission Webhook (verification) stages, chaosmeta -webhook will update the resource matching rules according to the fault definition, and intercept, tamper with, delay, and exception the user's Kubernetes resource creation request. This is very meaningful for failure drill scenarios related to Operator applications and Kubernetes' own cluster robustness.


- **chaosmeta-measure-operator**: This is the component used to perform measurement capabilities, mainly used in two phases: failure measurement and recovery measurement. The fault metric is an effectiveness measure of the fault injection effect, while the recovery metric is an effectiveness measure of the resilience of the defense platform. Measurement capabilities are the key capabilities to achieve automation and intelligence in chaos engineering.


*For example, the failure effect of a drill is expected to be that the number of successful requests for a certain service drops by 50%, and the corresponding defense platform is expected to be able to detect it within 5 minutes and recover within 10 minutes. The execution method is to achieve full CPU usage. Then the fault measurement phase must find the time point when the number of successful service requests drops by 50% compared to before the fault injection (fault effective point). In the recovery measurement phase, it is necessary to find the time point when the corresponding alarm is generated (fault discovery point), and also to find the time point after the fault discovery point to request a successful amount to restore the water level before the drill (fault recovery point). Finally, an analysis report of the exercise was generated, giving areas for improvement in the defense platform.*

- **chaosmeta-workflow-operator**: Provides fault orchestration capabilities. Because in reality, except for a single failure scenario. There are also demands for a large number of complex fault scenarios, which require simulation through serial and parallel combinations of different fault injection capabilities. And orchestration is not limited to fault injection, but can also include orchestration nodes with different capability types such as traffic injection, fault admission detection, fault measurement, recovery measurement, etc. This is also a key capability for automating drills.

- **chaosmeta-flow-operator**: This is a component used to perform traffic injection, mainly used to mock the traffic of the target services. Because when we conduct fault drills, we often need to meet the flow rate to achieve the effect of the fault. For example, if you want to trigger a service delay alarm for a certain service, it is not enough to inject the delay into the container network of this service. If there is no traffic request, the corresponding monitoring alarm will not be triggered.

##### Kernel layer (Kernel)
The Kernel layer mainly includes the implementation of single-machine fault injection capability, mainly including the **chaosmetad** component, which provides the method of resident HTTP service and command line execution, and also encapsulates the corresponding daemonset component (**chaosmeta-daemonset**). The training platform can be flexibly matched with different needs.

# Capabilities of the current version
The current version has released: user interface, fault injection scheduling engine, measurement engine, traffic injection engine, single machine fault injection tool and other components
#### User Interface
- Provides experiment orchestration capabilities and lowers the threshold for use (the current version of the UI does not yet support traffic injection type and measurement type nodes);
- Provides the ability to inject and filter remote targets of Pod/Node in the cluster (the UI will support targets outside the cluster in the future);
- Provides space management capabilities and can separate and manage data on demand;
- Provide account permission management system.
#### Fault injection capability
- System Resources Exception: CPU, memory, network, disk, process, file, etc.;
- Kernel Resource Exception: fd, nproc, etc.;
- JVM Dynamic Injection: function call delay, function return value tampering, function throwing exception, etc.;
- Container Fault Injection: kill container, suspend container, CPU, memory, network, disk, process, file, JVM injection and other experimental scenarios in the container;
- Kubernetes Injection: execute experimental scenarios such as CPU, memory, network, disk, process, file, JVM injection on any pod;
- Cloud-Native Faults: Abnormalities in cluster resources such as accumulation of a large number of Pending Pods and Completed Jobs; there are also abnormalities in instances of cloud-native resources such as Deployment, Node, and Pod, such as copy expansion and shrinking tampering of Deployment instances, and injection of Pod instance Finalizers.
#### Measuring Capabilities
- monitor: Make expected judgments on the values of monitoring items, such as whether the CPU usage monitoring value of a certain machine is greater than 90%. Prometheus is supported by default.
- pod: Make expected judgments on pod-related data, such as whether the number of pod instances of an application is greater than 3
- http: Make expected judgments on http requests. For example, when making a specified http request, whether the return status code is 200
- tcp: Make expected judgments on tcp requests, such as testing whether the 8080 port of a certain server is connectable
#### Traffic injection capability
- http: http traffic injection

# Getting Start
#### Quickly try the single-machine injection capability
```shell
# Download docker mirror and run container
docker run --privileged -it registry.cn-hangzhou.aliyuncs.com/chaosmeta/chaosmetad-demo:v0.5.3 /bin/bash

# Start the test service
cd /tmp && python -m SimpleHTTPServer 8080 > server.log 2>&1 &
curl 127.0.0.1:8080

# Create an experiment to inject a 2s network delay into the lo network card, and it will automatically recover after 10 minutes
chaosmetad inject network delay -i lo -l 2s --uid test-fg3g4 -t 10m

# View experiment information, test effect
chaosmetad query
curl 127.0.0.1:8080

# Manually recover the experiment
chaosmetad recover test-fg3g4
```
#### Fault Ability Usage
For details, see: [Function Instructions](https://chaosmeta.gitbook.io/chaosmeta-en/capability-instruction)
#### Installation Guide
For details, see: [Installation Guide](https://chaosmeta.gitbook.io/chaosmeta-en/installation-guide)
# Communicate
Welcome to submit defects, questions, suggestions and new features, all problems can be submitted to [Github Issues](https://github.com/traas-stack/chaosmeta/issues/new), you can also contact us in the following ways:
- DingTalk group: 21765030887
- Slack group: [ChaosMeta](https://app.slack.com/client/T057ERYMS8J/C057883SM38?geocode=zh-cn)
- WeChat public account: ChaosMeta混沌工程
- Twitter：AntChaosMeta
- Email: chaosmeta.io@gmail.com
- WeChat group: email communication/WeChat public account to obtain QR code invitation

# RoadMap
### Platform capabilities
The future evolution of ChaosMeta platform capabilities is divided into three stages
##### Phase 1 - Manual Configuration
The goal to be achieved is to open all the components in the architecture diagram to the outside world. At this time, it can support the complete life cycle of chaos engineering, enter the field of primary automated chaos engineering, and use the "risk catalog" as a theoretical reference. Once manual configuration, multiple times automatically.
The order of opening to the outside world is as follows (if you have relevant needs, you are welcome to submit an issue, and priority adjustments will be considered):
- [x] Stand-alone fault injection tool：chaosmetad
- [x] Fault Remote Injection Engine：chaosmeta-inject-operator
- [x] Platform Dashboard：chaosmeta-platform
- [x] Orchestration Engine：chaosmeta-workflow-operator
- [x] Measure Engine：chaosmeta-measure-operator
- [x] Traffic Injection Engine：chaosmeta-flow-operator
- [ ] Risk Catalog：Common Risk Scenario Handbook for Technical Components in Each Field
- [ ] Cloud Native Dynamic Access Fault Injection Capability：chaosmeta-webhook

##### Phase 2 - Automation
At this stage, the "Risk Catalog" will play a greater role. It not only gives the risk of a class of applications, but also the corresponding prevention and emergency recommendations, and the score of each item, and ChaosMeta will The "risk catalog" is integrated into a risk medical examination package of general components, which realizes the one-click "physical examination" capability, inputs target application information, and directly outputs a risk score and risk analysis report.
##### Phase 3 - intelligence
Explore the direction of combining artificial intelligence
### Fault Injection Capability
The following is just a classification of fault capabilities. For the specific atomic fault capabilities provided, please refer to the [description of fault capabilities](https://chaosmeta.gitbook.io/chaosmeta-en/capability-instruction) (welcome to submit issues and put forward new capability requirements, and those with higher requirements are given priority):
![](docs/static/roadmap.png)
# License
ChaosMeta follows the Apache 2.0 license, please read [LICENSE](LICENSE) for details
