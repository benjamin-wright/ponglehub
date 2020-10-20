{{- define "ponglehub.server" -}}
{{ $server := first . }}
{{ $top := (index . 1) }}
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
---
{{- if not (empty $server.host) }}
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    ingress.kubernetes.io/ssl-redirect: "true"
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
{{- end -}}