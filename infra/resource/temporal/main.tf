terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 3.2"
    }
  }
}

resource "kubernetes_job_v1" "register_namespace" {
  for_each = toset(var.temporal_namespaces)

  wait_for_completion = true

  metadata {
    name      = "temporal-ns-${each.key}"
    namespace = var.namespace
  }

  spec {
    backoff_limit           = 10
    active_deadline_seconds = 600

    template {
      metadata {
        labels = {
          "app.kubernetes.io/name"      = "temporal-ns-register"
          "app.kubernetes.io/component" = "namespace-bootstrap"
        }
      }

      spec {
        restart_policy = "OnFailure"

        container {
          name    = "register"
          image   = var.admintools_image
          command = ["/bin/sh", "-ec"]
          env {
            name  = "TEMPORAL_ADDRESS"
            value = var.temporal_address
          }
          args = [<<-EOT
            ns="${each.key}"
            echo "Waiting for Temporal frontend at $TEMPORAL_ADDRESS ..."
            i=0
            until temporal operator cluster health --address "$TEMPORAL_ADDRESS" >/dev/null 2>&1; do
              i=$((i+1))
              if [ "$i" -ge 60 ]; then
                echo "Temporal frontend not reachable in time" >&2
                exit 1
              fi
              sleep 5
            done
            if temporal operator namespace describe --namespace "$ns" --address "$TEMPORAL_ADDRESS" >/dev/null 2>&1; then
              echo "Temporal namespace '$ns' already exists."
            else
              echo "Creating Temporal namespace '$ns' (retention ${var.retention}) ..."
              temporal operator namespace create --namespace "$ns" --retention "${var.retention}" --address "$TEMPORAL_ADDRESS"
              echo "Created Temporal namespace '$ns'."
            fi
          EOT
          ]
        }
      }
    }
  }
}
