A simple infrastructure debug orchestrator to bring your thoughts into discrete steps and tying them all together to effectively find or resolve issues.

# Motivation
We all have nuances in how we debug issues and somehow we tend to think its an art. The aim with cortex is to bring some science and automation into how we debug infrastructural problems.The main task cortex tries to solve is bring some structure to the art and provide an easier way to mimic your thought process so that its easier to share with others
and not having to ponder "What did I do 2 weeks ago". The hope with this tool is that it would help the SRE to think in discrete steps that could then be collated, reused and executed in different ways which would help not just in expressing the tasks better, but also become primers for the juniors in the team to learn the different ways to debug 

# How does cortex work?
Cortex works on the following principles:
1. Every action we take can be a discretely identified task through and is called a neuron
2. Every action has an output that could end the debugging session or be passed to the next neuron, or be discarded
3. You could create a plan composing of multiple neurons to be run in parallel or sequence and is called a synapse
4. Each neuron could be fired independently by multiple synapse and needs to handle state accordingly
5. Synapse should make the determination of resolution based on analysis of outputs from the neurons

The architecture is quite simple, and think of synapse building directed acyclic graph of parallel vertices and sequential plans which are then executed in an event loop.

<img src="./assets/cortex.jpg" alt="Synapse"
	title="Architecture" width="500" height="500" />

## Creating neurons
Neurons are folders with the following naming conventions:
1. If it does not mutate anything, starts with "check_" as the suffix. For eg: "check_web_proxy_connection_config".
2. It its a mutating neuron i.e it updates a config or property, use "mutate_" as the suffix. For eg: "mutate_web_proxy_connection_config"

The name of the folder should be unique and give enough indication on the activity

To create a neuron, run 
`cortex create-neuron check_web_proxy_conn_config`
 for more options, run
`cortex create-neuron -h`

It would create a folder and bootstrap files as below:

```
check_web_proxy_conn_config
    |----- neuron.yaml
    |----- run.sh
    |----- run.ps1
```

sample neuron.yaml:

```
---
name: check_web_proxy_conn_config
description: "A longer description"
exec_file: run.sh
pre_exec_debug: "Going to check the web_proxy connection configuration"
assertExitStatus: [0, 137]
post_exec_success_debug: "All configurations checkout ok"
post_exec_fail_debug:
 120: "Found maxconn rate to be too low"
 110: "Found maxpipes to be too low"

```


As you notice, the exit code has a lot of importance in how your neurons communicate the a

## Creating synapse   

A sample synapse yaml when you want to fix something when you find and error:

```
---
name: app_network_latency
plan:
 - definition:
    - neuron: check_web_proxy_conn_config
      config:
        path: /usr/neurons/check_web_proxy_conn_config
        fix:
          - 120: mutate_web_proxy_conn_bump_maxconn_config
          - 110: mutate_web_proxy_conn_bump_maxpipes_config
    - neuron: check_api_gateway_conn_config
      config:
        path: /usr/neurons/check_web_proxy_conn_config
        fix:
          - 120: mutate_api_gateway_conn_bump_maxconn_config
          - 110: mutate_api_gateway_conn_bump_maxpipes_config
    - neuron: mutate_web_proxy_conn_bump_maxconn_config
      config:
        path: /usr/neurons/mutate_web_proxy_conn_bump_maxconn_config
     
 - serial
    - check_api_gateway_conn_config
    - check_web_proxy_cpu_usage
    - check_grafana_cpu_trend
    
```

 A synapse that only checks and does not mutate:               

```
---
name: app_network_latency
 - definition:
     - neuron: check_web_proxy_conn_config
      config:
        type: check
        path: /usr/neurons/check_web_proxy_conn_config
    - neuron: check_api_gateway_conn_config
      config:
        type: check
        path: /usr/neurons/check_web_proxy_conn_config
plan:
  - config:
     - exit_on_first_error: false
  - steps:
      serial
        - check_web_proxy_conn_config
        - check_api_gateway_conn_config 
```

 A synapse that that runs the checks in parallel:               
    
```
---
name: app_network_latency
 - definition:
     - neuron: check_web_proxy_conn_config
      config:
        path: /usr/neurons/check_web_proxy_conn_config
        fix:
          - 120: mutate_web_proxy_conn_bump_maxconn_config
          - 110: mutate_web_proxy_conn_bump_maxpipes_config
    - neuron: check_api_gateway_conn_config
      config:
        path: /usr/neurons/check_web_proxy_conn_config
        fix:
          - 120: mutate_api_gateway_conn_bump_maxconn_config
          - 110: mutate_api_gateway_conn_bump_maxpipes_config
    - neuron: mutate_web_proxy_conn_bump_maxconn_config
      config:
 plan:
   - config:
     - exit_on_first_error: false
   - parallel
     - check_web_proxy_conn_config
     - check_api_gateway_conn_config
     - check_web_proxy_cpu_usage
     - check_grafana_cpu_trend
```


