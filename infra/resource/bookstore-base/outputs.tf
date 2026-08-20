output "namespace" {
  description = "The bookstore namespace name"
  value       = kubernetes_namespace_v1.bookstore.metadata[0].name
}
