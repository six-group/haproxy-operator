# HAProxy Operator
HAProxy Operator is a Kubernetes-native solution designed to automate the deployment, configuration, and management of HAProxy instances using Custom Resources to abstract the key components such as backends, frontends, and listens.

## Instalation
### Helm
```console
helm repo add haproxy-operator https://six-group.github.io/haproxy-operator
helm install haproxy-operator six-group/haproxy-operator
```

## Usage

### HAProxy Instance (proxy.haproxy.com/v1alpha1)

An HAProxy instance refers to a single running instance of the HAProxy service. This service can be configured to manage the load balancing and distribution of network traffic among a set of servers or backends within or external to a Kubernetes cluster.

Each HAProxy instance has its own configuration file, named haproxy.cfg and stored as a Secret, which defines all the settings for that instance, including defaults, frontends, and backends. This configuration file specifies how incoming connections are handled, which algorithms are used for load balancing, and how to monitor the health of the backends. Multiple HAProxy instances can be run on the same namespace, each with its own configuration and each listening on different ports.

***Example:***

This is a configuration for an HAProxy instance with two sections: `global` and `defaults`. The `global` section sets process-wide parameters, including the number of threads, maximum concurrent connections, stats socket configuration, buffer sizes, SSL parameters, and logging settings. The `defaults` section sets default parameters for all other sections. It sets the mode to TCP, enables logging, and sets various timeout values for different types of connections and requests.

```
global
  nbthread 4
  maxconn 20000
  stats socket /var/lib/haproxy/run/haproxy.sock expose-fd listeners level admin mode 600
  stats timeout 300000
  tune.bufsize 32768
  tune.maxrewrite 8192
  tune.ssl.default-dh-param 2048
  ssl-default-bind-options ssl-min-ver TLSv1.2 
  ssl-default-bind-ciphers SHA256
  ssl-default-bind-ciphersuites TLS_SHA256
  log /var/lib/rsyslog/rsyslog.sock local0
  log-send-hostname

defaults unnamed_defaults_1
  mode tcp
  log global
  option tcplog
  timeout http-request 10000
  timeout connect 5000
  timeout client 30000
  timeout client-fin 1000
  timeout server 30000
  timeout server-fin 1000
  timeout tunnel 3600000
  timeout http-keep-alive 300000
 ```

```yaml title="haproxy.yaml"
apiVersion: proxy.haproxy.com/v1alpha1
kind: Instance
metadata:
  name: example
  namespace: default
spec:
  configuration:
    defaults:
      logging:
        enabled: true
        tcpLog: true
      mode: tcp
      timeouts:
        client: 30s
        client-fin: 1s
        connect: 5s
        http-keep-alive: 5m0s
        http-request: 10s
        server: 600s
        server-fin: 1s
        tunnel: 1h0m0s
    selector:
      matchLabels:
        proxy.haproxy.com/instance: example
    global:
      logging:
        address: /var/lib/rsyslog/rsyslog.sock
        enabled: true
        facility: local0
      ssl:
        defaultBindCipherSuites:
          - TLS_SHA256
        defaultBindCiphers:
          - SHA256
        defaultBindOptions:
          minVersion: TLSv1.2
      statsTimeout: 5m0s
      tune:
        bufsize: 32768
        maxrewrite: 8192
        ssl:
          defaultDHParam: 2048
      maxconn: 20000
      nbthread: 4
      reload: true
  image: 'haproxy:2.8.0'
  replicas: 2
  network:
    route:
      enabled: false
    service:
      enabled: false
```

