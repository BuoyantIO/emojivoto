FROM buoyantio/emojivoto-svc-base:v8

ARG svc_name

COPY $svc_name/target/ /usr/local/bin/

COPY --from=gcr.io/linkerd-io/await:v0.1.0 /linkerd-await /usr/local/bin
ENV NAME=$svc_name
ENTRYPOINT [ "linkerd-await", "--" ]
CMD [ "sh", "-c", "$NAME" ]
