module "cloud_run_fur_meds_pets_api" {
  source = "./cloud_run"

  project  = var.project
  location = var.region

  name  = "furmeds-api"
  image = "europe-west4-docker.pkg.dev/fur-meds-project/fur-meds/furmeds-api:${var.app_version}"

  envs = [
    {
      name  = "API_PORT"
      value = 80
    },
    {
      name  = "GIN_MODE"
      value = "release"
    }
  ]
}
