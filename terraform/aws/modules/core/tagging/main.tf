locals {
  original_tags = join(var.delimiter, compact(concat(tolist([var.customer, var.product, var.environment]), var.attributes)))
}

locals {
  transformed_tags = var.convert_case ? lower(local.original_tags) : local.original_tags
}

locals {
  id = var.enabled ? local.transformed_tags : ""

  customer    = var.enabled ? (var.convert_case ? lower(format("%v", var.customer)) : format("%v", var.customer)) : ""
  product     = var.enabled ? (var.convert_case ? lower(format("%v", var.product)) : format("%v", var.product)) : ""
  environment = var.enabled ? (var.convert_case ? lower(format("%v", var.environment)) : format("%v", var.environment)) : ""
  attributes  = var.enabled ? (var.convert_case ? lower(format("%v", join(var.delimiter, compact(var.attributes)))) : format("%v", join(var.delimiter, compact(var.attributes)))) : ""

  # Merge input tags with our tags.
  # Note: `Name` has a special meaning in AWS and we need to disamgiuate it by using the computed `id`
  tags = merge(
    tomap({
      "Name"        = local.id,
      "Customer"    = local.customer,
      "Product"     = local.product,
      "Environment" = local.environment,
      "ManagedBy"   = "Terraform",
      "CostCentre" = var.cost_centre }
    ), var.tags
  )
}
