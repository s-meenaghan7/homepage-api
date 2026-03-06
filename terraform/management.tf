# ==============================================================================
# Management resources for the homepage API
# ==============================================================================

resource "aws_servicecatalogappregistry_application" "homepage" {
  name        = "homepage"
  description = "Application registry for homepage resources"
}