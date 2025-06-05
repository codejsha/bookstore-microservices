output "harbor_users" {
  value = tolist([
    harbor_user.dev_admin.username,
    harbor_user.devops_admin.username
  ])
}
