FROM alpine:3.4
EXPOSE 8000
COPY /httpserver/static /httpserver/static
ADD babl-dashboard_linux_amd64 /bin/babl-dashboard
RUN chmod +x /bin/babl-dashboard
WORKDIR /
CMD ["/bin/babl-dashboard"]