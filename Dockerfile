FROM busybox
LISTEN 8000
ADD babl-dashboard_linux_amd64 /bin/babl-dashboard
CMD ["/bin/babl-dashboard"]