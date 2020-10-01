{{- define "ponglehub.job" -}}
{{ $job := first . }}
{{ $top := (index . 1) }}
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    app: {{ required "must enter a name property!" $job.name }}
  name: {{ required "must enter a name property!" $job.name }}
spec:
  backoffLimit: {{ $job.backoffLimit | default 0 }}
  completions: {{ $job.completions | default 1 }}
  parallelism: {{ $job.parallelism | default 1 }}
  template:
    metadata:
      labels:
        app: {{ required "must enter a name property!" $job.name }}
      annotations:
        linkerd.io/inject: disabled
    spec:
      {{- if $job.initContainers }}
      initContainers:
      {{- range $key, $container := $job.initContainers }}
      - name: {{ $container.name }}
        image: {{ required "must enter an image property!" $container.image}}
        imagePullPolicy: {{ $container.pullPolicy | default "Always" }}
        {{- if $container.env }}
        env:
          {{- tpl $container.env $top | nindent 10 }}
        {{- end }}
        resources:
        {{- if $job.resources }}
          {{- toYaml $job.resources | nindent 10 }}
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
      - name: {{ $job.name | default "job" }}
        image: {{ required "must enter an image property!" $job.image}}
        imagePullPolicy: {{ $job.pullPolicy | default "Always" }}
        {{- if $job.env }}
        env:
          {{- tpl $job.env $top | nindent 10 }}
        {{- end }}
        resources:
        {{- if $job.resources }}
          {{- toYaml $job.resources | nindent 10 }}
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
      restartPolicy: Never
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
{{- end -}}