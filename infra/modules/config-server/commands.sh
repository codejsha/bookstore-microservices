######################################################################

helm package ./helm

helm push --plain-http config-server-*.tgz oci://harbor.example.com/main
