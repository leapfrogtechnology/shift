package backend_ha_architecture

const InfrastructureTemplate = `
// Terraform State Backend Initialization
terraform {
  backend "remote" {
    organization = "lftechnology"
    token = "{{ info.Token }}"
    workspaces {
      name = "{{ info.Client.Project }}-{{ info.Client.Deployment.Name }}-{{ info.Client.Deployment.Type }}"
    }
  }
}

// Variables
variable "az_count" {
  type = "string"
  default = "2"
}

variable "cidr_block" {
  default = "172.22.0.0/16"
}

variable "tags" {
  type = "map"
  default = {
    Name = "{{ info.Client.Deployment.Name }}"
    Project = "{{ info.Client.Project }}"
  }
}
variable "region" {
  default = "us-east-1"
}
variable "alb_target_name" {
  type = "string"
  default = "{{ info.Client.Project }}-{{ info.Client.Deployment.Name }}-target-group"
}

variable "health_check_path" {
  type = "string"
  default = "{{ info.Client.Deployment.HealthCheckPath }}"
}

variable "repo_name" {
  type = "string"
{% if info.Client.Deployment.RepoName %}
  default = "{{ info.Client.Deployment.RepoName }}"
{% else %}
  default = "leapfrogtechnology/shift-ui:latest"
{% endif %}
}

// Provider Initialization
provider "aws" {
  region = var.region
  access_key = "{{ info.Client.Deployment.AccessKey }}"
  secret_key = "{{ info.Client.Deployment.SecretKey }}"
}

# Fetch AZ in current Region
data "aws_availability_zones" "available" {}

data "aws_acm_certificate" "cert" {
  domain = "*.shift.lftechnology.com"
}
resource "aws_vpc" "main" {
  cidr_block = var.cidr_block
  tags = var.tags
}

# Create a Private Subnet
resource "aws_subnet" "private" {
  count = var.az_count
  cidr_block = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index)
  vpc_id = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[count.index]
  tags = {
    Name = var.tags.Name
    Project = var.tags.Project
    Public = false
  }
}

# Create a Public Subnet
resource "aws_subnet" "public" {
  count = var.az_count
  cidr_block = cidrsubnet(aws_vpc.main.cidr_block, 8, var.az_count + count.index)
  vpc_id = aws_vpc.main.id
  availability_zone = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true
  tags = {
    Name = var.tags.Name
    Project = var.tags.Project
    Public = true
  }
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

// ALB
# ALB Security Group: Edit this to restrict access to the application
resource "aws_security_group" "lb" {
  name = "{{ info.Client.Project|lower }}-{{ info.Client.Deployment.Name|lower }}-security-group"
  description = "Security Group for {{ info.Client.Project }} {{ info.Client.Deployment.Name }}"
  vpc_id = aws_vpc.main.id
  ingress {
    protocol = "tcp"
    from_port = 80
    to_port = 80
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  ingress {
    protocol = "tcp"
    from_port = 443
    to_port = 443
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  egress {
    protocol = "-1"
    from_port = 0
    to_port = 0
    cidr_blocks = [
      "0.0.0.0/0"]
  }
}


resource "aws_alb" "main" {
  name = "{{ info.Client.Project|lower }}-{{ info.Client.Deployment.Name|lower }}-LB"
  subnets = aws_subnet.public.*.id
  security_groups = [
    aws_security_group.lb.id]
  tags = var.tags
}

resource "aws_alb_target_group" "app" {
  name = var.alb_target_name
  port = 80
  protocol = "HTTP"
  vpc_id = aws_vpc.main.id
  target_type = "ip"
  tags = var.tags
  depends_on = [aws_alb.main]
  health_check {
    healthy_threshold = "3"
    interval = "30"
    protocol = "HTTP"
    matcher = "200"
    timeout = "3"
    path = var.health_check_path
    unhealthy_threshold = "2"
  }
}

# Redirect all traffic from the ALB to the target group
resource "aws_alb_listener" "my_website_https" {
 load_balancer_arn = aws_alb.main.id
 port = "443"
 protocol = "HTTPS"
 ssl_policy = "ELBSecurityPolicy-2016-08"
 certificate_arn = data.aws_acm_certificate.cert.arn
 //  tags        = var.tags
 default_action {
   type = "forward"
   target_group_arn = aws_alb_target_group.app.id
 }
}


resource "aws_lb_listener" "my_website_http" {
  load_balancer_arn = aws_alb.main.id
  port = "80"
  protocol = "HTTP"
  //default_action {
  //  type = "forward"
  //  target_group_arn = aws_alb_target_group.app.id
  //}
    default_action {
      type = "redirect"
      target_group_arn = aws_alb_target_group.app.id

    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
      host = "#{host}"
      path = "/#{path}"
      query = "#{query}"
    }
  }
}
//Fargate

variable "fargate_cluster_name" {
  type = "string"
  default = "{{ info.Client.Project|lower }}-{{ info.Client.Deployment.Name|lower }}"
}

variable "fargate_container_port" {
  type = "string"
  default = "{{ info.Client.Deployment.Port }}"
}

variable "ecr_name" {
  type = "string"
  default = "{{ info.Client.Project|lower }}/{{ info.Client.Deployment.Name|lower }}-backend"
}

resource "aws_iam_role" "ECSAutoScalingRole" {
  name = "ECSAutoScalingRole-{{ info.Client.Project|lower }}-{{ info.Client.Deployment.Type|lower }}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs.application-autoscaling.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ECSAutoScalingPolicy" {
  name        = "ECSAutoScalingPolicy-{{ info.Client.Project|lower }}-{{ info.Client.Deployment.Type|lower }}"
  description = "ECSAutoScalingPolicy"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ECSTaskManagement",
            "Effect": "Allow",
            "Action": [
                "ec2:AttachNetworkInterface",
                "ec2:CreateNetworkInterface",
                "ec2:CreateNetworkInterfacePermission",
                "ec2:DeleteNetworkInterface",
                "ec2:DeleteNetworkInterfacePermission",
                "ec2:Describe*",
                "ec2:DetachNetworkInterface",
                "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
                "elasticloadbalancing:DeregisterTargets",
                "elasticloadbalancing:Describe*",
                "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
                "elasticloadbalancing:RegisterTargets",
                "route53:ChangeResourceRecordSets",
                "route53:CreateHealthCheck",
                "route53:DeleteHealthCheck",
                "route53:Get*",
                "route53:List*",
                "route53:UpdateHealthCheck",
                "servicediscovery:DeregisterInstance",
                "servicediscovery:Get*",
                "servicediscovery:List*",
                "servicediscovery:RegisterInstance",
                "servicediscovery:UpdateInstanceCustomHealthStatus"
            ],
            "Resource": "*"
        },
        {
            "Sid": "ECSTagging",
            "Effect": "Allow",
            "Action": [
                "ec2:CreateTags"
            ],
            "Resource": "arn:aws:ec2:*:*:network-interface/*"
        },
        {
            "Sid": "ECSLogsStream",
            "Effect": "Allow",
            "Action": [
                "logs:*"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach-auto-scaling" {
  role       = aws_iam_role.ECSAutoScalingRole.name
  policy_arn = aws_iam_policy.ECSAutoScalingPolicy.arn
}

resource "aws_iam_role" "ECSTasksExecutionRole" {
  name = "ECSTasksExecutionRole-{{ info.Client.Project|lower }}-{{ info.Client.Deployment.Type|lower }}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ECSTasksExecutionPolicy" {
  name        = "ECSTasksExecutionPolicy-{{ info.Client.Project|lower }}-{{ info.Client.Deployment.Type|lower }}"
  description = "ECSTasksExecutionPolicy"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ecr:GetAuthorizationToken",
                "ecr:BatchCheckLayerAvailability",
                "ecr:GetDownloadUrlForLayer",
                "ecr:BatchGetImage",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach-tasks-execution" {
  role       = aws_iam_role.ECSTasksExecutionRole.name
  policy_arn = aws_iam_policy.ECSTasksExecutionPolicy.arn
}

variable "container_port" {
  default = {{ info.Client.Deployment.Port }}
}

data "template_file" "container_definition" {
  template = file("sample.json.tpl")
  vars = {
    fargate_container_name = var.fargate_cluster_name
    fargate_cluster_name = var.fargate_cluster_name
    fargate_container_port = var.container_port
    repo_name = var.repo_name
  }
}

module "fargate" {
  source = "github.com/sparshadotel/terraform-fargate-module"
  faragte_app_name = var.fargate_cluster_name
  fargate_cluster_name = var.fargate_cluster_name
  fargate_container_name = var.fargate_cluster_name
  fargate_alb_security_group_id = aws_security_group.lb.id
  fargate_alb_target_group_id = aws_alb_target_group.app.id
  fargate_autoscaling_role_arn = aws_iam_role.ECSAutoScalingRole.arn
  fargate_task_execution_role_arn = aws_iam_role.ECSTasksExecutionRole.arn
  fargate_container_definitions = data.template_file.container_definition.rendered
  fargate_container_port = var.container_port
  fargate_subnets = aws_subnet.private.*.id
  vpc_id = aws_vpc.main.id
  tags = var.tags
  cloudwatch_log_group = var.fargate_cluster_name

}

resource "aws_ecr_repository" "repo" {
  name = var.ecr_name
  tags = var.tags
}

output "repoUrl" {
  value = aws_ecr_repository.repo.repository_url
}

output "backendServiceId" {
  value = module.fargate.fargate_service_id
}

output "backendClusterName" {
  value = var.fargate_cluster_name
}

output "backendTaskDefinitionId" {
  value = module.fargate.fargate_task_definition
}

output "backendContainerDefinition" {
  value = module.fargate.fargate_container_definition
}

output "appUrl" {
  value = aws_alb.main.dns_name
}
// Template
`

const ContainerTemplate = `[
    {
        "name": "${fargate_container_name}",
        "image": "${repo_name}",
        "cpu": 0,
        "memoryReservation": 256,
        "portMappings": [
            {
                "containerPort": ${fargate_container_port},
                "hostPort": ${fargate_container_port},
                "protocol": "tcp"
            }
        ],
        "essential": true,
        "environment": [
            {
                "name": "REGION_NAME",
                "value": "us-east-1"
            }
        ],
        "mountPoints": [],
        "volumesFrom": [],
        "logConfiguration": {
            "logDriver": "awslogs",
            "options": {
                "awslogs-group": "/ecs/${fargate_cluster_name}",
                "awslogs-region": "us-east-1",
                "awslogs-stream-prefix": "ecs"
            }
        }
    }
]`
