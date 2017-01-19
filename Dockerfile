FROM tutum/curl
EXPOSE 8000
COPY /httpserver/static /httpserver/static
COPY /scripts /scripts
COPY jq-linux64 /bin/jq
RUN chmod +x /bin/jq
ADD babl-dashboard_linux_amd64 /bin/babl-dashboard
RUN chmod +x /bin/babl-dashboard
WORKDIR /

CMD ["/bin/babl-dashboard"]