[API Reference Instance](docs/api-reference.md#instance) defines all the features that can be configured in an HAProxy instance.

### HAProxy Configuration (config.haproxy.com/v1alpha1)
For the dynamic configuration of HAProxy instances, custom resources have been created for each configuration section, i.e., `listen`, `frontend`, `backend`, and `resolver`.
These configuration resources are associated with particular instances by the use of label selectors. A label selector is specified within the `Instance` configuration, and the corresponding label is applied to each configuration resource to establish a relation.

An example of a label selector used within an `Instance` to match a specific HAProxy instance is provided below:
```yaml
selector:
  matchLabels:
    proxy.haproxy.com/instance: example
```
This approach allows HAProxy instances to be configured dynamically, with a focus on modularity and ease of management.


#### Frontend

`Frontend` defines how incoming connections are handled  based on the rules defined. It specifies the IP addresses and ports that HAProxy listens on and sets rules for what to do with connections once they are received. These rules can include Access Control Lists (ACLs), which allow you to route traffic based on various factors such as the client's IP address, the requested URL, or the type of protocol used. The HAProxy Operator allows you to define frontends in a declarative manner, specifying things like the port number and the default backend.

***Example 1:***

The HAProxy frontend 'example-1' operates in HTTP mode and listens for incoming connections on a Unix socket at /var/lib/haproxy/run/local.sock:9443. It has a certificate file configured, which is used to terminate TLS connections. It also has a default backend configured, which is used when no other rules match an incoming request.


 ```
 frontend example-1
  mode http
  bind unix@/var/lib/haproxy/run/local.sock:9443 name https crt /usr/local/etc/haproxy/ssl-certs.crt ssl accept-proxy crt-list /usr/local/etc/haproxy/cert_list.map
  errorfile 403 /usr/local/etc/haproxy/error-403.http
  use_backend %[base,map_reg(/usr/local/etc/haproxy/edge.map)] if { base,map_reg(/usr/local/etc/haproxy/edge.map) -m found }
  use_backend %[base,map_reg(/usr/local/etc/haproxy/reencrypt.map)] if { base,map_reg(/usr/local/etc/haproxy/reencrypt.map) -m found }
  default_backend default-namespace
 ```

```yaml title="haproxy.yaml"
apiVersion: config.haproxy.com/v1alpha1
kind: Frontend
metadata:
  name: example-1
  namespace: default
  labels:
    proxy.haproxy.com/instance: example
spec:
  backendSwitching:
    - backend:
        regexMapping:
          name: edge
          parameter: base
      condition: '{ base,map_reg(/usr/local/etc/haproxy/edge.map) -m found }'
      conditionType: if
    - backend:
        regexMapping:
          name: reencrypt
          parameter: base
      condition: '{ base,map_reg(/usr/local/etc/haproxy/reencrypt.map) -m found }'
      conditionType: if
  binds:
    - acceptProxy: true
      address: unix@/var/lib/haproxy/run/local.sock
      hidden: true
      name: https
      port: 9443
      ssl:
        certificate:
          name: ssl-certs
          valueFrom:
            - secretKeyRef:
                key: tls.crt
                name: ssl-certs
            - secretKeyRef:
                key: tls.key
                name: ssl-certs
        enabled: true
      sslCertificateList:
        name: cert_list
  defaultBackend:
    name: default-namespace
  errorFiles:
    - code: 403
      file:
        name: error-403
        value: |-
          HTTP/1.0 403 Forbidden
          Pragma: no-cache
          Cache-Control: private, max-age=0, no-cache, no-store
          Connection: close
          Content-Type: text/html

          <!DOCTYPE html>
          <html lang="en">
             <head>
                <title>403 Forbidden</title>
             </head>
          </html>
        valueFrom: {}
  mode: http
```

***Example 2:***

This is a HAProxy frontend configuration named 'example-2'. It operates in TCP mode, binds to a specific IP and port, and inspects TCP requests with a delay. It accepts requests with a specific SSL hello type.

```
frontend example-2
  mode tcp
  bind ${BIND_ADDRESS}:443 name public-ssl
  tcp-request inspect-delay 5000
  tcp-request content accept if { req_ssl_hello_type 1 }
  default_backend default-namespace
```

```yaml title="haproxy.yaml"
apiVersion: config.haproxy.com/v1alpha1
kind: Frontend
metadata:
  name: example-2
  namespace: default
  labels:
    proxy.haproxy.com/instance: example
spec:
  binds:
    - address: '${BIND_ADDRESS}'
      name: public-ssl
      port: 443
  defaultBackend:
    name: default-namespace
  mode: tcp
  tcpRequest:
    - timeout: 5s
      type: inspect-delay
    - action: accept
      condition: '{ req_ssl_hello_type 1 }'
      conditionType: if
      type: content
```

***Example 3:***

This is a HAProxy frontend configuration named 'example-3'. It operates in HTTP mode and binds to a specific IP and port. For every HTTP request, it immediately returns a HTTP 200 OK status with a JSON response indicating a successful health check.


```
frontend example-3
  mode http
  bind ${BIND_ADDRESS}:50055 name health
  http-request return status 200 content-type application/json string "{\"status\":\"OK\"}"
```

```yaml title="haproxy.yaml"
apiVersion: config.haproxy.com/v1alpha1
kind: Frontend
metadata:
  name: example-3
  namespace: default
  labels:
    proxy.haproxy.com/instance: example
spec:
  binds:
    - address: '${BIND_ADDRESS}'
      name: health
      port: 50055
  defaultBackend: {}
  httpRequest:
    return:
      content:
        format: string
        type: application/json
        value: '{\"status\":\"OK\"}'
      status: 200
  mode: http
```
[API Reference Frontend](docs/api-reference.md#frontend) defines all the features that can be configured in an HAProxy frontend.


#### Backend

`Backend` refers to a set of servers that will receive the forwarded requests. The backend section defines how to reach the server, how to check its health, and how to balance the load among the servers. It can contain one or more servers, each server representing an application server in your infrastructure. With the HAProxy Operator, you can define the desired state for your backends in OpenShift, and the operator will ensure that the actual state matches the desired state.

***Example 1:***

This is a HAProxy backend configuration named 'example-1'. It operates in TCP mode and defines an Access Control List (ACL) for a specific source IP. It enables connection redispatching with a maximum of 3 retries per request. It rejects TCP requests not matching the ACL. It defines a server with specific health check settings, initial address resolution disabled, a specific check interval, and specified resolvers for hostname resolution.

```
backend example-1
  mode tcp
  acl whitelist src 0.0.0.0
  option redispatch 3
  tcp-request content reject if !whitelist
  server web web.namespace.svc.cluster.local:443 check init-addr none inter 500 resolvers dns-namespace
```

```yaml title="haproxy.yaml"
apiVersion: config.haproxy.com/v1alpha1
kind: Backend
metadata:
  name: example-1
  namespace: default
  labels:
    proxy.haproxy.com/instance: example
spec:
  acl:
    - criterion: src
      name: whitelist
      values:
        - 0.0.0.0
  mode: tcp
  redispatch: true
  servers:
    - address: web.namespace.svc.cluster.local
      check:
        enabled: true
        inter: 500ms
      initAddr: none
      name: web
      port: 443
      resolvers:
        name: dns-namespace
  tcpRequest:
    - action: reject
      condition: '!whitelist'
      conditionType: if
      type: content
```

***Example 2:***

The HAProxy backend 'example-2' operates in HTTP mode. It has an Access Control List (ACL) named "whitelist" that matches when the source IP of the request is 0.0.0.0. It adds the X-Forwarded-For header to preserve the client's IP address, redistributes sessions in case of failure, and sets a health check timeout of 5 seconds. If a TCP request doesn't match the "whitelist" ACL, it's rejected. Various X-Forwarded-* and Forwarded headers are added to the HTTP request to convey information about the original request.  A server named "web" is defined within this backend, with health checks enabled and an interval of 500 milliseconds between checks. The server's hostname resolution uses the "dns-namespace" resolvers. SSL/TLS configuration and certificate verification are also specified for this server and its weight is set to 256.

```
backend example-2
  mode http
  acl whitelist src 0.0.0.0
  option forwardfor
  option redispatch 3
  timeout check 5000
  tcp-request content reject if !whitelist
  http-request add-header X-Forwarded-Host %[req.hdr(host)]
  http-request add-header X-Forwarded-Port %[dst_port]
  http-request add-header X-Forwarded-Proto http if !{ ssl_fc }
  http-request add-header X-Forwarded-Proto-Version h2 if { ssl_fc_alpn -i h2 }
  http-request add-header Forwarded for=%[src];host=%[req.hdr(host)];proto=%[req.hdr(X-Forwarded-Proto)]
  cookie e76a2f0f39106e5e833f1323866171d4 attr SameSite=None httponly indirect nocache insert secure
  server web web.namespace.svc.cluster.local:443 check ssl alpn http/1.1,h2 ca-file /usr/local/etc/haproxy/service-ca.crt cookie 4b24b04d486a91808d248592b93d2293 init-addr none inter 500 resolvers dns-namespace verify required verifyhost web.namespace.svc weight 256
```

```yaml title="haproxy.yaml"
apiVersion: config.haproxy.com/v1alpha1
kind: Backend
metadata:
  name: example-2
  namespace: default
  labels:
    proxy.haproxy.com/instance: example
spec:
  mode: http
  cookie:
    attribute:
      - SameSite=None
    httpOnly: true
    indirect: true
    mode:
      insert: true
      prefix: false
      rewrite: false
    name: app
    noCache: true
    secure: true
  forwardFor:
    enabled: true
  httpRequest:
    addHeader:
      - name: X-Forwarded-Host
        value:
          str: '%[req.hdr(host)]'
      - name: X-Forwarded-Port
        value:
          str: '%[dst_port]'
      - condition: '!{ ssl_fc }'
        conditionType: if
        name: X-Forwarded-Proto
        value:
          str: http
      - condition: '{ ssl_fc_alpn -i h2 }'
        conditionType: if
        name: X-Forwarded-Proto-Version
        value:
          str: h2
      - name: Forwarded
        value:
          str: 'for=%[src];host=%[req.hdr(host)];proto=%[req.hdr(X-Forwarded-Proto)]'
  acl:
    - criterion: src
      name: whitelist
      values:
        - 0.0.0.0
  redispatch: true
  tcpRequest:
    - action: reject
      condition: '!whitelist'
      conditionType: if
      type: content
  servers:
    - port: 443
      initAddr: none
      verifyHost: web.namespace.svc
      cookie: true
      check:
        enabled: true
        inter: 500ms
      name: web
      ssl:
        alpn:
          - http/1.1
          - h2
        caCertificate:
          name: service-ca.crt
          valueFrom:
            - configMapKeyRef:
                key: service-ca.crt
                name: openshift-service-ca.crt
        enabled: true
        verify: required
      resolvers:
        name: dns-namespace
      address: web.namespace.svc.cluster.local
      weight: 256
  timeouts:
    check: 5s
```

[API Reference Backend](docs/api-reference.md#backend) defines all the features that can be configured in an HAProxy backend.
