{{- define "ponglehub.db.ca" -}}
{{- $root := . }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: db-ca
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: db-certs
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs: [ create ]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: db-certs
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: db-certs
subjects:
- kind: ServiceAccount
  name: db-ca
---
apiVersion: batch/v1
kind: Job
metadata:
  name: db-ca
spec:
  backoffLimit: 0
  template:
    metadata:
      annotations:
        linkerd.io/inject: disabled
    spec:
      restartPolicy: Never
      serviceAccountName: db-ca
      securityContext:
        runAsUser: 1000
        fsGroup: 2000
      volumes:
      - name: certs
        emptyDir: {}
      initContainers:
      - name: make-ca
        image: cockroachdb/cockroach:v21.1.11
        command: [ /bin/sh ]
        args:
        - -c
        - |-
          cockroach cert create-ca --certs-dir=/output/certs --ca-key=/output/ca.key
          cockroach cert create-node \
            'cockroachdb-public' \
            'cockroachdb-public.{{ $root.Release.Namespace }}' \
            'cockroachdb-public.{{ $root.Release.Namespace }}.svc.cluster.local' \
            '*.cockroachdb' \
            '*.cockroachdb.{{ $root.Release.Namespace }}' \
            '*.cockroachdb.{{ $root.Release.Namespace }}.svc.cluster.local' \
            --certs-dir=/output/certs \
            --ca-key=/output/ca.key
          cockroach cert create-client root --certs-dir=/output/certs --ca-key=/output/ca.key
        volumeMounts:
        - name: certs
          mountPath: /output
        resources:
          requests:
            memory: 128Mi
            cpu: 0.1
          limits:
            memory: 128Mi
            cpu: 0.1
      containers:
      - name: write-secret
        image: bitnami/kubectl
        command: [ /bin/bash ]
        args:
        - -c
        - |-
          kubectl create secret tls cockroach-ca --cert=/output/certs/ca.crt --key=/output/ca.key || true
          kubectl create secret generic cockroach-node \
            --from-file=ca.crt=/output/certs/ca.crt \
            --from-file=tls.crt=/output/certs/node.crt \
            --from-file=tls.key=/output/certs/node.key \
            --type=kubernetes.io/tls
          kubectl create secret tls cockroach-root --cert=/output/certs/client.root.crt --key=/output/certs/client.root.key || true
        volumeMounts:
        - name: certs
          mountPath: /output
        resources:
          requests:
            memory: 128Mi
            cpu: 0.1
          limits:
            memory: 128Mi
            cpu: 0.1
{{- end -}}

{{- define "ponglehub.db.user" -}}
{{- $username := index . 0 }}
{{- $root := index . 1 }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: db-certs-{{ $username }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: db-certs-{{ $username }}
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs: [ get, create ]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: db-certs-{{ $username }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: db-certs-{{ $username }}
subjects:
- kind: ServiceAccount
  name: db-certs-{{ $username }}
  namespace: {{ $root.Release.Namespace }}
---
apiVersion: batch/v1
kind: Job
metadata:
  name: db-certs-{{ $username }}
spec:
  backoffLimit: 0
  template:
    metadata:
      annotations:
        linkerd.io/inject: disabled
    spec:
      restartPolicy: Never
      serviceAccountName: db-certs-{{ $username }}
      securityContext:
        runAsUser: 1000
        fsGroup: 2000
      volumes:
      - name: certs
        emptyDir: {}
      initContainers:
      - name: read-secret
        image: bitnami/kubectl
        command: [ /bin/bash ]
        args:
        - -c
        - |-
          mkdir -p /output/certs
          kubectl get secret -n auth-service cockroachdb-ca -o jsonpath="{.data.ca\.key}" | base64 -d > output/ca.key
          kubectl get secret -n auth-service cockroachdb-node -o jsonpath="{.data.ca\.crt}" | base64 -d > output/certs/ca.crt
        volumeMounts:
        - name: certs
          mountPath: /output
        resources:
          requests:
            memory: 128Mi
            cpu: 0.1
          limits:
            memory: 128Mi
            cpu: 0.1
      - name: make-certs
        image: cockroachdb/cockroach:v21.1.11
        command: [ /bin/sh ]
        args:
        - -c
        - |-
          cockroach cert create-client {{ $username }} --certs-dir=/output/certs --ca-key=/output/ca.key
        volumeMounts:
        - name: certs
          mountPath: /output
        resources:
          requests:
            memory: 128Mi
            cpu: 0.1
          limits:
            memory: 128Mi
            cpu: 0.1
      containers:
      - name: write-secret
        image: bitnami/kubectl
        command: [ /bin/bash ]
        args:
        - -c
        - |-
          kubectl create secret generic cockroachdb-{{ $username }} \
            --from-file=ca.crt=/output/certs/ca.crt \
            --from-file=client.{{ $username }}.crt=/output/certs/client.{{ $username }}.crt \
            --from-file=client.{{ $username }}.key=/output/certs/client.{{ $username }}.key || true
        volumeMounts:
        - name: certs
          mountPath: /output
        resources:
          requests:
            memory: 128Mi
            cpu: 0.1
          limits:
            memory: 128Mi
            cpu: 0.1
{{- end -}}