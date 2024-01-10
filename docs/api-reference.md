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

| Field | Description |
| --- | --- |
| `name` _string_ | Name |
| `criterion` _string_ | Criterion is the name of a sample fetch method, or one of its ACL specific declinations. |
| `values` _string array_ | Values are of the type supported by the criterion. |


#### Backend



Backend is the Schema for the backend API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1`
| `kind` _string_ | `Backend`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[BackendSpec](#backendspec)_ |  |


#### BackendReference





_Appears in:_
- [BackendSwitchingRule](#backendswitchingrule)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of a specific backend |
| `regexMapping` _[RegexBackendMapping](#regexbackendmapping)_ | Mapping of multiple backends |


#### BackendSpec



BackendSpec defines the desired state of Backend

_Appears in:_
- [Backend](#backend)

| Field | Description |
| --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform content switching and generally to take decisions based on content extracted from the request, the response or any environmental status |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel. The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit. More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |
| `checkTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | CheckTimeout sets an additional check timeout, but only after a connection has been already established. |
| `servers` _[Server](#server) array_ | Servers defines the backend servers and its configuration. |
| `serverTemplates` _[ServerTemplate](#servertemplate) array_ | ServerTemplates defines the backend server templates and its configuration. |
| `balance` _[Balance](#balance)_ | Balance defines the load balancing algorithm to be used in a backend. |
| `hostRegex` _string_ | HostRegex specifies a regular expression used for backend switching rules. |
| `hostCertificate` _[CertificateListElement](#certificatelistelement)_ | HostCertificate specifies a certificate for that host used in the crt-list of a frontend |
| `redispatch` _boolean_ | Redispatch enable or disable session redistribution in case of connection failure |
| `hashType` _[HashType](#hashtype)_ | HashType specifies a method to use for mapping hashes to servers |
| `cookie` _[Cookie](#cookie)_ | Cookie enables cookie-based persistence in a backend. |


#### BackendSwitchingRule





_Appears in:_
- [FrontendSpec](#frontendspec)

| Field | Description |
| --- | --- |
| `backend` _[BackendReference](#backendreference)_ | Backend reference used to resolve the backend name. |


#### Balance





_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `algorithm` _string_ | Algorithm is the algorithm used to select a server when doing load balancing. This only applies when no persistence information is available, or when a connection is redispatched to another server. |


#### BaseSpec





_Appears in:_
- [BackendSpec](#backendspec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform content switching and generally to take decisions based on content extracted from the request, the response or any environmental status |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel. The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit. More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |


#### Bind





_Appears in:_
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name for these sockets, which will be reported on the stats page. |
| `address` _string_ | Address can be a host name, an IPv4 address, an IPv6 address, or '*' (is equal to the special address "0.0.0.0"). |
| `port` _integer_ | Port |
| `portRangeEnd` _[int64](#int64)_ | PortRangeEnd if set it must be greater than Port |
| `transparent` _boolean_ | Transparent is an optional keyword which is supported only on certain Linux kernels. It indicates that the addresses will be bound even if they do not belong to the local machine, and that packets targeting any of these addresses will be intercepted just as if the addresses were locally configured. This normally requires that IP forwarding is enabled. Caution! do not use this with the default address '*', as it would redirect any traffic for the specified port. |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |
| `sslCertificateList` _[CertificateList](#certificatelist)_ | This setting is only available when support for OpenSSL was built in. It designates a list of PEM file with an optional ssl configuration and a SNI filter per certificate. |
| `hidden` _boolean_ | Hidden hides the bind and prevent exposing the Bind in services or routes |
| `acceptProxy` _boolean_ | AcceptProxy enforces the use of the PROXY protocol over any connection accepted by any of the sockets declared on the same line. |


#### CertificateListElement





_Appears in:_
- [BackendSpec](#backendspec)
- [CertificateList](#certificatelist)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `certificate` _[SSLCertificate](#sslcertificate)_ | Certificate that will be presented to clients who provide a valid TLSServerNameIndication field matching the SNIFilter. |
| `sniFilter` _string_ | SNIFilter specifies the filter for the SSL Certificate.  Wildcards are supported in the SNIFilter. Negative filter are also supported. |
| `alpn` _string array_ | Alpn enables the TLS ALPN extension and advertises the specified protocol list as supported on top of ALPN. |


#### Check





_Appears in:_
- [Server](#server)
- [ServerParams](#serverparams)
- [ServerTemplate](#servertemplate)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enable enables health checks on a server. If not set, no health checking is performed, and the server is always considered available. |
| `inter` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Inter sets the interval between two consecutive health checks. If left unspecified, the delay defaults to 2000 ms. |
| `rise` _[int64](#int64)_ | Rise specifies the number of consecutive successful health checks after a server will be considered as operational. This value defaults to 2 if unspecified. |
| `fall` _[int64](#int64)_ | Fall specifies the number of consecutive unsuccessful health checks after a server will be considered as dead. This value defaults to 3 if unspecified. |


#### Cookie





_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name of the cookie which will be monitored, modified or inserted in order to bring persistence. |
| `mode` _[CookieMode](#cookiemode)_ | Mode could be 'rewrite', 'insert', 'prefix'. Select one. |
| `indirect` _boolean_ | Indirect no cookie will be emitted to a client which already has a valid one for the server which has processed the request. |
| `noCache` _boolean_ | NoCache recommended in conjunction with the insert mode when there is a cache between the client and HAProx |
| `postOnly` _boolean_ | PostOnly ensures that cookie insertion will only be performed on responses to POST requests. |
| `preserve` _boolean_ | Preserve only be used with "insert" and/or "indirect". It allows the server to emit the persistence cookie itself. |
| `httpOnly` _boolean_ | HTTPOnly add an "HttpOnly" cookie attribute when a cookie is inserted. It doesn't share the cookie with non-HTTP components. |
| `secure` _boolean_ | Secure add a "Secure" cookie attribute when a cookie is inserted. The user agent never emits this cookie over non-secure channels. The cookie will be presented only over SSL/TLS connections. |
| `dynamic` _boolean_ | Dynamic activates dynamic cookies, when used, a session cookie is dynamically created for each server, based on the IP and port of the server, and a secret key. |
| `domain` _string array_ | Domain specify the domain at which a cookie is inserted. You can specify several domain names by invoking this option multiple times. |
| `maxIdle` _integer_ | MaxIdle cookies are ignored after some idle time. |
| `maxLife` _integer_ | MaxLife cookies are ignored after some life time. |
| `attribute` _string array_ | Attribute add an extra attribute when a cookie is inserted. |


#### CookieMode





_Appears in:_
- [Cookie](#cookie)

| Field | Description |
| --- | --- |
| `rewrite` _boolean_ | Rewrite the cookie will be provided by the server. |
| `insert` _boolean_ | Insert cookie will have to be inserted by haproxy in server responses. |
| `prefix` _boolean_ | Prefix is needed in some specific environments where the client does not support more than one single cookie and the application already needs it. |


#### Deny





_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enabled enables deny http request |


#### ErrorFile





_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [DefaultsConfiguration](#defaultsconfiguration)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `code` _integer_ | Code is the HTTP status code. |
| `file` _[StaticHTTPFile](#statichttpfile)_ | File designates a file containing the full HTTP response. |


#### ErrorFileValueFrom

_Underlying type:_ _[struct{ConfigMapKeyRef *k8s.io/api/core/v1.ConfigMapKeySelector "json:\"configMapKeyRef,omitempty\""}](#struct{configmapkeyref-*k8sioapicorev1configmapkeyselector-"json:\"configmapkeyref,omitempty\""})_



_Appears in:_
- [StaticHTTPFile](#statichttpfile)



#### Forwardfor





_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ |  |
| `except` _string_ | Pattern: ^[^\s]+$ |
| `header` _string_ | Pattern: ^[^\s]+$ |
| `ifnone` _boolean_ |  |


#### Frontend



Frontend is the Schema for the frontends API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1`
| `kind` _string_ | `Frontend`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[FrontendSpec](#frontendspec)_ |  |


#### FrontendSpec



FrontendSpec defines the desired state of Frontend

_Appears in:_
- [Frontend](#frontend)

| Field | Description |
| --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform content switching and generally to take decisions based on content extracted from the request, the response or any environmental status |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel. The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit. More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |
| `binds` _[Bind](#bind) array_ | Binds defines the frontend listening addresses, ports and its configuration. |
| `backendSwitching` _[BackendSwitchingRule](#backendswitchingrule) array_ | BackendSwitching rules specify the specific backend used if/unless an ACL-based condition is matched. |
| `defaultBackend` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#localobjectreference-v1-core)_ | DefaultBackend to use when no 'use_backend' rule has been matched. |


#### HTTPHeaderRule





_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description |
| --- | --- |
| `name` _string_ | Name specifies the header name |
| `value` _[HTTPHeaderValue](#httpheadervalue)_ | Value specifies the header value |


#### HTTPHeaderValue

_Underlying type:_ _[struct{Env *k8s.io/api/core/v1.EnvVar "json:\"env,omitempty\""; Str *string "json:\"str,omitempty\""; Format *string "json:\"format,omitempty\""}](#struct{env-*k8sioapicorev1envvar-"json:\"env,omitempty\"";-str-*string-"json:\"str,omitempty\"";-format-*string-"json:\"format,omitempty\""})_



_Appears in:_
- [HTTPHeaderRule](#httpheaderrule)



#### HTTPPathRule





_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description |
| --- | --- |
| `format` _string_ | Value specifies the path value |




#### HTTPRequestRules





_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `setHeader` _[HTTPHeaderRule](#httpheaderrule) array_ | SetHeader sets HTTP header fields |
| `setPath` _[HTTPPathRule](#httppathrule) array_ | SetPath sets request path |
| `addHeader` _[HTTPHeaderRule](#httpheaderrule) array_ | AddHeader appends HTTP header fields |
| `redirect` _[Redirect](#redirect) array_ | Redirect performs an HTTP redirection based on a redirect rule. |
| `replacePath` _[ReplacePath](#replacepath) array_ | ReplacePath matches the value of the path using a regex and completely replaces it with the specified format. The replacement does not modify the scheme, the authority and the query-string. |
| `deny` _[Deny](#deny)_ | Deny stops the evaluation of the rules and immediately rejects the request and emits an HTTP 403 error. Optionally the status code specified as an argument to deny_status. |
| `denyStatus` _[int64](#int64)_ | DenyStatus is the HTTP status code. |
| `return` _[HTTPReturn](#httpreturn)_ | Return stops the evaluation of the rules and immediately returns a response. |


#### HTTPReturn





_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description |
| --- | --- |
| `content` _[HTTPReturnContent](#httpreturncontent)_ | Content is a full HTTP response specifying the errorfile to use, or the response payload specifying the file or the string to use. |


#### HTTPReturnContent

_Underlying type:_ _[struct{Type string "json:\"type\""; Format string "json:\"format\""; Value string "json:\"value\""}](#struct{type-string-"json:\"type\"";-format-string-"json:\"format\"";-value-string-"json:\"value\""})_



_Appears in:_
- [HTTPReturn](#httpreturn)



#### HashType





_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `method` _string_ |  |
| `function` _string_ |  |
| `modifier` _string_ |  |


#### Hold





_Appears in:_
- [ResolverSpec](#resolverspec)

| Field | Description |
| --- | --- |
| `nx` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Nx defines interval between two successive name resolution when the last answer was nx. |
| `obsolete` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Obsolete defines interval between two successive name resolution when the last answer was obsolete. |
| `other` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Other defines interval between two successive name resolution when the last answer was other. |
| `refused` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Refused defines interval between two successive name resolution when the last answer was nx. |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Timeout defines interval between two successive name resolution when the last answer was timeout. |
| `valid` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Valid defines interval between two successive name resolution when the last answer was valid. |


#### Listen



Listen is the Schema for the frontends API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1`
| `kind` _string_ | `Listen`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ListenSpec](#listenspec)_ |  |


#### ListenSpec



ListenSpec defines the desired state of Listen

_Appears in:_
- [Listen](#listen)

| Field | Description |
| --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In TCP mode it is a layer 4 proxy. In HTTP mode it is a layer 7 proxy. |
| `httpRequest` _[HTTPRequestRules](#httprequestrules)_ | HTTPRequest rules define a set of rules which apply to layer 7 processing. |
| `tcpRequest` _[TCPRequestRule](#tcprequestrule) array_ | TCPRequest rules perform an action on an incoming connection depending on a layer 4 condition. |
| `acl` _[ACL](#acl) array_ | ACL (Access Control Lists) provides a flexible solution to perform content switching and generally to take decisions based on content extracted from the request, the response or any environmental status |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta))_ | Timeouts: check, connect, http-keep-alive, http-request, queue, server, tunnel. The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit. More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |
| `forwardFor` _[Forwardfor](#forwardfor)_ | Forwardfor enable insertion of the X-Forwarded-For header to requests sent to servers |
| `httpPretendKeepalive` _boolean_ | HTTPPretendKeepalive will keep the connection alive. It is recommended not to enable this option by default. |
| `binds` _[Bind](#bind) array_ | Binds defines the frontend listening addresses, ports and its configuration. |
| `servers` _[Server](#server) array_ | Servers defines the backend servers and its configuration. |
| `serverTemplates` _[ServerTemplate](#servertemplate) array_ | ServerTemplates defines the backend server templates and its configuration. |
| `checkTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | CheckTimeout sets an additional check timeout, but only after a connection has been already established. |
| `balance` _[Balance](#balance)_ | Balance defines the load balancing algorithm to be used in a backend. |
| `redispatch` _boolean_ | Redispatch enable or disable session redistribution in case of connection failure |
| `hashType` _[HashType](#hashtype)_ | HashType Specify a method to use for mapping hashes to servers |
| `cookie` _[Cookie](#cookie)_ | Cookie enables cookie-based persistence in a backend. |
| `hostCertificate` _[CertificateListElement](#certificatelistelement)_ | HostCertificate specifies a certificate for that host used in the crt-list of a frontend |


#### Nameserver





_Appears in:_
- [ResolverSpec](#resolverspec)

| Field | Description |
| --- | --- |
| `name` _string_ | Name specifies a unique name of the nameserver. |
| `address` _string_ | Address |
| `port` _integer_ | Port |




#### ProxyProtocol





_Appears in:_
- [Server](#server)
- [ServerParams](#serverparams)
- [ServerTemplate](#servertemplate)

| Field | Description |
| --- | --- |
| `v1` _boolean_ | V1 parameter enforces use of the PROXY protocol version 1. |
| `v2` _[ProxyProtocolV2](#proxyprotocolv2)_ | V2 parameter enforces use of the PROXY protocol version 2. |
| `v2SSL` _boolean_ | V2SSL parameter add the SSL information extension of the PROXY protocol to the PROXY protocol header. |
| `v2SSLCN` _boolean_ | V2SSLCN parameter add the SSL information extension of the PROXY protocol to the PROXY protocol header and he SSL information extension along with the Common Name from the subject of the client certificate (if any), is added to the PROXY protocol header. |




#### ProxyProtocolV2Options





_Appears in:_
- [ProxyProtocolV2](#proxyprotocolv2)

| Field | Description |
| --- | --- |
| `ssl` _boolean_ | Ssl is equivalent to use V2SSL. |
| `certCn` _boolean_ | CertCn is equivalent to use V2SSLCN. |
| `sslCipher` _boolean_ | SslCipher is the name of the used cipher. |
| `certSig` _boolean_ | CertSig is the signature algorithm of the used certificate. |
| `certKey` _boolean_ | CertKey is the key algorithm of the used certificate. |
| `authority` _boolean_ | Authority is the host name value passed by the client (only SNI from a TLS) |
| `crc32C` _boolean_ | Crc32c is the checksum of the PROXYv2 header. |
| `uniqueID` _boolean_ | UniqueId sends a unique ID generated using the frontend's "unique-id-format" within the PROXYv2 header. This unique-id is primarily meant for "mode tcp". It can lead to unexpected results in "mode http". |


#### Redirect





_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description |
| --- | --- |
| `code` _[int64](#int64)_ | Code indicates which type of HTTP redirection is desired. |
| `type` _[RedirectType](#redirecttype)_ | Type selects a mode and value to redirect |
| `value` _string_ | Value to redirect |
| `option` _[RedirectOption](#redirectoption)_ | Value to redirect |


#### RedirectCookie





_Appears in:_
- [RedirectOption](#redirectoption)

| Field | Description |
| --- | --- |
| `name` _string_ | Name |
| `value` _string_ | Value |




#### RedirectType

_Underlying type:_ _[struct{Location bool "json:\"location\""; Prefix bool "json:\"insert\""; Scheme bool "json:\"prefix\""}](#struct{location-bool-"json:\"location\"";-prefix-bool-"json:\"insert\"";-scheme-bool-"json:\"prefix\""})_



_Appears in:_
- [Redirect](#redirect)



#### RegexBackendMapping





_Appears in:_
- [BackendReference](#backendreference)

| Field | Description |
| --- | --- |
| `name` _string_ | Name to identify the mapping |
| `parameter` _string_ | Parameter which will be used for the mapping (default: base) |
| `selector` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#labelselector-v1-meta)_ | LabelSelector to select multiple backends |


#### ReplacePath





_Appears in:_
- [HTTPRequestRules](#httprequestrules)

| Field | Description |
| --- | --- |
| `matchRegex` _string_ | MatchRegex is a string pattern used to identify the paths that need to be replaced. |
| `replaceFmt` _string_ | ReplaceFmt defines the format string used to replace the values that match the pattern. |


#### Resolver



Resolver is the Schema for the Resolver API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `config.haproxy.com/v1alpha1`
| `kind` _string_ | `Resolver`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[ResolverSpec](#resolverspec)_ |  |


#### ResolverSpec



ResolverSpec defines the desired state of Resolver

_Appears in:_
- [Resolver](#resolver)

| Field | Description |
| --- | --- |
| `nameservers` _[Nameserver](#nameserver) array_ | Nameservers used to configure a nameservers. |
| `acceptedPayloadSize` _[int64](#int64)_ | AcceptedPayloadSize defines the maximum payload size accepted by HAProxy and announced to all the  name servers configured in this resolver. |
| `parseResolvConf` _boolean_ | ParseResolvConf if true, adds all nameservers found in /etc/resolv.conf to this resolvers nameservers list. |
| `resolveRetries` _[int64](#int64)_ | ResolveRetries defines the number <nb> of queries to send to resolve a server name before giving up. Default value: 3 |
| `hold` _[Hold](#hold)_ | Hold defines the period during which the last name resolution should be kept based on the last resolution status. |
| `timeouts` _[Timeouts](#timeouts)_ | Timeouts defines timeouts related to name resolution. |


#### Rule

_Underlying type:_ _[struct{ConditionType string "json:\"conditionType,omitempty\""; Condition string "json:\"condition,omitempty\""}](#struct{conditiontype-string-"json:\"conditiontype,omitempty\"";-condition-string-"json:\"condition,omitempty\""})_



_Appears in:_
- [BackendSwitchingRule](#backendswitchingrule)
- [Deny](#deny)
- [HTTPHeaderRule](#httpheaderrule)
- [HTTPPathRule](#httppathrule)
- [Redirect](#redirect)
- [ReplacePath](#replacepath)
- [TCPRequestRule](#tcprequestrule)



#### SSL





_Appears in:_
- [Bind](#bind)
- [Server](#server)
- [ServerParams](#serverparams)
- [ServerTemplate](#servertemplate)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enabled enables SSL deciphering on connections instantiated from this listener. A certificate is necessary. All contents in the buffers will appear in clear text, so that ACLs and HTTP processing will only have access to deciphered contents. SSLv3 is disabled per default, set MinVersion to SSLv3 to enable it. |
| `minVersion` _string_ | MinVersion enforces use of the specified version or upper on SSL connections instantiated from this listener. |
| `verify` _string_ | Verify is only available when support for OpenSSL was built in. If set to 'none', client certificate is not requested. This is the default. In other cases, a client certificate is requested. If the client does not provide a certificate after the request and if 'Verify' is set to 'required', then the handshake is aborted, while it would have succeeded if set to 'optional'. The verification of the certificate provided by the client using CAs from CACertificate. On verify failure the handshake abortes, regardless of the 'verify' option. |
| `caCertificate` _[SSLCertificate](#sslcertificate)_ | CACertificate configures the CACertificate used for the Server or Bind client certificate |
| `certificate` _[SSLCertificate](#sslcertificate)_ | Certificate configures a PEM based Certificate file containing both the required certificates and any associated private keys. |
| `sni` _string_ | SNI parameter evaluates the sample fetch expression, converts it to a string and uses the result as the host name sent in the SNI TLS extension to the server. |
| `alpn` _string array_ | Alpn enables the TLS ALPN extension and advertises the specified protocol list as supported on top of ALPN. |


#### SSLCertificate





_Appears in:_
- [CertificateListElement](#certificatelistelement)
- [GlobalConfiguration](#globalconfiguration)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `value` _string_ |  |
| `valueFrom` _[SSLCertificateValueFrom](#sslcertificatevaluefrom) array_ |  |


#### SSLCertificateValueFrom





_Appears in:_
- [SSLCertificate](#sslcertificate)

| Field | Description |
| --- | --- |
| `configMapKeyRef` _[ConfigMapKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#configmapkeyselector-v1-core)_ | ConfigMapKeyRef selects a key of a ConfigMap |
| `secretKeyRef` _[SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#secretkeyselector-v1-core)_ | SecretKeyRef selects a key of a secret in the pod namespace |


#### Server





_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |
| `weight` _[int64](#int64)_ | Weight parameter is used to adjust the server weight relative to other servers. All servers will receive a load proportional to their weight relative to the sum of all weights. |
| `check` _[Check](#check)_ | Check configures the health checks of the server. |
| `initAddr` _string_ | InitAddr indicates in what order the server address should be resolved upon startup if it uses an FQDN. Attempts are made to resolve the address by applying in turn each of the methods mentioned in the comma-delimited list. The first method which succeeds is used. |
| `resolvers` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#localobjectreference-v1-core)_ | Resolvers points to an existing resolvers to resolve current server hostname. |
| `sendProxy` _boolean_ | SendProxy enforces use of the PROXY protocol over any connection established to this server. The PROXY protocol informs the other end about the layer 3/4 addresses of the incoming connection, so that it can know the client address or the public address it accessed to, whatever the upper layer protocol. |
| `SendProxyV2` _[ProxyProtocol](#proxyprotocol)_ | SendProxyV2 preparing new update. |
| `verifyHost` _string_ | VerifyHost is only available when support for OpenSSL was built in, and only takes effect if pec.ssl.verify' is set to 'required'. This directive sets a default static hostname to check the server certificate against when no SNI was used to connect to the server. |
| `cookie` _boolean_ | Cookie sets the cookie value assigned to the server. |
| `name` _string_ | Name of the server. |
| `address` _string_ | Address can be a host name, an IPv4 address, an IPv6 address. |
| `port` _integer_ | Port |


#### ServerParams





_Appears in:_
- [Server](#server)
- [ServerTemplate](#servertemplate)

| Field | Description |
| --- | --- |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |
| `weight` _[int64](#int64)_ | Weight parameter is used to adjust the server weight relative to other servers. All servers will receive a load proportional to their weight relative to the sum of all weights. |
| `check` _[Check](#check)_ | Check configures the health checks of the server. |
| `initAddr` _string_ | InitAddr indicates in what order the server address should be resolved upon startup if it uses an FQDN. Attempts are made to resolve the address by applying in turn each of the methods mentioned in the comma-delimited list. The first method which succeeds is used. |
| `resolvers` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#localobjectreference-v1-core)_ | Resolvers points to an existing resolvers to resolve current server hostname. |
| `sendProxy` _boolean_ | SendProxy enforces use of the PROXY protocol over any connection established to this server. The PROXY protocol informs the other end about the layer 3/4 addresses of the incoming connection, so that it can know the client address or the public address it accessed to, whatever the upper layer protocol. |
| `SendProxyV2` _[ProxyProtocol](#proxyprotocol)_ | SendProxyV2 preparing new update. |
| `verifyHost` _string_ | VerifyHost is only available when support for OpenSSL was built in, and only takes effect if pec.ssl.verify' is set to 'required'. This directive sets a default static hostname to check the server certificate against when no SNI was used to connect to the server. |
| `cookie` _boolean_ | Cookie sets the cookie value assigned to the server. |


#### ServerTemplate





_Appears in:_
- [BackendSpec](#backendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `ssl` _[SSL](#ssl)_ | SSL configures OpenSSL |
| `weight` _[int64](#int64)_ | Weight parameter is used to adjust the server weight relative to other servers. All servers will receive a load proportional to their weight relative to the sum of all weights. |
| `check` _[Check](#check)_ | Check configures the health checks of the server. |
| `initAddr` _string_ | InitAddr indicates in what order the server address should be resolved upon startup if it uses an FQDN. Attempts are made to resolve the address by applying in turn each of the methods mentioned in the comma-delimited list. The first method which succeeds is used. |
| `resolvers` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#localobjectreference-v1-core)_ | Resolvers points to an existing resolvers to resolve current server hostname. |
| `sendProxy` _boolean_ | SendProxy enforces use of the PROXY protocol over any connection established to this server. The PROXY protocol informs the other end about the layer 3/4 addresses of the incoming connection, so that it can know the client address or the public address it accessed to, whatever the upper layer protocol. |
| `SendProxyV2` _[ProxyProtocol](#proxyprotocol)_ | SendProxyV2 preparing new update. |
| `verifyHost` _string_ | VerifyHost is only available when support for OpenSSL was built in, and only takes effect if pec.ssl.verify' is set to 'required'. This directive sets a default static hostname to check the server certificate against when no SNI was used to connect to the server. |
| `cookie` _boolean_ | Cookie sets the cookie value assigned to the server. |
| `prefix` _string_ | Prefix for the server names to be built. |
| `numMin` _[int64](#int64)_ | NumMin is the min number of servers as server name suffixes this template initializes. |
| `num` _integer_ | Num is the max number of servers as server name suffixes this template initializes. |
| `fqdn` _string_ | FQDN for all the servers this template initializes. |
| `port` _integer_ | Port |


#### StaticHTTPFile





_Appears in:_
- [ErrorFile](#errorfile)

| Field | Description |
| --- | --- |
| `name` _string_ |  |
| `value` _string_ |  |
| `valueFrom` _[ErrorFileValueFrom](#errorfilevaluefrom)_ |  |


#### StatusPhase

_Underlying type:_ _string_

StatusPhase is a label for the phase of an object at the current time.

_Appears in:_
- [Status](#status)



#### TCPRequestRule





_Appears in:_
- [BackendSpec](#backendspec)
- [BaseSpec](#basespec)
- [FrontendSpec](#frontendspec)
- [ListenSpec](#listenspec)

| Field | Description |
| --- | --- |
| `type` _string_ | Type specifies the type of the tcp-request rule. |
| `action` _string_ | Action defines the action to perform if the condition applies. |
| `timeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Timeout sets timeout for the action |


#### Timeouts





_Appears in:_
- [ResolverSpec](#resolverspec)

| Field | Description |
| --- | --- |
| `resolve` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Resolve time to trigger name resolutions when no other time applied. Default value: 1s |
| `retry` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | Retry time between two DNS queries, when no valid response have been received. Default value: 1s |



## proxy.haproxy.com/v1alpha1

Package v1alpha1 contains API Schema definitions for the proxy v1alpha1 API group

### Resource Types
- [Instance](#instance)



#### Configuration





_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description |
| --- | --- |
| `global` _[GlobalConfiguration](#globalconfiguration)_ | Global contains the global HAProxy configuration settings |
| `defaults` _[DefaultsConfiguration](#defaultsconfiguration)_ | Defaults presets settings for all frontend, backend and listen |
| `selector` _[LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#labelselector-v1-meta)_ | LabelSelector to select other configuration objects of the config.haproxy.com API |


#### DefaultsConfiguration





_Appears in:_
- [Configuration](#configuration)

| Field | Description |
| --- | --- |
| `mode` _string_ | Mode can be either 'tcp' or 'http'. In tcp mode it is a layer 4 proxy. In http mode it is a layer 7 proxy. |
| `errorFiles` _[ErrorFile](#errorfile) array_ | ErrorFiles custom error files to be used |
| `timeouts` _object (keys:string, values:[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta))_ | Timeouts: check, client, client-fin, connect, http-keep-alive, http-request, queue, server, server-fin, tunnel. The timeout value specified in milliseconds by default, but can be in any other unit if the number is suffixed by the unit. More info: https://cbonte.github.io/haproxy-dconv/2.6/configuration.html |
| `logging` _[DefaultsLoggingConfiguration](#defaultsloggingconfiguration)_ | Logging is used to configure default logging for all proxies. |
| `additionalParameters` _string_ | AdditionalParameters can be used to specify any further configuration statements which are not covered in this section explicitly. |


#### DefaultsLoggingConfiguration





_Appears in:_
- [DefaultsConfiguration](#defaultsconfiguration)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enabled will enable logs for all proxies |
| `httpLog` _boolean_ | HTTPLog enables HTTP log format which is the most complete and the best suited for HTTP proxies. It provides the same level of information as the TCP format with additional features which are specific to the HTTP protocol. |
| `tcpLog` _boolean_ | TCPLog enables advanced logging of TCP connections with session state and timers. By default, the log output format is very poor, as it only contains the source and destination addresses, and the instance name. |


#### GlobalConfiguration





_Appears in:_
- [Configuration](#configuration)

| Field | Description |
| --- | --- |
| `reload` _boolean_ | Reload enables auto-reload of the configuration using sockets. Requires an image that supports this feature. |
| `statsTimeout` _[Duration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#duration-v1-meta)_ | StatsTimeout sets the timeout on the stats socket. Default is set to 10 seconds. |
| `logging` _[GlobalLoggingConfiguration](#globalloggingconfiguration)_ | Logging is used to enable and configure logging in the global section of the HAProxy configuration. |
| `additionalParameters` _string_ | AdditionalParameters can be used to specify any further configuration statements which are not covered in this section explicitly. |
| `additionalCertificates` _[SSLCertificate](#sslcertificate) array_ | AdditionalCertificates can be used to include global ssl certificates which can bes used in any listen |
| `maxconn` _[int64](#int64)_ | Maxconn sets the maximum per-process number of concurrent connections. Proxies will stop accepting connections when this limit is reached. |
| `nbthread` _[int64](#int64)_ | Nbthread this setting is only available when support for threads was built in. It makes HAProxy run on specified number of threads. |
| `tune` _[GlobalTuneOptions](#globaltuneoptions)_ | TuneOptions sets the global tune options. |
| `ssl` _[GlobalSSL](#globalssl)_ | GlobalSSL sets the global SSL options. |
| `hardStopAfter` _[Duration](#duration)_ | HardStopAfter is the maximum time the instance will remain alive when a soft-stop is received. |


#### GlobalLoggingConfiguration





_Appears in:_
- [GlobalConfiguration](#globalconfiguration)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enabled will toggle the creation of a global syslog server. |
| `address` _string_ | Address can be a filesystem path to a UNIX domain socket or a remote syslog target (IPv4/IPv6 address optionally followed by a colon and a UDP port). |
| `facility` _string_ | Facility must be one of the 24 standard syslog facilities. |
| `level` _string_ | Level can be specified to filter outgoing messages. By default, all messages are sent. |
| `format` _string_ | Format is the log format used when generating syslog messages. |
| `sendHostname` _boolean_ | SendHostname sets the hostname field in the syslog header.  Generally used if one is not relaying logs through an intermediate syslog server. |
| `hostname` _string_ | Hostname specifies a value for the syslog hostname header, otherwise uses the hostname of the system. |


#### GlobalSSL





_Appears in:_
- [GlobalConfiguration](#globalconfiguration)

| Field | Description |
| --- | --- |
| `defaultBindCiphers` _string array_ | DefaultBindCiphers sets the list of cipher algorithms ("cipher suite") that are negotiated during the SSL/TLS handshake up to TLSv1.2 for all binds which do not explicitly define theirs. |
| `defaultBindCipherSuites` _string array_ | DefaultBindCipherSuites sets the default list of cipher algorithms ("cipher suite") that are negotiated during the TLSv1.3 handshake for all binds which do not explicitly define theirs. |
| `defaultBindOptions` _[GlobalSSLDefaultBindOptions](#globalssldefaultbindoptions)_ | DefaultBindOptions sets default ssl-options to force on all binds. |


#### GlobalSSLDefaultBindOptions

_Underlying type:_ _[struct{MinVersion *string "json:\"minVersion,omitempty\""}](#struct{minversion-*string-"json:\"minversion,omitempty\""})_



_Appears in:_
- [GlobalSSL](#globalssl)



#### GlobalSSLTuneOptions

_Underlying type:_ _[struct{CacheSize *int64 "json:\"cacheSize,omitempty\""; Keylog string "json:\"keylog,omitempty\""; Lifetime *k8s.io/apimachinery/pkg/apis/meta/v1.Duration "json:\"lifetime,omitempty\""; ForcePrivateCache bool "json:\"forcePrivateCache,omitempty\""; MaxRecord *int64 "json:\"maxRecord,omitempty\""; DefaultDHParam int64 "json:\"defaultDHParam,omitempty\""; CtxCacheSize int64 "json:\"ctxCacheSize,omitempty\""; CaptureBufferSize *int64 "json:\"captureBufferSize,omitempty\""}](#struct{cachesize-*int64-"json:\"cachesize,omitempty\"";-keylog-string-"json:\"keylog,omitempty\"";-lifetime-*k8sioapimachinerypkgapismetav1duration-"json:\"lifetime,omitempty\"";-forceprivatecache-bool-"json:\"forceprivatecache,omitempty\"";-maxrecord-*int64-"json:\"maxrecord,omitempty\"";-defaultdhparam-int64-"json:\"defaultdhparam,omitempty\"";-ctxcachesize-int64-"json:\"ctxcachesize,omitempty\"";-capturebuffersize-*int64-"json:\"capturebuffersize,omitempty\""})_



_Appears in:_
- [GlobalTuneOptions](#globaltuneoptions)



#### GlobalTuneOptions





_Appears in:_
- [GlobalConfiguration](#globalconfiguration)

| Field | Description |
| --- | --- |
| `maxrewrite` _[int64](#int64)_ | Maxrewrite sets the reserved buffer space to this size in bytes. The reserved space is used for header rewriting or appending. The first reads on sockets will never fill more than bufsize-maxrewrite. |
| `bufsize` _[int64](#int64)_ | Bufsize sets the buffer size to this size (in bytes). Lower values allow more sessions to coexist in the same amount of RAM, and higher values allow some applications with very large cookies to work. |
| `ssl` _[GlobalSSLTuneOptions](#globalssltuneoptions)_ | SSL sets the SSL tune options. |


#### Instance



Instance is the Schema for the instances API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `proxy.haproxy.com/v1alpha1`
| `kind` _string_ | `Instance`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[InstanceSpec](#instancespec)_ |  |


#### InstancePhase

_Underlying type:_ _string_

InstancePhase is a label for the phase of a Instance at the current time.

_Appears in:_
- [InstanceStatus](#instancestatus)



#### InstanceSpec



InstanceSpec defines the desired state of Instance

_Appears in:_
- [Instance](#instance)

| Field | Description |
| --- | --- |
| `replicas` _integer_ | Replicas is the desired number of replicas of the HAProxy Instance. |
| `network` _[Network](#network)_ | Network contains the configuration of Route, Services and other network related configuration. |
| `configuration` _[Configuration](#configuration)_ | Configuration is used to bootstrap the global and defaults section of the HAProxy configuration. |
| `image` _string_ | Image specifies the HaProxy image including th tag. |
| `sidecars` _[Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#container-v1-core) array_ | Sidecars additional sidecar containers |
| `serviceAccountName` _string_ | ServiceAccountName is the name of the ServiceAccount to use to run this Instance. |
| `allowPrivilegedPorts` _boolean_ | AllowPrivilegedPorts allows to bind sockets with port numbers less than 1024. |
| `placement` _[Placement](#placement)_ | Placement define how the instance's pods should be scheduled. |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#pullpolicy-v1-core)_ | ImagePullPolicy one of Always, Never, IfNotPresent. |
| `metrics` _[Metrics](#metrics)_ | Metrics defines the metrics endpoint and scraping configuration. |
| `labels` _object (keys:string, values:string)_ | Labels additional labels for the ha-proxy pods |


#### Metrics





_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enabled will enable metrics globally for Instance. |
| `address` _string_ | Address to bind the metrics endpoint (default: '0.0.0.0'). |
| `port` _integer_ | Port specifies the port used for metrics. |
| `relabelings` _RelabelConfig array_ | RelabelConfigs to apply to samples before scraping. More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config |
| `interval` _[Duration](#duration)_ | Interval at which metrics should be scraped If not specified Prometheus' global scrape interval is used. |


#### Network





_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description |
| --- | --- |
| `hostNetwork` _boolean_ | HostNetwork will enable the usage of host network. |
| `hostIPs` _object (keys:string, values:string)_ | HostIPs defines an environment variable BIND_ADDRESS in the instance based on the provided host to IP mapping |
| `route` _[RouteSpec](#routespec)_ | Route defines the desired state for OpenShift Routes. |
| `service` _[ServiceSpec](#servicespec)_ | Service defines the desired state for a Service. |


#### Placement





_Appears in:_
- [InstanceSpec](#instancespec)

| Field | Description |
| --- | --- |
| `nodeSelector` _object (keys:string, values:string)_ | NodeSelector is a selector which must be true for the pod to fit on a node. |
| `topologySpreadConstraints` _[TopologySpreadConstraint](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.25/#topologyspreadconstraint-v1-core) array_ | TopologySpreadConstraints describes how a group of pods ought to spread across topology domains. Scheduler will schedule pods in a way which abides by the constraints. |


#### RouteSpec





_Appears in:_
- [Network](#network)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enabled will toggle the creation of OpenShift Routes. |
| `tls` _[TLSConfig](#tlsconfig)_ | TLS provides the ability to configure certificates and termination for the route. |


#### ServiceSpec





_Appears in:_
- [Network](#network)

| Field | Description |
| --- | --- |
| `enabled` _boolean_ | Enabled will toggle the creation of a Service. |


