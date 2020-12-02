{{- define "ponglehub.server" -}}
{{ $server := first . }}
{{ $top := (index . 1) }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ $server.name }}
automountServiceAccountToken: true
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: {{ required "must enter a name property!" $server.name }}
  name: {{ required "must enter a name property!" $server.name }}
spec:
  replicas: {{ $server.replicas | default 1 }}
  selector:
    matchLabels:
      app: {{ $server.name }}
  template:
    metadata:
      labels:
        app: {{ required "must enter a name property!" $server.name }}
      {{- if $server.annotations }}
      annotations:
      {{- range $key, $value := $server.annotations }}
        {{ $key }}: {{ $value }}
      {{- end }}
      {{- end }}
    spec:
      serviceAccountName: {{ $server.name }}
      {{- if $server.initContainers }}
      initContainers:
      {{- range $key, $container := $server.initContainers }}
      - name: {{ $container.name }}
        image: {{ required "must enter an image property!" $container.image}}
        imagePullPolicy: {{ $container.pullPolicy | default "Always" }}
        {{- if $container.env }}
        env:
        {{- range $name, $value := $container.env }}
        - name: {{ $name }}
          value: {{ tpl $value $top | quote }}
        {{- end }}
        {{- end }}
        resources:
        {{- if $server.resources }}
          {{- toYaml $server.resources | nindent 10 }}
        {{- else }}
          requests:
            memory: 32Mi
            cpu: 0.1
          limits:
            memory: 32Mi
            cpu: 0.1
        {{- end }}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      {{- end }}
      {{- end}}
      containers:
      - name: server
        image: {{ required "must enter an image property!" $server.image}}
        imagePullPolicy: {{ $server.pullPolicy | default "Always" }}
        {{- if $server.command }}
        command:
          {{- range $server.command }}
        - {{ . }}
          {{- end }}
        {{- end }}
        {{- if $server.args }}
        args:
          {{- range $server.args }}
        - {{ . | quote }}
          {{- end }}
        {{- end }}
        {{- if $server.env }}
        env:
        {{- range $key, $value := $server.env }}
        - name: {{ $key }}
          value: {{ tpl $value $top | quote }}
        {{- end }}
        {{- end }}
        resources:
        {{- if $server.resources }}
          {{- toYaml $server.resources | nindent 10 }}
        {{- else }}
          requests:
            memory: 32Mi
            cpu: 0.1
          limits:
            memory: 32Mi
            cpu: 0.1
        {{- end }}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30

---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: {{ $server.name }}
  name: {{ $server.name }}
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: {{ $server.port | default 80 }}
  selector:
    app: {{ $server.name }}
  type: ClusterIP
{{- if not (empty $server.host) }}
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    ingress.kubernetes.io/ssl-redirect: "true"
    ingress.kubernetes.io/auth-url: http://gatekeeper.ponglehub.svc.cluster.local
    ingress.kubernetes.io/auth-type: forward
  labels:
    app: {{ $server.name }}
  name: {{ $server.name }}
spec:
  rules:
  - host: {{ required "must provide a host!" $server.host }}
    http:
      paths:
      - backend:
          serviceName: {{ $server.name }}
          servicePort: 80
        path: /
        pathType: ImplementationSpecific
  tls:
  - hosts:
    - {{ $server.host }}
    secretName: ssl-secret
{{- end }}
{{- if not (empty $server.rbac) }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ $server.name }}
rules:
- apiGroups:
  {{- range $server.rbac.apiGroups }}
  - {{ . }}
  {{- end }}
  resources:
  {{- range $server.rbac.resources }}
  - {{ . }}
  {{- end }}
  verbs:
  {{- range $server.rbac.verbs }}
  - {{ . }}
  {{- end }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ $server.name }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ $server.name }}
subjects:
- kind: ServiceAccount
  name: {{ $server.name }}
  namespace: {{ $top.Release.Namespace }}
{{- end }}
{{- end -}}