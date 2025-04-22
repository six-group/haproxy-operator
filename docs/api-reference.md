# API Reference

## Packages
- [config.haproxy.com/v1alpha1](#confighaproxycomv1alpha1)
- [proxy.haproxy.com/v1alpha1](#proxyhaproxycomv1alpha1)


## config.haproxy.com/v1alpha1

Package v1alpha1 contains API Schema definitions for the config v1alpha1 API group

### Resource Types
- [Backend](#backend)
- [Frontend](#frontend)
- [Listen](#listen)
- [Resolver](#resolver)



#### ACL







_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name |  | Pattern: `^[^\s]+$` <br /> |
| `criterion` _string_ | Criterion is the name of a sample fetch method, or one of its ACL<br />specific declinations. |  | Pattern: `^[^\s]+$` <br /> |
| `values` _string array_ | Values are of the type supported by the criterion. |  |  |


#### Backend



Backend is the Schema for the backend API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1` | | |
| `kind` _string_ | `Backend` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[BackendSpec](#backendspec)_ |  |  |  |


#### BackendReference







_Appears in:_
- [BackendSwitchingRule](#backendswitchingrule)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name of a specific backend |  |  |
| `regexMapping` _[RegexBackendMapping](#regexbackendmapping)_ | Mapping of multiple backends |  |  |


#### BackendSpec



BackendSpec defines the desired state of Backend



_Appears in:_
- [Backend](#backend)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. | http | Enum: [http tcp] <br /> |
| `httpResponse` _[HTTPResponseRules](#httpresponserules)_ | HTTPResponse rules define a set of rules which apply to layer 7 processing. |  |  |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |  |  |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |  |  |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform<br />content switching and generally to take decisions based on content extracted<br />from the request, the response or any environmental status |  |  |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel.<br />The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit.<br />More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |  |  |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |  |  |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |  |  |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |  |  |
| `httpLog` _boolean_ | HTTPLog enables HTTP log format which is the most complete and the best suited for HTTP proxies. It provides<br />the same level of information as the TCP format with additional features which<br />are specific to the HTTP protocol. |  |  |
| `tcpLog` _boolean_ | TCPLog enables advanced logging of TCP connections with session state and timers. By default, the log output format<br />is very poor, as it only contains the source and destination addresses, and the instance name. |  |  |
| `checkTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | CheckTimeout sets an additional check timeout, but only after a connection has been already<br />established. |  |  |
| `servers` _[Server](#server) array_ | Servers defines the backend servers and its configuration. |  |  |
| `serverTemplates` _[ServerTemplate](#servertemplate) array_ | ServerTemplates defines the backend server templates and its configuration. |  |  |
| `balance` _[Balance](#balance)_ | Balance defines the load balancing algorithm to be used in a backend. |  |  |
| `hostRegex` _string_ | HostRegex specifies a regular expression used for backend switching rules. |  |  |
| `hostCertificate` _[CertificateListElement](#certificatelistelement)_ | HostCertificate specifies a certificate for that host used in the crt-list of a frontend |  |  |
| `redispatch` _boolean_ | Redispatch enable or disable session redistribution in case of connection failure |  |  |
| `hashType` _[HashType](#hashtype)_ | HashType specifies a method to use for mapping hashes to servers |  |  |
| `cookie` _[Cookie](#cookie)_ | Cookie enables cookie-based persistence in a backend. |  |  |
| `httpchk` _[HTTPChk](#httpchk)_ | HTTPChk Enables HTTP protocol to check on the servers health |  |  |
| `tcpCheck` _boolean_ | TCPCheck Perform health checks using tcp-check send/expect sequences |  |  |


#### BackendSwitchingRule







_Appears in:_
- [FrontendSpec](#frontendspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `backend` _[BackendReference](#backendreference)_ | Backend reference used to resolve the backend name. |  |  |


#### Balance







_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `algorithm` _string_ | Algorithm is the algorithm used to select a server when doing load balancing. This only applies when no persistence information is available, or when a connection is redispatched to another server. |  | Enum: [roundrobin static-rr leastconn first source uri hdr random rdp-cookie] <br /> |


#### BaseSpec







_Appears in:_
- [BackendSpec](#backendspec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. | http | Enum: [http tcp] <br /> |
| `httpResponse` _[HTTPResponseRules](#httpresponserules)_ | HTTPResponse rules define a set of rules which apply to layer 7 processing. |  |  |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |  |  |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |  |  |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform<br />content switching and generally to take decisions based on content extracted<br />from the request, the response or any environmental status |  |  |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel.<br />The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit.<br />More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |  |  |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |  |  |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |  |  |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |  |  |
| `httpLog` _boolean_ | HTTPLog enables HTTP log format which is the most complete and the best suited for HTTP proxies. It provides<br />the same level of information as the TCP format with additional features which<br />are specific to the HTTP protocol. |  |  |
| `tcpLog` _boolean_ | TCPLog enables advanced logging of TCP connections with session state and timers. By default, the log output format<br />is very poor, as it only contains the source and destination addresses, and the instance name. |  |  |


#### Bind







_Appears in:_
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name for these sockets, which will be reported on the stats page. |  |  |
| `address` _string_ | Address can be a host name, an IPv4 address, an IPv6 address, or '*' (is equal to the special address "0.0.0.0"). |  | Pattern: `^[^\s]+$` <br /> |
| `port` _integer_ | Port |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `portRangeEnd` _integer_ | PortRangeEnd if set it must be greater than Port |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `transparent` _boolean_ | Transparent is an optional keyword which is supported only on certain Linux kernels. It<br />indicates that the addresses will be bound even if they do not belong to the<br />local machine, and that packets targeting any of these addresses will be<br />intercepted just as if the addresses were locally configured. This normally<br />requires that IP forwarding is enabled. Caution! do not use this with the<br />default address '*', as it would redirect any traffic for the specified port. |  |  |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |  |  |
| `hidden` _boolean_ | Hidden hides the bind and prevent exposing the Bind in services or routes |  |  |
| `acceptProxy` _boolean_ | AcceptProxy enforces the use of the PROXY protocol over any connection accepted by any of<br />the sockets declared on the same line. |  |  |


#### CertificateListElement







_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `certificate` _[SSLCertificate](#sslcertificate)_ | Certificate that will be presented to clients who provide a valid<br />TLSServerNameIndication field matching the SNIFilter. |  |  |
| `sniFilter` _string_ | SNIFilter specifies the filter for the SSL Certificate.  Wildcards are supported in the SNIFilter. Negative filter are also supported. |  |  |
| `alpn` _string array_ | Alpn enables the TLS ALPN extension and advertises the specified protocol<br />list as supported on top of ALPN. |  |  |
| `ocsp` _boolean_ | Ocsp Enable OCSP stapling for a specific certificate |  |  |
| `ocsp_file` _[OcspFile](#ocspfile)_ | OcspFile you can save the OCSP response to a file so that HAProxy loads it during startup. |  |  |


#### Check







_Appears in:_
- [Server](#server)
- [ServerParams](#serverparams)
- [ServerTemplate](#servertemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enable enables health checks on a server. If not set, no health checking is performed, and the server is always<br />considered available. |  |  |
| `inter` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Inter sets the interval between two consecutive health checks. If left unspecified, the delay defaults to 2000 ms. |  |  |
| `rise` _integer_ | Rise specifies the number of consecutive successful health checks after a server will be considered as operational.<br />This value defaults to 2 if unspecified. |  |  |
| `fall` _integer_ | Fall specifies the number of consecutive unsuccessful health checks after a server will be considered as dead.<br />This value defaults to 3 if unspecified. |  |  |


#### Cookie







_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name of the cookie which will be monitored, modified or inserted in order to bring persistence. |  |  |
| `mode` _[CookieMode](#cookiemode)_ | Mode could be 'rewrite', 'insert', 'prefix'. Select one. |  |  |
| `indirect` _boolean_ | Indirect no cookie will be emitted to a client which already has a valid one<br />for the server which has processed the request. |  |  |
| `noCache` _boolean_ | NoCache recommended in conjunction with the insert mode when there is a cache<br />between the client and HAProx |  |  |
| `postOnly` _boolean_ | PostOnly ensures that cookie insertion will only be performed on responses to POST requests. |  |  |
| `preserve` _boolean_ | Preserve only be used with "insert" and/or "indirect". It allows the server<br />to emit the persistence cookie itself. |  |  |
| `httpOnly` _boolean_ | HTTPOnly add an "HttpOnly" cookie attribute when a cookie is inserted.<br />It doesn't share the cookie with non-HTTP components. |  |  |
| `secure` _boolean_ | Secure add a "Secure" cookie attribute when a cookie is inserted. The user agent<br />never emits this cookie over non-secure channels. The cookie will be presented<br />only over SSL/TLS connections. |  |  |
| `dynamic` _boolean_ | Dynamic activates dynamic cookies, when used, a session cookie is dynamically created for each server,<br />based on the IP and port of the server, and a secret key. |  |  |
| `domain` _string array_ | Domain specify the domain at which a cookie is inserted. You can specify<br />several domain names by invoking this option multiple times. |  |  |
| `maxIdle` _integer_ | MaxIdle cookies are ignored after some idle time. |  |  |
| `maxLife` _integer_ | MaxLife cookies are ignored after some life time. |  |  |
| `attribute` _string array_ | Attribute add an extra attribute when a cookie is inserted. |  |  |


#### CookieMode







_Appears in:_
- [Cookie](#cookie)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `rewrite` _boolean_ | Rewrite the cookie will be provided by the server. |  |  |
| `insert` _boolean_ | Insert cookie will have to be inserted by haproxy in server responses. |  |  |
| `prefix` _boolean_ | Prefix is needed in some specific environments where the client does not support<br />more than one single cookie and the application already needs it. |  |  |


#### Deny







_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `enabled` _boolean_ | Enabled enables deny http request |  |  |
| `denyStatus` _integer_ | DenyStatus is the HTTP status code. |  | Maximum: 599 <br />Minimum: 200 <br /> |


#### ErrorFile







_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [DefaultsConfiguration](#defaultsconfiguration)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `code` _integer_ | Code is the HTTP status code. |  | Enum: [200 400 401 403 404 405 407 408 410 413 425 429 500 501 502 503 504] <br /> |
| `file` _[StaticHTTPFile](#statichttpfile)_ | File designates a file containing the full HTTP response. |  |  |


#### ErrorFileValueFrom







_Appears in:_
- [StaticHTTPFile](#statichttpfile)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `configMapKeyRef` _[ConfigMapKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#configmapkeyselector-v1-core)_ | ConfigMapKeyRef selects a key of a ConfigMap. |  |  |


#### Forwardfor







_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ |  |  |  |
| `except` _string_ | Pattern: ^[^\s]+$ |  |  |
| `header` _string_ | Pattern: ^[^\s]+$ |  |  |
| `ifnone` _boolean_ |  |  |  |


#### Frontend



Frontend is the Schema for the frontends API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1` | | |
| `kind` _string_ | `Frontend` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[FrontendSpec](#frontendspec)_ |  |  |  |


#### FrontendSpec



FrontendSpec defines the desired state of Frontend



_Appears in:_
- [Frontend](#frontend)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. | http | Enum: [http tcp] <br /> |
| `httpResponse` _[HTTPResponseRules](#httpresponserules)_ | HTTPResponse rules define a set of rules which apply to layer 7 processing. |  |  |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |  |  |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |  |  |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform<br />content switching and generally to take decisions based on content extracted<br />from the request, the response or any environmental status |  |  |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel.<br />The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit.<br />More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |  |  |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |  |  |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |  |  |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |  |  |
| `httpLog` _boolean_ | HTTPLog enables HTTP log format which is the most complete and the best suited for HTTP proxies. It provides<br />the same level of information as the TCP format with additional features which<br />are specific to the HTTP protocol. |  |  |
| `tcpLog` _boolean_ | TCPLog enables advanced logging of TCP connections with session state and timers. By default, the log output format<br />is very poor, as it only contains the source and destination addresses, and the instance name. |  |  |
| `binds` _[Bind](#bind) array_ | Binds defines the frontend listening addresses, ports and its configuration. |  | MinItems: 1 <br /> |
| `backendSwitching` _[BackendSwitchingRule](#backendswitchingrule) array_ | BackendSwitching rules specify the specific backend used if/unless an ACL-based condition is matched. |  |  |
| `defaultBackend` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | DefaultBackend to use when no 'use_backend' rule has been matched. |  |  |


#### HTTPChk







_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `uri` _string_ | URI |  |  |
| `method` _string_ | Method http method<br />Enum: [HEAD PUT POST GET TRACE PATCH DELETE CONNECT OPTIONS] |  | Enum: [HEAD PUT POST GET TRACE PATCH DELETE CONNECT OPTIONS] <br /> |


#### HTTPDeleteHeaderRule







_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `name` _string_ | Name specifies the header name |  |  |
| `method` _string_ | Method is the matching applied on the header name |  | Enum: [str beg end sub reg] <br /> |


#### HTTPHeaderRule







_Appears in:_
- [HTTPRequestRules](#httprequestrules)
- [HTTPResponseRules](#httpresponserules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `name` _string_ | Name specifies the header name |  |  |
| `value` _[HTTPHeaderValue](#httpheadervalue)_ | Value specifies the header value |  |  |


#### HTTPHeaderValue







_Appears in:_
- [HTTPHeaderRule](#httpheaderrule)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `env` _[EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#envvar-v1-core)_ | Env variable with the header value |  |  |
| `str` _string_ | Str with the header value |  |  |
| `format` _string_ | Format specifies the format of the header value (implicit default is '%s') |  |  |


#### HTTPPathRule







_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `format` _string_ | Value specifies the path value |  |  |




#### HTTPRequestRules







_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `setHeader` _[HTTPHeaderRule](#httpheaderrule) array_ | SetHeader sets HTTP header fields |  |  |
| `setPath` _[HTTPPathRule](#httppathrule) array_ | SetPath sets request path |  |  |
| `addHeader` _[HTTPHeaderRule](#httpheaderrule) array_ | AddHeader appends HTTP header fields |  |  |
| `delHeader` _[HTTPDeleteHeaderRule](#httpdeleteheaderrule) array_ | DelHeader removes all HTTP header fields |  |  |
| `redirect` _[Redirect](#redirect) array_ | Redirect performs an HTTP redirection based on a redirect rule. |  |  |
| `replacePath` _[ReplacePath](#replacepath) array_ | ReplacePath matches the value of the path using a regex and completely replaces it with the specified format.<br />The replacement does not modify the scheme, the authority and the query-string. |  |  |
| `deny` _[Deny](#deny) array_ | Deny stops the evaluation of the rules and immediately rejects the request and emits an HTTP 403 error.<br />Optionally the status code specified as an argument to deny_status. |  |  |
| `return` _[HTTPReturn](#httpreturn)_ | Return stops the evaluation of the rules and immediately returns a response. |  |  |


#### HTTPResponseRules







_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `setHeader` _[HTTPHeaderRule](#httpheaderrule) array_ | SetHeader sets HTTP header fields |  |  |


#### HTTPReturn







_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `content` _[HTTPReturnContent](#httpreturncontent)_ | Content is a full HTTP response specifying the errorfile to use, or the response payload specifying the file or the string to use. |  |  |


#### HTTPReturnContent







_Appears in:_
- [HTTPReturn](#httpreturn)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `type` _string_ | Type specifies the content-type of the HTTP response. |  |  |
| `format` _string_ | ContentFormat defines the format of the Content. Can be one an errorfile or a string. |  | Enum: [default-errorfile errorfile errorfiles file lf-file string lf-string] <br /> |
| `value` _string_ | Value specifying the file or the string to use. |  |  |


#### HashType







_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `method` _string_ |  |  | Enum: [map-based consistent] <br /> |
| `function` _string_ |  |  | Enum: [sdbm djb2 wt6 crc32] <br /> |
| `modifier` _string_ |  |  | Enum: [avalanche] <br /> |


#### Hold







_Appears in:_
- [ResolverSpec](#resolverspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `nx` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Nx defines interval between two successive name resolution when the last answer was nx. |  |  |
| `obsolete` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Obsolete defines interval between two successive name resolution when the last answer was obsolete. |  |  |
| `other` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Other defines interval between two successive name resolution when the last answer was other. |  |  |
| `refused` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Refused defines interval between two successive name resolution when the last answer was nx. |  |  |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Timeout defines interval between two successive name resolution when the last answer was timeout. |  |  |
| `valid` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Valid defines interval between two successive name resolution when the last answer was valid. |  |  |


#### Listen



Listen is the Schema for the frontends API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1` | | |
| `kind` _string_ | `Listen` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ListenSpec](#listenspec)_ |  |  |  |


#### ListenSpec



ListenSpec defines the desired state of Listen



_Appears in:_
- [Listen](#listen)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. | http | Enum: [http tcp] <br /> |
| `httpResponse` _[HTTPResponseRules](#httpresponserules)_ | HTTPResponse rules define a set of rules which apply to layer 7 processing. |  |  |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |  |  |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |  |  |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform<br />content switching and generally to take decisions based on content extracted<br />from the request, the response or any environmental status |  |  |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel.<br />The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit.<br />More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |  |  |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |  |  |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |  |  |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |  |  |
| `httpLog` _boolean_ | HTTPLog enables HTTP log format which is the most complete and the best suited for HTTP proxies. It provides<br />the same level of information as the TCP format with additional features which<br />are specific to the HTTP protocol. |  |  |
| `tcpLog` _boolean_ | TCPLog enables advanced logging of TCP connections with session state and timers. By default, the log output format<br />is very poor, as it only contains the source and destination addresses, and the instance name. |  |  |
| `binds` _[Bind](#bind) array_ | Binds defines the frontend listening addresses, ports and its configuration. |  | MinItems: 1 <br /> |
| `servers` _[Server](#server) array_ | Servers defines the backend servers and its configuration. |  |  |
| `serverTemplates` _[ServerTemplate](#servertemplate) array_ | ServerTemplates defines the backend server templates and its configuration. |  |  |
| `checkTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | CheckTimeout sets an additional check timeout, but only after a connection has been already<br />established. |  |  |
| `balance` _[Balance](#balance)_ | Balance defines the load balancing algorithm to be used in a backend. |  |  |
| `redispatch` _boolean_ | Redispatch enable or disable session redistribution in case of connection failure |  |  |
| `hashType` _[HashType](#hashtype)_ | HashType Specify a method to use for mapping hashes to servers |  |  |
| `cookie` _[Cookie](#cookie)_ | Cookie enables cookie-based persistence in a backend. |  |  |
| `hostCertificate` _[CertificateListElement](#certificatelistelement)_ | HostCertificate specifies a certificate for that host used in the crt-list of a frontend |  |  |
| `httpCheck` _[HTTPChk](#httpchk)_ | HTTPCheck Enables HTTP protocol to check on the servers health |  |  |
| `tcpCheck` _boolean_ | TCPCheck Perform health checks using tcp-check send/expect sequences |  |  |


#### Nameserver







_Appears in:_
- [ResolverSpec](#resolverspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name specifies a unique name of the nameserver. |  | Pattern: `^[A-Za-z0-9-_.:]+$` <br /> |
| `address` _string_ | Address |  | Pattern: `^[^\s]+$` <br /> |
| `port` _integer_ | Port |  | Maximum: 65535 <br />Minimum: 1 <br /> |




#### OcspFile







_Appears in:_
- [CertificateListElement](#certificatelistelement)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name |  |  |
| `value` _string_ | Value |  |  |


#### ProxyProtocol







_Appears in:_
- [Server](#server)
- [ServerParams](#serverparams)
- [ServerTemplate](#servertemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `v1` _boolean_ | V1 parameter enforces use of the PROXY protocol version 1. |  |  |
| `v2` _[ProxyProtocolV2](#proxyprotocolv2)_ | V2 parameter enforces use of the PROXY protocol version 2. |  |  |
| `v2SSL` _boolean_ | V2SSL parameter add the SSL information extension of the PROXY protocol to the PROXY protocol header. |  |  |
| `v2SSLCN` _boolean_ | V2SSLCN parameter add the SSL information extension of the PROXY protocol to the PROXY protocol header and he SSL information extension<br />along with the Common Name from the subject of the client certificate (if any), is added to the PROXY protocol header. |  |  |


#### ProxyProtocolV2







_Appears in:_
- [ProxyProtocol](#proxyprotocol)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled enables the PROXY protocol version 2. |  |  |
| `options` _[ProxyProtocolV2Options](#proxyprotocolv2options)_ | Options is a list of options to add to the PROXY protocol header. |  |  |


#### ProxyProtocolV2Options







_Appears in:_
- [ProxyProtocolV2](#proxyprotocolv2)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ssl` _boolean_ | Ssl is equivalent to use V2SSL. |  |  |
| `certCn` _boolean_ | CertCn is equivalent to use V2SSLCN. |  |  |
| `sslCipher` _boolean_ | SslCipher is the name of the used cipher. |  |  |
| `certSig` _boolean_ | CertSig is the signature algorithm of the used certificate. |  |  |
| `certKey` _boolean_ | CertKey is the key algorithm of the used certificate. |  |  |
| `authority` _boolean_ | Authority is the host name value passed by the client (only SNI from a TLS) |  |  |
| `crc32C` _boolean_ | Crc32c is the checksum of the PROXYv2 header. |  |  |
| `uniqueID` _boolean_ | UniqueId sends a unique ID generated using the frontend's "unique-id-format" within the PROXYv2 header.<br />This unique-id is primarily meant for "mode tcp". It can lead to unexpected results in "mode http". |  |  |


#### Redirect







_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `code` _integer_ | Code indicates which type of HTTP redirection is desired. |  | Enum: [301 302 303 307 308] <br /> |
| `type` _[RedirectType](#redirecttype)_ | Type selects a mode and value to redirect |  |  |
| `value` _string_ | Value to redirect |  |  |
| `option` _[RedirectOption](#redirectoption)_ | Value to redirect |  |  |


#### RedirectCookie







_Appears in:_
- [RedirectOption](#redirectoption)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name |  |  |
| `value` _string_ | Value |  |  |


#### RedirectOption







_Appears in:_
- [Redirect](#redirect)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `dropQuery` _boolean_ | DropQuery removes the query string from the original URL when performing the concatenation. |  |  |
| `appendSlash` _boolean_ | AppendSlash adds a / character at the end of the URL. |  |  |
| `SetCookie` _[RedirectCookie](#redirectcookie)_ | SetCookie adds header to the redirection. It will be added with NAME (and optionally "=value") |  |  |
| `ClearCookie` _[RedirectCookie](#redirectcookie)_ | ClearCookie is to instruct the browser to delete the cookie. It will be added with NAME (and optionally "=").<br />To add "=" type any string in the value field |  |  |


#### RedirectType







_Appears in:_
- [Redirect](#redirect)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `location` _boolean_ | Location replaces the entire location of a URL. |  |  |
| `insert` _boolean_ | Prefix adds a prefix to the URL's location. |  |  |
| `prefix` _boolean_ | Scheme redirects to a different scheme. |  |  |


#### RegexBackendMapping







_Appears in:_
- [BackendReference](#backendreference)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ | Name to identify the mapping |  |  |
| `parameter` _string_ | Parameter which will be used for the mapping (default: base) | base |  |
| `selector` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#labelselector-v1-meta)_ | LabelSelector to select multiple backends |  |  |


#### ReplacePath







_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `matchRegex` _string_ | MatchRegex is a string pattern used to identify the paths that need to be replaced. |  |  |
| `replaceFmt` _string_ | ReplaceFmt defines the format string used to replace the values that match the pattern. |  |  |


#### Resolver



Resolver is the Schema for the Resolver API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1` | | |
| `kind` _string_ | `Resolver` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[ResolverSpec](#resolverspec)_ |  |  |  |


#### ResolverSpec



ResolverSpec defines the desired state of Resolver



_Appears in:_
- [Resolver](#resolver)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `nameservers` _[Nameserver](#nameserver) array_ | Nameservers used to configure a nameservers. |  |  |
| `acceptedPayloadSize` _integer_ | AcceptedPayloadSize defines the maximum payload size accepted by HAProxy and announced to all the  name servers<br />configured in this resolver. |  | Maximum: 8192 <br />Minimum: 512 <br /> |
| `parseResolvConf` _boolean_ | ParseResolvConf if true, adds all nameservers found in /etc/resolv.conf to this resolvers nameservers list. |  |  |
| `resolveRetries` _integer_ | ResolveRetries defines the number <nb> of queries to send to resolve a server name before giving up. Default value: 3 |  | Minimum: 1 <br /> |
| `hold` _[Hold](#hold)_ | Hold defines the period during which the last name resolution should be kept based on the last resolution status. |  |  |
| `timeouts` _[Timeouts](#timeouts)_ | Timeouts defines timeouts related to name resolution. |  |  |


#### Rule







_Appears in:_
- [BackendSwitchingRule](#backendswitchingrule)
- [Deny](#deny)
- [HTTPDeleteHeaderRule](#httpdeleteheaderrule)
- [HTTPHeaderRule](#httpheaderrule)
- [HTTPPathRule](#httppathrule)
- [Redirect](#redirect)
- [ReplacePath](#replacepath)
- [TCPRequestRule](#tcprequestrule)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |


#### SSL







_Appears in:_
- [Bind](#bind)
- [Server](#server)
- [ServerParams](#serverparams)
- [ServerTemplate](#servertemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled enables SSL deciphering on connections instantiated from this listener. A<br />certificate is necessary. All contents in the buffers will<br />appear in clear text, so that ACLs and HTTP processing will only have access<br />to deciphered contents. SSLv3 is disabled per default, set MinVersion to SSLv3<br />to enable it. |  |  |
| `minVersion` _string_ | MinVersion enforces use of the specified version or upper on SSL connections<br />instantiated from this listener. |  | Enum: [SSLv3 TLSv1.0 TLSv1.1 TLSv1.2 TLSv1.3] <br /> |
| `verify` _string_ | Verify is only available when support for OpenSSL was built in. If set<br />to 'none', client certificate is not requested. This is the default. In other<br />cases, a client certificate is requested. If the client does not provide a<br />certificate after the request and if 'Verify' is set to 'required', then the<br />handshake is aborted, while it would have succeeded if set to 'optional'. The verification<br />of the certificate provided by the client using CAs from CACertificate.<br />On verify failure the handshake abortes, regardless of the 'verify' option. |  | Enum: [none optional required] <br /> |
| `caCertificate` _[SSLCertificate](#sslcertificate)_ | CACertificate configures the CACertificate used for the Server or Bind client certificate |  |  |
| `certificate` _[SSLCertificate](#sslcertificate)_ | Certificate configures a PEM based Certificate file containing both the required certificates and any<br />associated private keys. |  |  |
| `sni` _string_ | SNI parameter evaluates the sample fetch expression, converts it to a<br />string and uses the result as the host name sent in the SNI TLS extension to<br />the server. |  |  |
| `alpn` _string array_ | Alpn enables the TLS ALPN extension and advertises the specified protocol<br />list as supported on top of ALPN. |  |  |


#### SSLCertificate







_Appears in:_
- [CertificateListElement](#certificatelistelement)
- [GlobalConfiguration](#globalconfiguration)
- [SSL](#ssl)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  |  |
| `value` _string_ |  |  |  |
| `valueFrom` _[SSLCertificateValueFrom](#sslcertificatevaluefrom) array_ |  |  |  |


#### SSLCertificateValueFrom







_Appears in:_
- [SSLCertificate](#sslcertificate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `configMapKeyRef` _[ConfigMapKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#configmapkeyselector-v1-core)_ | ConfigMapKeyRef selects a key of a ConfigMap |  |  |
| `secretKeyRef` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#secretkeyselector-v1-core)_ | SecretKeyRef selects a key of a secret in the pod namespace |  |  |


#### Server







_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |  |  |
| `weight` _integer_ | Weight parameter is used to adjust the server weight relative to<br />other servers. All servers will receive a load proportional to their weight<br />relative to the sum of all weights. |  | Maximum: 256 <br />Minimum: 0 <br /> |
| `check` _[Check](#check)_ | Check configures the health checks of the server. |  |  |
| `initAddr` _string_ | InitAddr indicates in what order the server address should be resolved upon startup if it uses an FQDN.<br />Attempts are made to resolve the address by applying in turn each of the methods mentioned in the comma-delimited<br />list. The first method which succeeds is used. |  |  |
| `resolvers` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Resolvers points to an existing resolvers to resolve current server hostname. |  |  |
| `sendProxy` _boolean_ | SendProxy enforces use of the PROXY protocol over any<br />connection established to this server. The PROXY protocol informs the other<br />end about the layer 3/4 addresses of the incoming connection, so that it can<br />know the client address or the public address it accessed to, whatever the<br />upper layer protocol. |  |  |
| `SendProxyV2` _[ProxyProtocol](#proxyprotocol)_ | SendProxyV2 preparing new update. |  |  |
| `verifyHost` _string_ | VerifyHost is only available when support for OpenSSL was built in, and<br />only takes effect if pec.ssl.verify' is set to 'required'. This directive sets<br />a default static hostname to check the server certificate against when no<br />SNI was used to connect to the server. |  |  |
| `checkSNI` _string_ | CheckSNI This option allows you to specify the SNI to be used when doing health checks over SSL |  |  |
| `cookie` _boolean_ | Cookie sets the cookie value assigned to the server. |  |  |
| `name` _string_ | Name of the server. |  |  |
| `address` _string_ | Address can be a host name, an IPv4 address, an IPv6 address. |  | Pattern: `^[^\s]+$` <br /> |
| `port` _integer_ | Port |  | Maximum: 65535 <br />Minimum: 1 <br /> |


#### ServerParams







_Appears in:_
- [Server](#server)
- [ServerTemplate](#servertemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |  |  |
| `weight` _integer_ | Weight parameter is used to adjust the server weight relative to<br />other servers. All servers will receive a load proportional to their weight<br />relative to the sum of all weights. |  | Maximum: 256 <br />Minimum: 0 <br /> |
| `check` _[Check](#check)_ | Check configures the health checks of the server. |  |  |
| `initAddr` _string_ | InitAddr indicates in what order the server address should be resolved upon startup if it uses an FQDN.<br />Attempts are made to resolve the address by applying in turn each of the methods mentioned in the comma-delimited<br />list. The first method which succeeds is used. |  |  |
| `resolvers` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Resolvers points to an existing resolvers to resolve current server hostname. |  |  |
| `sendProxy` _boolean_ | SendProxy enforces use of the PROXY protocol over any<br />connection established to this server. The PROXY protocol informs the other<br />end about the layer 3/4 addresses of the incoming connection, so that it can<br />know the client address or the public address it accessed to, whatever the<br />upper layer protocol. |  |  |
| `SendProxyV2` _[ProxyProtocol](#proxyprotocol)_ | SendProxyV2 preparing new update. |  |  |
| `verifyHost` _string_ | VerifyHost is only available when support for OpenSSL was built in, and<br />only takes effect if pec.ssl.verify' is set to 'required'. This directive sets<br />a default static hostname to check the server certificate against when no<br />SNI was used to connect to the server. |  |  |
| `checkSNI` _string_ | CheckSNI This option allows you to specify the SNI to be used when doing health checks over SSL |  |  |
| `cookie` _boolean_ | Cookie sets the cookie value assigned to the server. |  |  |


#### ServerTemplate







_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |  |  |
| `weight` _integer_ | Weight parameter is used to adjust the server weight relative to<br />other servers. All servers will receive a load proportional to their weight<br />relative to the sum of all weights. |  | Maximum: 256 <br />Minimum: 0 <br /> |
| `check` _[Check](#check)_ | Check configures the health checks of the server. |  |  |
| `initAddr` _string_ | InitAddr indicates in what order the server address should be resolved upon startup if it uses an FQDN.<br />Attempts are made to resolve the address by applying in turn each of the methods mentioned in the comma-delimited<br />list. The first method which succeeds is used. |  |  |
| `resolvers` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core)_ | Resolvers points to an existing resolvers to resolve current server hostname. |  |  |
| `sendProxy` _boolean_ | SendProxy enforces use of the PROXY protocol over any<br />connection established to this server. The PROXY protocol informs the other<br />end about the layer 3/4 addresses of the incoming connection, so that it can<br />know the client address or the public address it accessed to, whatever the<br />upper layer protocol. |  |  |
| `SendProxyV2` _[ProxyProtocol](#proxyprotocol)_ | SendProxyV2 preparing new update. |  |  |
| `verifyHost` _string_ | VerifyHost is only available when support for OpenSSL was built in, and<br />only takes effect if pec.ssl.verify' is set to 'required'. This directive sets<br />a default static hostname to check the server certificate against when no<br />SNI was used to connect to the server. |  |  |
| `checkSNI` _string_ | CheckSNI This option allows you to specify the SNI to be used when doing health checks over SSL |  |  |
| `cookie` _boolean_ | Cookie sets the cookie value assigned to the server. |  |  |
| `prefix` _string_ | Prefix for the server names to be built. |  | Pattern: `^[^\s]+$` <br /> |
| `numMin` _integer_ | NumMin is the min number of servers as server name suffixes this template initializes. |  |  |
| `num` _integer_ | Num is the max number of servers as server name suffixes this template initializes. |  |  |
| `fqdn` _string_ | FQDN for all the servers this template initializes. |  |  |
| `port` _integer_ | Port |  | Maximum: 65535 <br />Minimum: 1 <br /> |


#### StaticHTTPFile







_Appears in:_
- [ErrorFile](#errorfile)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `name` _string_ |  |  |  |
| `value` _string_ |  |  |  |
| `valueFrom` _[ErrorFileValueFrom](#errorfilevaluefrom)_ |  |  |  |




#### TCPRequestRule







_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditionType` _string_ | ConditionType specifies the type of the condition matching ('if' or 'unless') |  | Enum: [if unless] <br /> |
| `condition` _string_ | Condition is a condition composed of ACLs. |  |  |
| `type` _string_ | Type specifies the type of the tcp-request rule. |  | Enum: [connection content inspect-delay session] <br /> |
| `action` _string_ | Action defines the action to perform if the condition applies. |  | Enum: [accept capture do-resolve expect-netscaler-cip expect-proxy reject sc-inc-gpc0 sc-inc-gpc1 sc-set-gpt0 send-spoe-group set-dst-port set-dst set-priority set-src set-var silent-drop track-sc0 track-sc1 track-sc2 unset-var use-service lua] <br /> |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Timeout sets timeout for the action |  |  |


#### Timeouts







_Appears in:_
- [ResolverSpec](#resolverspec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `resolve` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Resolve time to trigger name resolutions when no other time applied. Default value: 1s |  |  |
| `retry` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Retry time between two DNS queries, when no valid response have been received. Default value: 1s |  |  |



## proxy.haproxy.com/v1alpha1

Package v1alpha1 contains API Schema definitions for the proxy v1alpha1 API group

### Resource Types
- [Instance](#instance)



#### Configuration







_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `global` _[GlobalConfiguration](#globalconfiguration)_ | Global contains the global HAProxy configuration settings |  |  |
| `defaults` _[DefaultsConfiguration](#defaultsconfiguration)_ | Defaults presets settings for all frontend, backend and listen |  |  |
| `selector` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#labelselector-v1-meta)_ | LabelSelector to select other configuration objects of the config.haproxy.com API |  |  |


#### DefaultsConfiguration







_Appears in:_
- [Configuration](#configuration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In tcp mode it is a layer 4 proxy. In http mode it is a layer 7 proxy. | http | Enum: [http tcp] <br /> |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |  |  |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta))_ | Timeouts: check, client, client-fin, connect, http-keep-alive, http-request, queue, server, server-fin, tunnel.<br />The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit.<br />More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html | \{ client:5s connect:5s server:10s \} |  |
| `logging` _[DefaultsLoggingConfiguration](#defaultsloggingconfiguration)_ | Logging is used to configure default logging for all proxies. |  |  |
| `additionalParameters` _string_ | AdditionalParameters can be used to specify any further configuration statements which are not covered in this section explicitly. |  |  |


#### DefaultsLoggingConfiguration







_Appears in:_
- [DefaultsConfiguration](#defaultsconfiguration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled will enable logs for all proxies |  |  |
| `httpLog` _boolean_ | HTTPLog enables HTTP log format which is the most complete and the best suited for HTTP proxies. It provides<br />the same level of information as the TCP format with additional features which<br />are specific to the HTTP protocol. |  |  |
| `tcpLog` _boolean_ | TCPLog enables advanced logging of TCP connections with session state and timers. By default, the log output format<br />is very poor, as it only contains the source and destination addresses, and the instance name. |  |  |


#### GlobalConfiguration







_Appears in:_
- [Configuration](#configuration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `reload` _boolean_ | Reload enables auto-reload of the configuration using sockets. Requires an image that supports this feature. | false |  |
| `statsTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | StatsTimeout sets the timeout on the stats socket. Default is set to 10 seconds. |  |  |
| `logging` _[GlobalLoggingConfiguration](#globalloggingconfiguration)_ | Logging is used to enable and configure logging in the global section of the HAProxy configuration. |  |  |
| `additionalParameters` _string_ | AdditionalParameters can be used to specify any further configuration statements which are not covered in this section explicitly. |  |  |
| `additionalCertificates` _[SSLCertificate](#sslcertificate) array_ | AdditionalCertificates can be used to include global ssl certificates which can bes used in any listen |  |  |
| `maxconn` _integer_ | Maxconn sets the maximum per-process number of concurrent connections. Proxies will stop accepting connections when this limit is reached. |  |  |
| `nbthread` _integer_ | Nbthread this setting is only available when support for threads was built in. It makes HAProxy run on specified number of threads. |  |  |
| `tune` _[GlobalTuneOptions](#globaltuneoptions)_ | TuneOptions sets the global tune options. |  |  |
| `ssl` _[GlobalSSL](#globalssl)_ | GlobalSSL sets the global SSL options. |  |  |
| `hardStopAfter` _[Duration](#duration)_ | HardStopAfter is the maximum time the instance will remain alive when a soft-stop is received. |  |  |
| `ocsp` _[GlobalOCSPConfiguration](#globalocspconfiguration)_ | Ocsp is used to enable stapling at the global level for all certificates in the configuration. |  |  |


#### GlobalLoggingConfiguration







_Appears in:_
- [GlobalConfiguration](#globalconfiguration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled will toggle the creation of a global syslog server. |  |  |
| `address` _string_ | Address can be a filesystem path to a UNIX domain socket or a remote syslog target (IPv4/IPv6 address optionally followed by a colon and a UDP port). | /var/lib/rsyslog/rsyslog.sock | Pattern: `^[^\s]+$` <br /> |
| `facility` _string_ | Facility must be one of the 24 standard syslog facilities. | local0 | Enum: [kern user mail daemon auth syslog lpr news uucp cron auth2 ftp ntp audit alert cron2 local0 local1 local2 local3 local4 local5 local6 local7] <br /> |
| `level` _string_ | Level can be specified to filter outgoing messages. By default, all messages are sent. |  | Enum: [emerg alert crit err warning notice info debug] <br /> |
| `format` _string_ | Format is the log format used when generating syslog messages. |  | Enum: [rfc3164 rfc5424 short raw] <br /> |
| `sendHostname` _boolean_ | SendHostname sets the hostname field in the syslog header.  Generally used if one is not relaying logs through an<br />intermediate syslog server. |  |  |
| `hostname` _string_ | Hostname specifies a value for the syslog hostname header, otherwise uses the hostname of the system. |  |  |


#### GlobalOCSPConfiguration







_Appears in:_
- [GlobalConfiguration](#globalconfiguration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `mode` _boolean_ | Mode Enable automatic OCSP response update when set to 'on', disable it otherwise.<br />Its value defaults to 'off'. |  |  |
| `maxDelay` _integer_ | MaxDelay sets the maximum interval between two automatic updates of the same OCSP<br />response. This time is expressed in seconds and defaults to 3600 (1 hour). |  |  |
| `minDelay` _integer_ | MinDelay sets the minimum interval between two automatic updates of the same OCSP<br />response. This time is expressed in seconds and defaults to 300 (5 minutes). |  |  |
| `httpproxy` _[OcspUpdateOptionsHttpproxy](#ocspupdateoptionshttpproxy)_ | HttpProxy Allow to use an HTTP proxy for the OCSP updates. This only works with HTTP,<br />HTTPS is not supported. This option will allow the OCSP updater to send<br />absolute URI in the request to the proxy. |  |  |


#### GlobalSSL







_Appears in:_
- [GlobalConfiguration](#globalconfiguration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `defaultBindCiphers` _string array_ | DefaultBindCiphers sets the list of cipher algorithms ("cipher suite") that are negotiated during the SSL/TLS handshake up to TLSv1.2 for all<br />binds which do not explicitly define theirs. |  |  |
| `defaultBindCipherSuites` _string array_ | DefaultBindCipherSuites sets the default list of cipher algorithms ("cipher suite") that are negotiated<br />during the TLSv1.3 handshake for all binds which do not explicitly define theirs. |  |  |
| `defaultBindOptions` _[GlobalSSLDefaultBindOptions](#globalssldefaultbindoptions)_ | DefaultBindOptions sets default ssl-options to force on all binds. |  |  |


#### GlobalSSLDefaultBindOptions







_Appears in:_
- [GlobalSSL](#globalssl)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `minVersion` _string_ | MinVersion enforces use of the specified version or upper on SSL connections<br />instantiated from this listener. |  | Enum: [SSLv3 TLSv1.0 TLSv1.1 TLSv1.2 TLSv1.3] <br /> |


#### GlobalSSLTuneOptions







_Appears in:_
- [GlobalTuneOptions](#globaltuneoptions)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `cacheSize` _integer_ | CacheSize sets the size of the global SSL session cache, in a number of blocks. A block<br />is large enough to contain an encoded session without peer certificate.  An<br />encoded session with peer certificate is stored in multiple blocks depending<br />on the size of the peer certificate. The default value may be forced<br />at build time, otherwise defaults to 20000.  Setting this value to 0 disables the SSL session cache. |  |  |
| `keylog` _string_ | Keylog activates the logging of the TLS keys. It should be used with<br />care as it will consume more memory per SSL session and could decrease<br />performances. This is disabled by default. |  |  |
| `lifetime` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#duration-v1-meta)_ | Lifetime sets how long a cached SSL session may remain valid. This time defaults to 5 min. It is important<br />to understand that it does not guarantee that sessions will last that long, because if the cache is<br />full, the longest idle sessions will be purged despite their configured lifetime. |  |  |
| `forcePrivateCache` _boolean_ | ForcePrivateCache disables SSL session cache sharing between all processes. It<br />should normally not be used since it will force many renegotiations due to<br />clients hitting a random process. |  |  |
| `maxRecord` _integer_ | MaxRecord sets the maximum amount of bytes passed to SSL_write() at a time. Default<br />value 0 means there is no limit. Over SSL/TLS, the client can decipher the<br />data only once it has received a full record. |  |  |
| `defaultDHParam` _integer_ | DefaultDHParam sets the maximum size of the Diffie-Hellman parameters used for generating<br />the ephemeral/temporary Diffie-Hellman key in case of DHE key exchange. The<br />final size will try to match the size of the server's RSA (or DSA) key (e.g,<br />a 2048 bits temporary DH key for a 2048 bits RSA key), but will not exceed<br />this maximum value. Default value if 2048. |  |  |
| `ctxCacheSize` _integer_ | CtxCacheSize sets the size of the cache used to store generated certificates to <number><br />entries. This is an LRU cache. Because generating an SSL certificate<br />dynamically is expensive, they are cached. The default cache size is set to 1000 entries. |  |  |
| `captureBufferSize` _integer_ | CaptureBufferSize sets the maximum size of the buffer used for capturing client hello cipher<br />list, extensions list, elliptic curves list and elliptic curve point<br />formats. If the value is 0 (default value) the capture is disabled,<br />otherwise a buffer is allocated for each SSL/TLS connection. |  |  |


#### GlobalTuneOptions







_Appears in:_
- [GlobalConfiguration](#globalconfiguration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `maxrewrite` _integer_ | Maxrewrite sets the reserved buffer space to this size in bytes. The reserved space is<br />used for header rewriting or appending. The first reads on sockets will never<br />fill more than bufsize-maxrewrite. |  |  |
| `buffers_limit` _integer_ | BuffersLimit Sets a hard limit on the number of buffers which may be allocated per process.<br />The default value is zero which means unlimited. The limit will automatically<br />be re-adjusted to satisfy the reserved buffers for emergency situations so<br />that the user doesn't have to perform complicated calculations. |  |  |
| `bufsize` _integer_ | Bufsize sets the buffer size to this size (in bytes). Lower values allow more<br />sessions to coexist in the same amount of RAM, and higher values allow some<br />applications with very large cookies to work. |  |  |
| `buffers_reserve` _integer_ | BuffersReserve Sets the number of per-thread buffers which are pre-allocated and<br />reserved for use only during memory shortage conditions resulting in failed memory<br />allocations. The minimum value is 2 and the default is 4. |  |  |
| `ssl` _[GlobalSSLTuneOptions](#globalssltuneoptions)_ | SSL sets the SSL tune options. |  |  |


#### Instance



Instance is the Schema for the instances API





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `proxy.haproxy.com/v1alpha1` | | |
| `kind` _string_ | `Instance` | | |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[InstanceSpec](#instancespec)_ |  |  |  |




#### InstanceSpec



InstanceSpec defines the desired state of Instance



_Appears in:_
- [Instance](#instance)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `replicas` _integer_ | Replicas is the desired number of replicas of the HAProxy Instance. | 1 |  |
| `network` _[Network](#network)_ | Network contains the configuration of Route, Services and other network related configuration. |  |  |
| `configuration` _[Configuration](#configuration)_ | Configuration is used to bootstrap the global and defaults section of the HAProxy configuration. |  |  |
| `rolloutOnConfigChange` _boolean_ | RolloutOnConfigChange enable rollout on config changes |  |  |
| `image` _string_ | Image specifies the HaProxy image including th tag. | haproxy:latest |  |
| `resources` _[ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)_ | Resources defines the resource requirements for the HAProxy pods. |  |  |
| `sidecars` _[Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#container-v1-core) array_ | Sidecars additional sidecar containers |  |  |
| `serviceAccountName` _string_ | ServiceAccountName is the name of the ServiceAccount to use to run this Instance. |  |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#localobjectreference-v1-core) array_ | ImagePullSecrets is an optional list of secret names in the same namespace to use for pulling any of the images used. |  |  |
| `allowPrivilegedPorts` _boolean_ | AllowPrivilegedPorts allows to bind sockets with port numbers less than 1024. |  |  |
| `placement` _[Placement](#placement)_ | Placement define how the instance's pods should be scheduled. |  |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#pullpolicy-v1-core)_ | ImagePullPolicy one of Always, Never, IfNotPresent. |  |  |
| `metrics` _[Metrics](#metrics)_ | Metrics defines the metrics endpoint and scraping configuration. |  |  |
| `labels` _object (keys:string, values:string)_ | Labels additional labels for the ha-proxy pods |  |  |
| `env` _object (keys:string, values:string)_ | Env additional environment variables |  |  |
| `readinessProbe` _[Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#probe-v1-core)_ | ReadinessProbe the readiness probe for the main container |  |  |
| `livenessProbe` _[Probe](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#probe-v1-core)_ | LivenessProbe the liveness probe for the main container |  |  |
| `podDisruptionBudget` _[PodDisruptionBudget](#poddisruptionbudget)_ | PodDisruptionBudget defines pod disruptions options |  |  |


#### Metrics







_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled will enable metrics globally for Instance. |  |  |
| `address` _string_ | Address to bind the metrics endpoint (default: '0.0.0.0'). | 0.0.0.0 |  |
| `port` _integer_ | Port specifies the port used for metrics. |  |  |
| `relabelings` _RelabelConfig array_ | RelabelConfigs to apply to samples before scraping.<br />More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config |  |  |
| `interval` _[Duration](#duration)_ | Interval at which metrics should be scraped<br />If not specified Prometheus' global scrape interval is used. |  |  |


#### Network







_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `hostNetwork` _boolean_ | HostNetwork will enable the usage of host network. |  |  |
| `hostIPs` _object (keys:string, values:string)_ | HostIPs defines an environment variable BIND_ADDRESS in the instance based on the provided host to IP mapping |  |  |
| `route` _[RouteSpec](#routespec)_ | Route defines the desired state for OpenShift Routes. |  |  |
| `service` _[ServiceSpec](#servicespec)_ | Service defines the desired state for a Service. |  |  |


#### OcspUpdateOptionsHttpproxy







_Appears in:_
- [GlobalOCSPConfiguration](#globalocspconfiguration)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `address` _string_ | Address can be a host name, an IPv4 address or an IPv6 address |  | Pattern: `^[^\s]+$` <br /> |
| `port` _integer_ | Port |  |  |


#### Placement







_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `nodeSelector` _object (keys:string, values:string)_ | NodeSelector is a selector which must be true for the pod to fit on a node. |  |  |
| `topologySpreadConstraints` _[TopologySpreadConstraint](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#topologyspreadconstraint-v1-core) array_ | TopologySpreadConstraints describes how a group of pods ought to spread across topology<br />domains. Scheduler will schedule pods in a way which abides by the constraints. |  |  |


#### PodDisruptionBudget







_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `minAvailable` _[IntOrString](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#intorstring-intstr-util)_ | An eviction is allowed if at least “minAvailable“ pods selected by “selector” will still be available after the eviction |  |  |
| `maxUnavailable` _[IntOrString](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#intorstring-intstr-util)_ | An eviction is allowed if at most “maxUnavailable“ pods selected by “selector” are unavailable after the eviction |  |  |


#### RouteSpec







_Appears in:_
- [Network](#network)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled will toggle the creation of OpenShift Routes. |  |  |
| `tls` _[TLSConfig](#tlsconfig)_ | TLS provides the ability to configure certificates and termination for the route. |  |  |


#### ServiceSpec







_Appears in:_
- [Network](#network)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `enabled` _boolean_ | Enabled will toggle the creation of a Service. |  |  |
| `type` _[ServiceType](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#servicetype-v1-core)_ | Type will define the Service Type. | ClusterIP | Enum: [ClusterIP NodePort LoadBalancer] <br /> |
| `annotations` _object (keys:string, values:string)_ | Annotations to be added to Service. |  |  |


