{{- define "ponglehub.static" -}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ required "must enter a name property!" .name }}
  labels:
    app: {{ .name }}
spec:
  selector:
    matchLabels:
      app: {{ .name }}
  template:
    metadata:
      labels:
        app: {{ .name }}
    spec:
      containers:
      - name: server
        image: {{ required "must enter an image property!"  .image }}
        imagePullPolicy: {{ .pullPolicy | default "Always" }}
        ports:
        - containerPort: {{ .port | default 80 }}
          name: http
          protocol: TCP
        resources:
        {{- if .resources }}
          {{- toYaml .resources | nindent 10 }}
        {{- else }}
          requests:
            memory: 32Mi
            cpu: 0.1
          limits:
            memory: 32Mi
            cpu: 0.1
        {{- end }}
      restartPolicy: Always
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: {{ .name }}
  name: {{ .name }}
spec:
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: http
  selector:
    app: {{ .name }}
  type: ClusterIP
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    ingress.kubernetes.io/ssl-redirect: "true"
    kubernetes.io/ingress.class: traefik
  labels:
    app: {{ .name }}
  name: {{ .name }}
spec:
  rules:
  - host: {{ required "must provide a host!" .host }}
    http:
      paths:
      - backend:
          serviceName: {{ .name }}
          servicePort: 80
        path: /
        pathType: ImplementationSpecific
  tls:
  - hosts:
    - {{ .host }}
{{- end -}}