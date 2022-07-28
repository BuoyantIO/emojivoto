FROM datawire/emojivoto-svc-base:v11.1

ARG svc_name

COPY $svc_name/target/ /usr/local/bin/

COPY .ambassador-internal/openapi.yaml /.ambassador-internal/openapi.yaml

# ARG variables arent available for ENTRYPOINT
ENV SVC_NAME $svc_name
ENTRYPOINT cd /usr/local/bin && $SVC_NAME
