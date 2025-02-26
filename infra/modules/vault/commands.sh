######################################################################

UNSEAL_KEY="<UNSEAL_KEY>"
kubectl exec -n vault vault-0 -- vault operator unseal ${UNSEAL_KEY}
kubectl exec -n vault vault-1 -- vault operator unseal ${UNSEAL_KEY}
kubectl exec -n vault vault-2 -- vault operator unseal ${UNSEAL_KEY}
