variable "project" {
  type = string
}

variable "name" {
  type = string
}

variable "location" {
  type = string
}

variable "protocol" {
  type        = string
  description = "grpc or http"
}

variable "envs" {
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "auth" {
  type    = bool
  default = true
}
