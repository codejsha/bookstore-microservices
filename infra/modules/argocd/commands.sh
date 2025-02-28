######################################################################

argocd app delete -y admin-command
argocd app delete -y admin-query
argocd app delete -y catalog-command
argocd app delete -y catalog-query
argocd app delete -y customer-command
argocd app delete -y customer-query
argocd app delete -y order-command
argocd app delete -y order-query
argocd app delete -y payment-command
argocd app delete -y payment-query
argocd proj delete bookstore
