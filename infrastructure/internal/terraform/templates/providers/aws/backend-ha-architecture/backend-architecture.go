package backend_ha_architecture

const InfrastructureTemplate = `
// Terraform State Backend Initialization
terraform {
  backend "remote" {
    organization = "lftechnology"
    token = "{{ info.TERRAFORM_TOKEN }}"
    workspaces {
      name = "{{ info.CLIENT_NAME }}-backend"
    }
  }
}
// Provider Initialization
provider "aws" {
  region     = "{{ info.AWS_REGION }}"
  access_key = "{{ info.AWS_ACCESS_KEY }}"
  secret_key = "{{ info.AWS_SECRET_KEY }}"
}

# Variables
variable "az_count" {
  type = "string"
  default = "2"
}

variable "tags" {
  type = "map"
  default = {
    Name = "{{ info.RESOURCE_NAME }}"
    Project = "{{ info.PROJECT_NAME }}"
  }
}

# Fetch AZ in current Region
data "aws_availability_zones" "available" {}

resource "aws_vpc" "main" {
  cidr_block = "{{ info.CIDR_BLOCK }}"
  tags = var.tags
}

# Create a Private Subnet
resource "aws_subnet" "private" {
  count = var.az_count
  cidr_block = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index)
  vpc_id = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[count.index]
}

# Create a Public Subnet
resource "aws_subnet" "public" {
  count = var.az_count
  cidr_block = cidrsubnet(aws_vpc.main.cidr_block, 8, var.az_count + count.index)
  vpc_id = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
}

# Internet Gateway for Public Subnet
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

# Route the Public Subnet through IGW
resource "aws_route" "internet_access" {
  route_table_id = aws_vpc.main.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id = aws_internet_gateway.gw.id
}

# Create a NAT gateway with an EIP for each private subnet to get internet connectivity
resource "aws_eip" "gw" {
  count = var.az_count
  vpc = true
  depends_on = [
    "aws_internet_gateway.gw"
  ]
}

resource "aws_nat_gateway" "gw" {
  count = var.az_count
  subnet_id = element(aws_subnet.public.*.id, count.index)
  allocation_id = element(aws_eip.gw.*.id, count.index)
}

# Create a new route table for the private subnets, make it route non-local traffic through the NAT gateway to the internet
resource "aws_route_table" "private" {
  count = var.az_count
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway_id = element(aws_nat_gateway.gw.*.id, count.index)
  }
}

# Explicitly associate the newly created route tables to the private subnets (so they don't default to the main route table)
resource "aws_route_table_association" "private" {
  count = var.az_count
  subnet_id = element(aws_subnet.private.*.id, count.index)
  route_table_id = element(aws_route_table.private.*.id, count.index)
}

# Outputs

output "vpc_id" {
  value = aws_vpc.main.id
}

output "vpc_cidr_block" {
  value = aws_vpc.main.cidr_block
}

output "private_subnets" {
  value = aws_subnet.private.*.id
}

output "public_subnets" {
  value = aws_subnet.public.*.id
}
`
