{{- define "ponglehub.job" -}}
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    app: {{ required "must enter a name property!" .name }}
  name: {{ required "must enter a name property!" .name }}
spec:
  backoffLimit: {{ .backoffLimit | default 3 }}
  completions: {{ .completions | default 1 }}
  parallelism: {{ .parallelism | default 1 }}
  template:
    metadata:
      labels:
        app: {{ required "must enter a name property!" .name }}
    spec:
      containers:
      - name: {{ .containerName | default "job" }}
        image: {{ required "must enter an image property!" .image}}
        imagePullPolicy: {{ .pullPolicy | default "Always" }}
        {{- if .env }}
        env:
          {{- toYaml .env | nindent 10 }}
        {{- end }}
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
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: OnFailure
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
{{- end -}}